package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

// Config 应用配置结构
type Config struct {
	Server     ServerConfig     `yaml:"server"`
	Database   DatabaseConfig   `yaml:"database"`
	Redis      RedisConfig      `yaml:"redis"`
	JWT        JWTConfig        `yaml:"jwt"`
	Crawler    CrawlerConfig    `yaml:"crawler"`
	NLPService NLPServiceConfig `yaml:"nlp_service"`
	Log        LogConfig        `yaml:"log"`
	Bootstrap  BootstrapConfig  `yaml:"bootstrap"`
	CORS       CORSConfig       `yaml:"cors"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Mode         string `yaml:"mode"`
	ReadTimeout  int    `yaml:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	DBName          string `yaml:"dbname"`
	SSLMode         string `yaml:"sslmode"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// GetDSN 返回PostgreSQL连接字符串
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode,
	)
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
	PoolSize int    `yaml:"pool_size"`
}

// GetAddr 返回Redis地址
func (r *RedisConfig) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

// GetExpireDuration 返回过期时间
func (j *JWTConfig) GetExpireDuration() time.Duration {
	return time.Duration(j.ExpireHours) * time.Hour
}

// CrawlerConfig 爬虫配置
type CrawlerConfig struct {
	WorkerCount    int       `yaml:"worker_count"`
	RequestTimeout int       `yaml:"request_timeout"`
	UserAgents     []string  `yaml:"user_agents"`
	RateLimit      RateLimit `yaml:"rate_limit"`
}

// RateLimit 频率限制配置
type RateLimit struct {
	RequestsPerSecond float64 `yaml:"requests_per_second"`
	Burst             int     `yaml:"burst"`
}

// NLPServiceConfig NLP服务配置
type NLPServiceConfig struct {
	URL     string `yaml:"url"`
	Timeout int    `yaml:"timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	Output     string `yaml:"output"`
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age"`
}

// BootstrapConfig 启动阶段配置
type BootstrapConfig struct {
	SeedCompetitionsOnStart *bool `yaml:"seed_competitions_on_start"`
}

// ShouldSeedCompetitions 判断启动时是否执行 competitions seed
func (b BootstrapConfig) ShouldSeedCompetitions() bool {
	if b.SeedCompetitionsOnStart == nil {
		return true
	}
	return *b.SeedCompetitionsOnStart
}

// CORSConfig CORS 配置
type CORSConfig struct {
	AllowedOrigins  []string `yaml:"allowed_origins"`
	AllowCredentials bool    `yaml:"allow_credentials"`
	AllowedHeaders  []string `yaml:"allowed_headers"`
	AllowedMethods  []string `yaml:"allowed_methods"`
}

// globalConfig 全局配置实例
// 注意: 此变量在Load()函数中设置，在Get()函数中读取
// 并发安全: Load()通常在应用启动时调用一次，之后只读，因此不需要加锁
// 如果需要在运行时动态更新配置，需要添加sync.RWMutex保护
var globalConfig *Config

