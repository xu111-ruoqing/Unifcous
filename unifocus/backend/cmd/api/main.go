package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unifocus/backend/internal/api/handlers"
	"github.com/unifocus/backend/internal/api/middleware"
	"github.com/unifocus/backend/internal/config"
	"github.com/unifocus/backend/internal/crawler"
	"github.com/unifocus/backend/internal/infrastructure/client"
	"github.com/unifocus/backend/internal/seed"
	"github.com/unifocus/backend/internal/repository/postgres"
	"github.com/unifocus/backend/internal/repository/redis"
	"github.com/unifocus/backend/internal/service"
	"github.com/unifocus/backend/pkg/jwt"
	"github.com/unifocus/backend/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Log); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting UniFocus API Server...")
	logger.Infof("Environment: %s", cfg.Server.Mode)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	db, err := postgres.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("Failed to close database: %v", err)
		}
	}()

	// Seed competitions master data (idempotent via name_key).
	if cfg.Bootstrap.ShouldSeedCompetitions() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		if err := seed.SeedCompetitions(ctx, db); err != nil {
			cancel()
			logger.Fatalf("Failed to seed competitions: %v", err)
		}
		cancel()
		logger.Info("Competitions seed completed")
	} else {
		logger.Info("Competitions seed skipped by configuration")
	}

	// 初始化Redis连接
	rdb, err := redis.NewClient(&cfg.Redis, "unifocus")
	if err != nil {
		logger.Fatalf("Failed to initialize redis: %v", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			logger.Errorf("Failed to close redis: %v", err)
		}
	}()

	// 初始化服务层
	userRepo := postgres.NewUserRepository(db)
	oppRepo := postgres.NewOpportunityRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	crawlTaskRepo := postgres.NewCrawlTaskRepository(db)
	recRepo := postgres.NewRecognitionRepository(db)
	competitionRepo := postgres.NewCompetitionRepository(db)

	scheduler := crawler.NewScheduler(crawlTaskRepo, nil) // Pipeline paused for next step

	jwtMgr := jwt.NewManager(&cfg.JWT)
	authService := service.NewAuthService(userRepo, jwtMgr)
	recService := service.NewRecognitionService(recRepo)
	oppService := service.NewOpportunityService(oppRepo, recService)
	scoringService := service.NewScoringService()

	nlpClient := client.NewNLPClient(&cfg.NLPService)
	profileService := service.NewProfileService(profileRepo, nlpClient)

	crawlerService := service.NewCrawlerService(crawlTaskRepo, scheduler)

	// 创建路由（传入数据库和Redis实例供后续使用）
	router := setupRouter(cfg, db, rdb, authService, oppService, profileService, crawlerService, competitionRepo, scoringService, oppRepo, userRepo, profileRepo)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 启动服务器（优雅关闭）
	go func() {
		logger.Infof("Server is running on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭（5秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

// setupRouter 设置路由
func setupRouter(cfg *config.Config, db *postgres.DB, rdb *redis.Client, authService *service.AuthService, oppService *service.OpportunityService, profileService *service.ProfileService, crawlerService service.CrawlerService, competitionRepo *postgres.CompetitionRepository, scoringService *service.ScoringService, oppRepo *postgres.OpportunityRepository, userRepo *postgres.UserRepository, profileRepo *postgres.ProfileRepository) *gin.Engine {
	router := gin.New()

	// 中间件
	router.Use(gin.Recovery())
	router.Use(corsMiddleware(cfg.CORS))
	router.Use(middleware.MetricsMiddleware())
	router.Use(loggerMiddleware())

	// 初始化handlers
	authHandler := handlers.NewAuthHandler(authService)
	oppHandler := handlers.NewOpportunityHandler(oppService)
	profileHandler := handlers.NewProfileHandler(profileService)
	metricsHandler := handlers.NewMetricsHandler()
	crawlTaskHandler := handlers.NewCrawlTaskHandler(crawlerService)
	competitionHandler := handlers.NewCompetitionHandler(competitionRepo)
	scoringHandler := handlers.NewScoringHandler(scoringService, oppRepo, userRepo, profileRepo)

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "1.0.0",
			"time":    time.Now().Unix(),
		})
	})

	router.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		if err := db.HealthCheck(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"error":  fmt.Sprintf("database health check failed: %v", err),
			})
			return
		}
		if err := rdb.HealthCheck(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not_ready",
				"error":  fmt.Sprintf("redis health check failed: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
			"time":   time.Now().Unix(),
		})
	})

	// 监控指标
	router.GET("/api/v1/metrics", metricsHandler.GetMetrics)

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 认证路由（无需JWT）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 公开的机会查询路由（无需认证）
		v1.GET("/opportunities", oppHandler.List)
		v1.GET("/opportunities/:id", oppHandler.GetByID)
		v1.GET("/opportunities/:id/score", scoringHandler.ScoreOpportunity)
		v1.GET("/opportunities/scored", scoringHandler.ScoreList)

		// 竞赛管理（无需认证）
		v1.GET("/competitions", competitionHandler.List)
		v1.GET("/competitions/:id", competitionHandler.GetByID)
		v1.POST("/competitions", competitionHandler.Create)
		v1.PUT("/competitions/:id", competitionHandler.Update)
		v1.DELETE("/competitions/:id", competitionHandler.Delete)

		// 需要认证的路由
		authorized := v1.Group("")
		authorized.Use(middleware.AuthMiddleware(authService))
		{
			// 机会管理（需要认证）
			authorized.POST("/opportunities", oppHandler.Create)
			authorized.PUT("/opportunities/:id", oppHandler.Update)
			authorized.DELETE("/opportunities/:id", oppHandler.Delete)

			// 用户画像管理
			authorized.GET("/users/me/profile", profileHandler.GetProfile)
			authorized.PUT("/users/me/profile", profileHandler.UpdateProfile)
			authorized.POST("/users/me/profile/resume", profileHandler.UploadResume)

			// 爬虫任务管理
			crawl := authorized.Group("/crawl-tasks")
			{
				crawl.POST("", crawlTaskHandler.Create)
				crawl.GET("", crawlTaskHandler.List)
				crawl.GET("/:id", crawlTaskHandler.GetByID)
				crawl.PUT("/:id", crawlTaskHandler.Update)
				crawl.DELETE("/:id", crawlTaskHandler.Delete)
				crawl.POST("/start", crawlTaskHandler.StartScheduler)
				crawl.POST("/stop", crawlTaskHandler.StopScheduler)
			}
		}
	}

	return router
}

// corsMiddleware CORS中间件
func corsMiddleware(cfg config.CORSConfig) gin.HandlerFunc {
	allowHeaders := strings.Join(cfg.AllowedHeaders, ", ")
	allowMethods := strings.Join(cfg.AllowedMethods, ", ")

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowedOrigin := ""

		if origin != "" {
			if isOriginAllowed(origin, cfg.AllowedOrigins) {
				allowedOrigin = origin
			} else if isOriginAllowed("*", cfg.AllowedOrigins) {
				allowedOrigin = "*"
			} else if c.Request.Method == http.MethodOptions {
				c.AbortWithStatus(http.StatusForbidden)
				return
			}
		} else if isOriginAllowed("*", cfg.AllowedOrigins) {
			allowedOrigin = "*"
		}

		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			c.Writer.Header().Add("Vary", "Origin")
			if cfg.AllowCredentials && allowedOrigin != "*" {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethods)

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if strings.TrimSpace(allowedOrigin) == origin {
			return true
		}
	}
	return false
}

// loggerMiddleware 日志中间件
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Infof("[%s] %s %s %d %v",
			method,
			path,
			clientIP,
			statusCode,
			latency,
		)
	}
}
