package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/unifocus/backend/internal/domain"
	"github.com/unifocus/backend/internal/service"
)

// CrawlTaskHandler handles HTTP requests for crawl tasks
type CrawlTaskHandler struct {
	service service.CrawlerService
}

// NewCrawlTaskHandler creates a new crawl task handler
func NewCrawlTaskHandler(service service.CrawlerService) *CrawlTaskHandler {
	return &CrawlTaskHandler{service: service}
}

// Create creates a new crawl task
func (h *CrawlTaskHandler) Create(c *gin.Context) {
	var task domain.CrawlTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateTask(c.Request.Context(), &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetByID retrieves a crawl task by ID
func (h *CrawlTaskHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	task, err := h.service.GetTask(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// List retrieves a list of crawl tasks
func (h *CrawlTaskHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	tasks, err := h.service.ListTasks(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list tasks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  tasks,
		"page":  page,
		"limit": limit,
	})
}

// Update updates a crawl task
func (h *CrawlTaskHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var task domain.CrawlTask
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task.ID = id

	if err := h.service.UpdateTask(c.Request.Context(), &task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Delete deletes a crawl task
func (h *CrawlTaskHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := h.service.DeleteTask(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// StartScheduler starts the crawler scheduler
func (h *CrawlTaskHandler) StartScheduler(c *gin.Context) {
	h.service.StartScheduler()
	c.JSON(http.StatusOK, gin.H{"message": "Scheduler started"})
}

// StopScheduler stops the crawler scheduler
func (h *CrawlTaskHandler) StopScheduler(c *gin.Context) {
	h.service.StopScheduler()
	c.JSON(http.StatusOK, gin.H{"message": "Scheduler stopped"})
}
