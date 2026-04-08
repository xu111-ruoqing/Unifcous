package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unifocus/backend/internal/domain"
	"github.com/unifocus/backend/internal/repository/postgres"
	"github.com/unifocus/backend/internal/service"
)

// ScoringHandler 可达性评分 HTTP 处理器
type ScoringHandler struct {
	scoringService *service.ScoringService
	oppRepo        *postgres.OpportunityRepository
	userRepo       *postgres.UserRepository
	profileRepo    *postgres.ProfileRepository
}

// NewScoringHandler 创建评分处理器
func NewScoringHandler(
	scoringService *service.ScoringService,
	oppRepo *postgres.OpportunityRepository,
	userRepo *postgres.UserRepository,
	profileRepo *postgres.ProfileRepository,
) *ScoringHandler {
	return &ScoringHandler{
		scoringService: scoringService,
		oppRepo:        oppRepo,
		userRepo:       userRepo,
		profileRepo:    profileRepo,
	}
}

// ScoreOpportunity 对单个机会计算可达性评分
// GET /api/v1/opportunities/:id/score?user_id=xxx
func (h *ScoringHandler) ScoreOpportunity(c *gin.Context) {
	oppID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid opportunity ID"})
		return
	}

	opp, err := h.oppRepo.GetByID(c.Request.Context(), oppID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "opportunity not found"})
		return
	}

	user, profile := h.resolveUserContext(c)
	result := h.scoringService.Score(c.Request.Context(), opp, profile, user)

	c.JSON(http.StatusOK, gin.H{
		"opportunity_id": oppID,
		"score_result":   result,
	})
}

// ScoreList 批量评分并返回排序后的机会列表
// GET /api/v1/opportunities/scored?user_id=xxx&limit=20&offset=0
func (h *ScoringHandler) ScoreList(c *gin.Context) {
	var filter domain.OpportunityFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}

	opps, total, err := h.oppRepo.List(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, profile := h.resolveUserContext(c)
	scored := h.scoringService.ScoreBatch(c.Request.Context(), opps, profile, user)

	c.JSON(http.StatusOK, gin.H{
		"data":   scored,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// resolveUserContext 从查询参数解析用户上下文（可选）
func (h *ScoringHandler) resolveUserContext(c *gin.Context) (*domain.User, *domain.UserProfile) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		return nil, nil
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		return nil, nil
	}

	user, err := h.userRepo.GetByID(c.Request.Context(), userID)
	if err != nil {
		return nil, nil
	}

	profile, err := h.profileRepo.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		return user, nil
	}

	return user, profile
}
