package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/services"
)

type AttendanceHandler struct {
	attendanceService *services.AttendanceService
}

// NewAttendanceHandler creates a new instance of AttendanceHandler
func NewAttendanceHandler(attendanceService *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{attendanceService: attendanceService}
}

// CreateAttendanceLog handles processing of hex data and saving the attendance log
func (h *AttendanceHandler) CreateAttendanceLog(c *gin.Context) {
	var requestBody struct {
		SerialNumber string `json:"serial_number"` // Serial number of the device
		HexData      string `json:"hex_data"`      // Hex data to process
	}

	// Bind the JSON request body to the struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Process the hex data and create the attendance log
	attendanceLog, err := h.attendanceService.CreateAttendanceLog(c.Request.Context(), requestBody.SerialNumber, requestBody.HexData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attendance log saved successfully",
		"data":    attendanceLog,
	})
}

// GetAttendanceLogByID retrieves a specific attendance log by its ID
func (h *AttendanceHandler) GetAttendanceLogByID(c *gin.Context) {
	id := c.Param("id")

	attendanceLog, err := h.attendanceService.GetAttendanceLogByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attendanceLog})
}

// GetAllAttendanceLogs retrieves all attendance logs with optional filters
func (h *AttendanceHandler) GetAllAttendanceLogs(c *gin.Context) {
	// Extract filters from query parameters
	filters := map[string]interface{}{}
	if serial := c.Query("serial_number"); serial != "" {
		filters["serial_number"] = serial
	}
	if userID := c.Query("user_id"); userID != "" {
		filters["user_id"] = userID
	}

	attendanceLogs, err := h.attendanceService.GetAllAttendanceLogs(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendance logs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attendanceLogs})
}

// DeleteAttendanceLog deletes an attendance log by its ID
func (h *AttendanceHandler) DeleteAttendanceLog(c *gin.Context) {
	id := c.Param("id")

	err := h.attendanceService.DeleteAttendanceLog(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendance log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attendance log deleted successfully"})
}
