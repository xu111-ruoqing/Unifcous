package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unifocus/backend/internal/domain"
	"github.com/unifocus/backend/internal/service"
)

// OpportunityHandler handles opportunity HTTP requests
type OpportunityHandler struct {
	oppService *service.OpportunityService
}

// NewOpportunityHandler creates a new opportunity handler
func NewOpportunityHandler(oppService *service.OpportunityService) *OpportunityHandler {
	return &OpportunityHandler{
		oppService: oppService,
	}
}

// Create handles opportunity creation
// @Summary Create a new opportunity
// @Description Create a new opportunity (requires authentication)
// @Tags opportunities
// @Accept json
// @Produce json
// @Param request body domain.CreateOpportunityRequest true "Opportunity creation request"
// @Success 201 {object} domain.Opportunity
// @Failure 400 {object} map[string]string
// @Router /api/v1/opportunities [post]
func (h *OpportunityHandler) Create(c *gin.Context) {
	var req domain.CreateOpportunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	opp, err := h.oppService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, opp)
}

// GetByID handles getting an opportunity by ID
// @Summary Get opportunity by ID
// @Description Retrieve an opportunity by its ID
// @Tags opportunities
// @Produce json
// @Param id path int true "Opportunity ID"
// @Success 200 {object} domain.Opportunity
// @Failure 404 {object} map[string]string
// @Router /api/v1/opportunities/{id} [get]
func (h *OpportunityHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid opportunity ID"})
		return
	}

	profileID, _ := strconv.ParseInt(c.Query("profile_id"), 10, 64) // Optional, 0 if error or empty

	opp, err := h.oppService.GetByID(c.Request.Context(), id, profileID)
	if err != nil {
		if err.Error() == "opportunity not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, opp)
}

// List handles listing opportunities with filtering
// @Summary List opportunities
// @Description Retrieve a list of opportunities with optional filtering and pagination
// @Tags opportunities
// @Produce json
// @Param type query string false "Opportunity type"
// @Param competition_level query string false "Competition level"
// @Param major query string false "Target major"
// @Param deadline_after query string false "Deadline after date (YYYY-MM-DD)"
// @Param deadline_before query string false "Deadline before date (YYYY-MM-DD)"
// @Param is_active query bool false "Is active"
// @Param tags query []string false "Tags"
// @Param limit query int false "Limit (default: 20)"
// @Param offset query int false "Offset (default: 0)"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/opportunities [get]
func (h *OpportunityHandler) List(c *gin.Context) {
	var filter domain.OpportunityFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	opportunities, total, err := h.oppService.List(c.Request.Context(), &filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   opportunities,
		"total":  total,
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// Update handles opportunity update
// @Summary Update an opportunity
// @Description Update an existing opportunity (requires authentication)
// @Tags opportunities
// @Accept json
// @Produce json
// @Param id path int true "Opportunity ID"
// @Param request body domain.CreateOpportunityRequest true "Opportunity update request"
// @Success 200 {object} domain.Opportunity
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/opportunities/{id} [put]
func (h *OpportunityHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid opportunity ID"})
		return
	}

	var req domain.CreateOpportunityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	opp, err := h.oppService.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err.Error() == "opportunity not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, opp)
}

// Delete handles opportunity deletion
// @Summary Delete an opportunity
// @Description Soft delete an opportunity (requires authentication)
// @Tags opportunities
// @Produce json
// @Param id path int true "Opportunity ID"
// @Success 204 "No Content"
// @Failure 404 {object} map[string]string
// @Router /api/v1/opportunities/{id} [delete]
func (h *OpportunityHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid opportunity ID"})
		return
	}

	if err := h.oppService.Delete(c.Request.Context(), id); err != nil {
		if err.Error() == "opportunity not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