// Load 加载配置文件
// 配置加载优先级:
// 1. 尝试加载 .env(.local) 文件（仅填充未设置的环境变量）
// 2. 如果configPath为空，先读 CONFIG_PATH，再根据 APP_ENV 选择配置文件（默认dev）
// 3. 读取YAML配置文件
// 4. 使用os.ExpandEnv解析环境变量（支持${VAR}格式）
// 5. 解析YAML内容到Config结构体
// 6. 用环境变量覆盖关键配置
// 7. 应用默认值并验证配置
func Load(configPath string) (*Config, error) {
	loadDotEnvFiles()

	// 如果路径为空，根据环境变量选择配置文件
	if configPath == "" {
		configPath = resolveConfigPath()
	}

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析环境变量（支持${VAR}格式，如 ${DATABASE_PASSWORD}）
	expandedData := os.ExpandEnv(string(data))

	// 解析YAML
	var cfg Config
	if err := yaml.Unmarshal([]byte(expandedData), &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := applyEnvOverrides(&cfg); err != nil {
		return nil, fmt.Errorf("failed to apply env overrides: %w", err)
	}
	applyDefaults(&cfg)

	// 验证配置
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	globalConfig = &cfg
	return &cfg, nil
}

// Get 获取全局配置
func Get() *Config {
	if globalConfig == nil {
		panic("config not loaded, call Load() first")
	}
	return globalConfig
}

func resolveConfigPath() string {
	if configPath, ok := lookupEnv("CONFIG_PATH"); ok {
		return configPath
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	return filepath.Join("configs", fmt.Sprintf("config.%s.yaml", env))
}

func loadDotEnvFiles() {
	for _, path := range []string{
		".env.local",
		".env",
		filepath.Join("..", ".env.local"),
		filepath.Join("..", ".env"),
	} {
		_ = godotenv.Load(path)
	}
}

func applyDefaults(c *Config) {
	if c.Server.Host == "" {
		c.Server.Host = "0.0.0.0"
	}
	if c.Server.Port <= 0 {
		c.Server.Port = 8080
	}
	if c.Server.Mode == "" {
		c.Server.Mode = "debug"
	}
	if c.Server.ReadTimeout <= 0 {
		c.Server.ReadTimeout = 60
	}
	if c.Server.WriteTimeout <= 0 {
		c.Server.WriteTimeout = 60
	}

	if c.Database.Port <= 0 {
		c.Database.Port = 5432
	}
	if c.Database.SSLMode == "" {
		c.Database.SSLMode = "disable"
	}
	if c.Database.MaxOpenConns <= 0 {
		c.Database.MaxOpenConns = 25
	}
	if c.Database.MaxIdleConns <= 0 {
		c.Database.MaxIdleConns = 5
	}
	if c.Database.ConnMaxLifetime <= 0 {
		c.Database.ConnMaxLifetime = 300
	}

	if c.Redis.Port <= 0 {
		c.Redis.Port = 6379
	}
	if c.Redis.PoolSize <= 0 {
		c.Redis.PoolSize = 10
	}

	if c.JWT.ExpireHours <= 0 {
		c.JWT.ExpireHours = 168
	}

	if c.Crawler.WorkerCount <= 0 {
		c.Crawler.WorkerCount = 5
	}
	if c.Crawler.RequestTimeout <= 0 {
		c.Crawler.RequestTimeout = 30
	}
	if c.Crawler.RateLimit.RequestsPerSecond <= 0 {
		c.Crawler.RateLimit.RequestsPerSecond = 2
	}
	if c.Crawler.RateLimit.Burst <= 0 {
		c.Crawler.RateLimit.Burst = 5
	}

	if c.NLPService.Timeout <= 0 {
		c.NLPService.Timeout = 60
	}

	if c.Log.Level == "" {
		c.Log.Level = "info"
	}
	if c.Log.Output == "" {
		c.Log.Output = "stdout"
	}
	if c.Log.FilePath == "" {
		c.Log.FilePath = "logs/app.log"
	}

	if len(c.CORS.AllowedOrigins) == 0 {
		c.CORS.AllowedOrigins = []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		}
	}
	if len(c.CORS.AllowedHeaders) == 0 {
		c.CORS.AllowedHeaders = []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Origin",
			"Cache-Control",
			"X-Requested-With",
		}
	}
	if len(c.CORS.AllowedMethods) == 0 {
		c.CORS.AllowedMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}
}

func applyEnvOverrides(c *Config) error {
	var err error

	setStringFromEnv(&c.Server.Host, "SERVER_HOST")
	if err = setIntFromEnv(&c.Server.Port, "SERVER_PORT", "API_PORT"); err != nil {
		return err
	}
	setStringFromEnv(&c.Server.Mode, "SERVER_MODE")
	if err = setIntFromEnv(&c.Server.ReadTimeout, "SERVER_READ_TIMEOUT", "API_READ_TIMEOUT"); err != nil {
		return err
	}
	if err = setIntFromEnv(&c.Server.WriteTimeout, "SERVER_WRITE_TIMEOUT", "API_WRITE_TIMEOUT"); err != nil {
		return err
	}

	setStringFromEnv(&c.Database.Host, "DB_HOST")
	if err = setIntFromEnv(&c.Database.Port, "DB_PORT"); err != nil {
		return err
	}
	setStringFromEnv(&c.Database.User, "DB_USER")
	setStringFromEnv(&c.Database.Password, "DB_PASSWORD")
	setStringFromEnv(&c.Database.DBName, "DB_NAME")
	setStringFromEnv(&c.Database.SSLMode, "DB_SSLMODE")
	if err = setIntFromEnv(&c.Database.MaxOpenConns, "DB_MAX_OPEN_CONNS"); err != nil {
		return err
	}
	if err = setIntFromEnv(&c.Database.MaxIdleConns, "DB_MAX_IDLE_CONNS"); err != nil {
		return err
	}
	if err = setIntFromEnv(&c.Database.ConnMaxLifetime, "DB_CONN_MAX_LIFETIME"); err != nil {
		return err
	}

	setStringFromEnv(&c.Redis.Host, "REDIS_HOST")
	if err = setIntFromEnv(&c.Redis.Port, "REDIS_PORT"); err != nil {
		return err
	}
	setStringFromEnv(&c.Redis.Password, "REDIS_PASSWORD")
	if err = setIntFromEnv(&c.Redis.DB, "REDIS_DB"); err != nil {
		return err
	}
	if err = setIntFromEnv(&c.Redis.PoolSize, "REDIS_POOL_SIZE"); err != nil {
		return err
	}

	setStringFromEnv(&c.JWT.Secret, "JWT_SECRET")
	if err = setIntFromEnv(&c.JWT.ExpireHours, "JWT_EXPIRE_HOURS"); err != nil {
		return err
	}

	if err = setIntFromEnv(&c.Crawler.WorkerCount, "CRAWLER_WORKER_COUNT"); err != nil {
		return err
	}
	if err = setIntFromEnv(&c.Crawler.RequestTimeout, "CRAWLER_REQUEST_TIMEOUT"); err != nil {
		return err
	}

	setStringFromEnv(&c.NLPService.URL, "NLP_SERVICE_URL")
	if err = setIntFromEnv(&c.NLPService.Timeout, "NLP_SERVICE_TIMEOUT"); err != nil {
		return err
	}

	setStringFromEnv(&c.Log.Level, "LOG_LEVEL")
	setStringFromEnv(&c.Log.Output, "LOG_OUTPUT")
	setStringFromEnv(&c.Log.FilePath, "LOG_FILE_PATH")

	setStringSliceFromEnv(&c.CORS.AllowedOrigins, "CORS_ALLOWED_ORIGINS")
	if err = setBoolFromEnv(&c.CORS.AllowCredentials, "CORS_ALLOW_CREDENTIALS"); err != nil {
		return err
	}
	setStringSliceFromEnv(&c.CORS.AllowedHeaders, "CORS_ALLOWED_HEADERS")
	setStringSliceFromEnv(&c.CORS.AllowedMethods, "CORS_ALLOWED_METHODS")
	if err = setBoolPtrFromEnv(&c.Bootstrap.SeedCompetitionsOnStart, "SEED_COMPETITIONS_ON_START"); err != nil {
		return err
	}

	return nil
}

