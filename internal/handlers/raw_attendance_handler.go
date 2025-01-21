package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/models"
	"point-system-api/internal/services"
)

// RawAttendanceHandler handles HTTP requests for raw attendance-related operations.
type RawAttendanceHandler struct {
	rawAttendanceService services.RawAttendanceService
}

// NewRawAttendanceHandler creates a new instance of RawAttendanceHandler.
func NewRawAttendanceHandler(rawAttendanceService services.RawAttendanceService) *RawAttendanceHandler {
	return &RawAttendanceHandler{
		rawAttendanceService: rawAttendanceService,
	}
}

// CreateRawAttendance handles the creation of a new raw attendance record.
func (h *RawAttendanceHandler) CreateRawAttendance(c *gin.Context) {
	var rawAttendance struct {
		models.RawAttendance
		UserID    string `json:"userID"`
		WorkDayID string `json:"workDayID"`
		Punch     string `json:"punch"`
		Status    string
	}

	if err := c.ShouldBindJSON(&rawAttendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload " + err.Error()})
		return
	}

	UserID, err := strconv.Atoi(rawAttendance.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UserID"})
		return
	}

	WorkDayID, err := strconv.Atoi(rawAttendance.WorkDayID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid WorkDayID"})
		return
	}

	Punch, err := strconv.Atoi(rawAttendance.Punch)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Punch"})
		return
	}

	rawAttendance.RawAttendance.Punch = uint(Punch)
	rawAttendance.RawAttendance.UserID = uint(UserID)
	rawAttendance.RawAttendance.WorkDayID = uint(WorkDayID)
	rawAttendance.RawAttendance.Status = 1

	if err := h.rawAttendanceService.CreateRawAttendance(c.Request.Context(), &rawAttendance.RawAttendance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("CREATE_RAWATTENDANCE")
	c.JSON(http.StatusCreated, gin.H{"message": "Raw attendance created successfully"})
}

// CreateManyRawAttendances handles the creation of multiple raw attendance records.
func (h *RawAttendanceHandler) CreateManyRawAttendances(c *gin.Context) {
	var rawAttendances []*models.RawAttendance
	if err := c.ShouldBindJSON(&rawAttendances); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.rawAttendanceService.CreateManyRawAttendances(c.Request.Context(), rawAttendances); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("CREATE_RAWATTENDANCE")
	c.JSON(http.StatusCreated, gin.H{"message": "Raw attendances created successfully"})
}

// GetRawAttendanceByID retrieves a raw attendance record by its ID.
func (h *RawAttendanceHandler) GetRawAttendanceByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid raw attendance ID"})
		return
	}

	rawAttendance, err := h.rawAttendanceService.GetRawAttendanceByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if rawAttendance == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Raw attendance not found"})
		return
	}

	c.JSON(http.StatusOK, rawAttendance)
}

// GetRawAttendancesByWorkDayID retrieves all raw attendance records for a specific work day.
func (h *RawAttendanceHandler) GetRawAttendancesByWorkDayID(c *gin.Context) {
	workDayID, err := strconv.Atoi(c.Param("workDayID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work day ID"})
		return
	}

	rawAttendances, err := h.rawAttendanceService.GetRawAttendancesByWorkDayID(c.Request.Context(), uint(workDayID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rawAttendances)
}

// UpdateRawAttendance updates an existing raw attendance record.
func (h *RawAttendanceHandler) UpdateRawAttendance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid raw attendance ID"})
		return
	}

	var rawAttendance models.RawAttendance
	if err := c.ShouldBindJSON(&rawAttendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	rawAttendance.ID = uint(id)
	if err := h.rawAttendanceService.UpdateRawAttendance(c.Request.Context(), &rawAttendance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("UPDATE_RAWATTENDANCE")
	c.JSON(http.StatusOK, gin.H{"message": "Raw attendance updated successfully"})
}

// DeleteRawAttendance deletes a raw attendance record by its ID.
func (h *RawAttendanceHandler) DeleteRawAttendance(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid raw attendance ID"})
		return
	}

	if err := h.rawAttendanceService.DeleteRawAttendance(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("DELETE_RAWATTENDANCE")
	c.JSON(http.StatusOK, gin.H{"message": "Raw attendance deleted successfully"})
}
