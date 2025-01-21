package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/models"
	"point-system-api/internal/services"
)

// WorkDayHandler handles HTTP requests for workday-related operations.
type WorkDayHandler struct {
	workDayService services.WorkDayService
}

// NewWorkDayHandler creates a new instance of WorkDayHandler.
func NewWorkDayHandler(workDayService services.WorkDayService) *WorkDayHandler {
	return &WorkDayHandler{
		workDayService: workDayService,
	}
}

// CreateWorkDay handles the creation of a new workday.
func (h *WorkDayHandler) CreateWorkDay(c *gin.Context) {
	var workday models.WorkDay
	if err := c.ShouldBindJSON(&workday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload " + err.Error()})
		return
	}

	if err := h.workDayService.CreateWorkDay(c.Request.Context(), &workday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("CREATE_WORKDAY")
	c.JSON(http.StatusCreated, gin.H{"message": "Workday created successfully"})
}

// GetWorkDayByID retrieves a workday by its ID.
func (h *WorkDayHandler) GetWorkDayByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workday ID"})
		return
	}

	workday, err := h.workDayService.GetWorkDayByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if workday == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workday not found"})
		return
	}

	c.JSON(http.StatusOK, workday)
}

// ListWorkDays retrieves all workdays.
func (h *WorkDayHandler) ListWorkDays(c *gin.Context) {
	workdays, err := h.workDayService.ListWorkDays(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workdays)
}

// UpdateWorkDay handles updating a workday by its ID.
func (h *WorkDayHandler) UpdateWorkDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workday ID"})
		return
	}

	var workday models.WorkDay
	if err := c.ShouldBindJSON(&workday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	workday.ID = uint(id)
	if err := h.workDayService.UpdateWorkDay(c.Request.Context(), &workday); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("UPDATE_WORKDAY")
	c.JSON(http.StatusOK, gin.H{"message": "Workday updated successfully"})
}

// DeleteWorkDay handles deleting a workday by its ID.
func (h *WorkDayHandler) DeleteWorkDay(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workday ID"})
		return
	}

	if err := h.workDayService.DeleteWorkDay(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("DELETE_WORKDAY")
	c.JSON(http.StatusOK, gin.H{"message": "Workday deleted successfully"})
}