func lookupEnv(keys ...string) (string, bool) {
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if !ok {
			continue
		}
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		return trimmed, true
	}
	return "", false
}

func setStringFromEnv(target *string, keys ...string) {
	if value, ok := lookupEnv(keys...); ok {
		*target = value
	}
}

func setIntFromEnv(target *int, keys ...string) error {
	value, ok := lookupEnv(keys...)
	if !ok {
		return nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("invalid integer for %s: %q", keys[0], value)
	}
	*target = parsed
	return nil
}

func setBoolFromEnv(target *bool, keys ...string) error {
	value, ok := lookupEnv(keys...)
	if !ok {
		return nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("invalid boolean for %s: %q", keys[0], value)
	}
	*target = parsed
	return nil
}

func setBoolPtrFromEnv(target **bool, keys ...string) error {
	value, ok := lookupEnv(keys...)
	if !ok {
		return nil
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fmt.Errorf("invalid boolean for %s: %q", keys[0], value)
	}
	*target = &parsed
	return nil
}

func setStringSliceFromEnv(target *[]string, keys ...string) {
	value, ok := lookupEnv(keys...)
	if !ok {
		return
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		result = append(result, trimmed)
	}

	if len(result) > 0 {
		*target = result
	}
}

// validate 验证配置有效性
func (c *Config) validate() error {
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	if c.Server.Mode != "debug" && c.Server.Mode != "release" && c.Server.Mode != "test" {
		return fmt.Errorf("invalid server mode: %s", c.Server.Mode)
	}

	if c.Database.Host == "" {
		return fmt.Errorf("database host cannot be empty")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user cannot be empty")
	}
	if c.Database.DBName == "" {
		return fmt.Errorf("database name cannot be empty")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("invalid database port: %d", c.Database.Port)
	}

	if c.Redis.Host == "" {
		return fmt.Errorf("redis host cannot be empty")
	}
	if c.Redis.Port <= 0 || c.Redis.Port > 65535 {
		return fmt.Errorf("invalid redis port: %d", c.Redis.Port)
	}

	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT secret cannot be empty")
	}
	if c.Server.Mode == "release" && len(c.JWT.Secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 chars in release mode")
	}

	if c.NLPService.URL == "" {
		return fmt.Errorf("nlp service url cannot be empty")
	}

	if len(c.CORS.AllowedOrigins) == 0 {
		return fmt.Errorf("cors allowed_origins cannot be empty")
	}
	if c.CORS.AllowCredentials && containsWildcard(c.CORS.AllowedOrigins) {
		return fmt.Errorf("cors allowed_origins cannot include '*' when allow_credentials=true")
	}

	return nil
}

func containsWildcard(origins []string) bool {
	for _, origin := range origins {
		if strings.TrimSpace(origin) == "*" {
			return true
		}
	}
	return false
}
