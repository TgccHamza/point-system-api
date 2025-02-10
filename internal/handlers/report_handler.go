package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"point-system-api/internal/services"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportService services.ReportService
}

func NewReportHandler(reportService services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) GenerateReport(c *gin.Context) {
	companyID := c.Param("companyID")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Parse companyID
	companyIDUint, err := strconv.ParseUint(companyID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
		return
	}

	// Generate report
	report, err := h.reportService.GenerateReport(context.Background(), uint(companyIDUint), startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate report"})
		return
	}

	// Respond with JSON
	c.JSON(http.StatusOK, report)
}
