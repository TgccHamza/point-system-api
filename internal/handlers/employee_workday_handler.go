package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/services"
)

// EmployeeWorkDayHandler handles HTTP requests for employee workday-related operations.
type EmployeeWorkDayHandler struct {
	employeeWorkDayService services.EmployeeWorkDayService
}

// NewEmployeeWorkDayHandler creates a new instance of EmployeeWorkDayHandler.
func NewEmployeeWorkDayHandler(employeeWorkDayService services.EmployeeWorkDayService) *EmployeeWorkDayHandler {
	return &EmployeeWorkDayHandler{
		employeeWorkDayService: employeeWorkDayService,
	}
}

// GenerateEmployeeWorkDay generates EmployeeWorkDay records based on WorkDay and RawAttendance.
func (h *EmployeeWorkDayHandler) GenerateEmployeeWorkDay(c *gin.Context) {
	workDayID, err := strconv.Atoi(c.Param("workDayID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid work day ID"})
		return
	}

	if err := h.employeeWorkDayService.GenerateEmployeeWorkDay(c.Request.Context(), uint(workDayID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee work days generated successfully"})
}

// UpdateEmployeeWorkDay updates an existing EmployeeWorkDay record with notes and status.
func (h *EmployeeWorkDayHandler) UpdateEmployeeWorkDay(c *gin.Context) {
	employeeWorkDayID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee work day ID"})
		return
	}

	var request struct {
		Notes  string `json:"notes"`
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.employeeWorkDayService.UpdateEmployeeWorkDay(c.Request.Context(), uint(employeeWorkDayID), request.Notes, request.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("UPDATE_WORKDAYEMPLOYEES")
	c.JSON(http.StatusOK, gin.H{"message": "Employee work day updated successfully"})
}
