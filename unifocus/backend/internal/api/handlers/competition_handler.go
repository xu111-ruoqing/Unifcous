package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unifocus/backend/internal/domain"
	"github.com/unifocus/backend/internal/repository/postgres"
)

type CompetitionHandler struct {
	repo *postgres.CompetitionRepository
}

func NewCompetitionHandler(repo *postgres.CompetitionRepository) *CompetitionHandler {
	return &CompetitionHandler{repo: repo}
}

type competitionDTO struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Level             string `json:"level"`
	OfficialURL       string `json:"official_url"`
	TypicalTimeWindow string `json:"typical_time_window"`
	TimelineHint      string `json:"timeline_hint"`
	Note              string `json:"note"`
}

type competitionInput struct {
	Name              string `json:"name" binding:"required"`
	Level             string `json:"level"`
	OfficialURL       string `json:"official_url"`
	TypicalTimeWindow string `json:"typical_time_window"`
	TimelineHint      string `json:"timeline_hint"`
	Note              string `json:"note"`
}

func toDTO(c domain.Competition) competitionDTO {
	return competitionDTO{
		ID:                c.ID,
		Name:              c.Name,
		Level:             c.Level,
		OfficialURL:       c.OfficialURL,
		TypicalTimeWindow: c.TypicalTimeWindow,
		TimelineHint:      c.TimelineHint,
		Note:              c.Note,
	}
}

func (h *CompetitionHandler) List(c *gin.Context) {
	q := c.Query("q")
	level := c.Query("level")
	limit := parseIntWithDefault(c.Query("limit"), 200)
	offset := parseIntWithDefault(c.Query("offset"), 0)

	items, err := h.repo.ListCompetitions(c.Request.Context(), q, level, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	out := make([]competitionDTO, 0, len(items))
	for _, it := range items {
		out = append(out, toDTO(it))
	}
	c.JSON(http.StatusOK, out)
}

func (h *CompetitionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	comp, err := h.repo.GetCompetition(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if comp == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, toDTO(*comp))
}

func (h *CompetitionHandler) Create(c *gin.Context) {
	var input competitionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comp := domain.Competition{
		Name:              input.Name,
		Level:             input.Level,
		OfficialURL:       input.OfficialURL,
		TypicalTimeWindow: input.TypicalTimeWindow,
		TimelineHint:      input.TimelineHint,
		Note:              input.Note,
	}
	if err := h.repo.CreateCompetition(c.Request.Context(), &comp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, toDTO(comp))
}

func (h *CompetitionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var input competitionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comp := domain.Competition{
		ID:                id,
		Name:              input.Name,
		Level:             input.Level,
		OfficialURL:       input.OfficialURL,
		TypicalTimeWindow: input.TypicalTimeWindow,
		TimelineHint:      input.TimelineHint,
		Note:              input.Note,
	}
	if err := h.repo.UpdateCompetition(c.Request.Context(), &comp); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toDTO(comp))
}

func (h *CompetitionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.repo.DeleteCompetition(c.Request.Context(), id); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

func parseIntWithDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 0 {
		return def
	}
	return n
}
