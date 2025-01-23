package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/services"
)

type AttendanceHandler struct {
	attendanceService services.AttendanceService
}

// NewAttendanceHandler creates a new instance of AttendanceHandler
func NewAttendanceHandler(attendanceService services.AttendanceService) *AttendanceHandler {
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

	manager.broadcast <- []byte("CREATE_ATTENDANCELOG")
	c.JSON(http.StatusOK, gin.H{
		"message": "Attendance log saved successfully",
		"data":    attendanceLog,
	})
}

// GetAttendanceLogByID retrieves a specific attendance log by its ID
func (h *AttendanceHandler) GetAttendanceLogByID(c *gin.Context) {
	id := c.Param("id")

	// Convert id from string to uint
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	attendanceLog, err := h.attendanceService.GetAttendanceByID(c.Request.Context(), uint(parsedID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance log not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": attendanceLog})
}

// GetAllAttendanceLogs retrieves all attendance logs with optional filters
func (h *AttendanceHandler) ListAttendanceLogs(c *gin.Context) {
	// Extract query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	startDate := c.Query("startDate")
	endDate := c.Query("endDate")
	punch := c.Query("punch")
	filters := map[string]interface{}{}

	// Add search filter if provided
	if search != "" {
		filters["search"] = search
	}

	if startDate != "" {
		filters["start_date"] = startDate
	}

	if endDate != "" {
		filters["end_date"] = endDate
	}

	if punch != "" {
		filters["punch"] = punch
	}

	// Add filters from query parameters
	if serial := c.Query("serial_number"); serial != "" {
		filters["serial_number"] = serial
	}
	if companyID := c.Query("company_id"); companyID != "" {
		companyIDInt, err := strconv.Atoi(companyID)
		if err == nil {
			filters["company_id"] = companyIDInt
		}
	}

	// Get attendance logs with filters and pagination
	attendanceLogs, total, err := h.attendanceService.GetAllAttendanceLogs(c.Request.Context(), page, limit, filters, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  attendanceLogs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// DeleteAttendanceLog deletes an attendance log by its ID
func (h *AttendanceHandler) DeleteAttendanceLog(c *gin.Context) {
	id := c.Param("id")
	// Convert id from string to uint
	parsedID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	err = h.attendanceService.DeleteAttendanceLog(c.Request.Context(), uint(parsedID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendance log"})
		return
	}

	manager.broadcast <- []byte("DELETE_ATTENDANCELOG")
	c.JSON(http.StatusOK, gin.H{"message": "Attendance log deleted successfully"})
}
