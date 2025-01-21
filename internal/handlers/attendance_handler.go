package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/services"
)

// AttendanceHandler handles HTTP requests for attendance-related operations.
type AttendanceHandler struct {
	attendanceService services.AttendanceService
}

// NewAttendanceHandler creates a new instance of AttendanceHandler.
func NewAttendanceHandler(attendanceService services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{
		attendanceService: attendanceService,
	}
}

// ProcessHexData handles the processing of hex data.
func (h *AttendanceHandler) ProcessHexData(c *gin.Context) {
	var requestBody struct {
		SerialNumber string `json:"serial_number"` // Serial number of the device
		HexData      string `json:"hex_data"`      // Hex data to process
	}

	// Bind the JSON request body to the struct
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Process the hex data and save the attendance log
	attendanceLog, err := h.attendanceService.CreateAttendanceLog(c.Request.Context(), requestBody.SerialNumber, requestBody.HexData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the response as JSON
	c.JSON(http.StatusOK, gin.H{
		"message": "Attendance log saved successfully",
		"data":    attendanceLog,
	})
}
