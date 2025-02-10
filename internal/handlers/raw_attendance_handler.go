package handlers

import (
	"database/sql"
	"net/http"
	"point-system-api/internal/models"
	"point-system-api/internal/services"
	"point-system-api/internal/types"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type RawAttendanceHandler struct {
	rawAttendanceService services.RawAttendanceService
}

func transformRawAttendance(ra *models.RawAttendance) types.RawAttendanceResponse {
	var totalHours *float64
	if ra.TotalHours.Valid {
		totalHours = &ra.TotalHours.Float64
	}

	var employeeName *string
	if ra.EmployeeName.Valid {
		employeeName = &ra.EmployeeName.String
	}

	var position *string
	if ra.Position.Valid {
		position = &ra.Position.String
	}

	var startAt *string
	if ra.StartAt.Valid {
		startAt = &ra.StartAt.String
	}

	var endAt *string
	if ra.EndAt.Valid {
		endAt = &ra.EndAt.String
	}

	var status *string
	if ra.Status.Valid {
		status = &ra.Status.String
	}

	var notes *string
	if ra.Notes.Valid {
		notes = &ra.Notes.String
	}

	createdAt := ra.CreatedAt.Format(time.RFC3339)
	updatedAt := ra.UpdatedAt.Format(time.RFC3339)

	workDayID := ra.WorkDayID
	companyID := ra.CompanyID
	userID := ra.UserID

	return types.RawAttendanceResponse{
		ID:           ra.ID,
		CreatedAt:    &createdAt,
		UpdatedAt:    &updatedAt,
		WorkDayID:    &workDayID,
		CompanyID:    &companyID,
		UserID:       &userID,
		EmployeeName: employeeName,
		Position:     position,
		StartAt:      startAt,
		EndAt:        endAt,
		TotalHours:   totalHours,
		Status:       status,
		Notes:        notes,
	}
}

func NewRawAttendanceHandler(rawAttendanceService services.RawAttendanceService) *RawAttendanceHandler {
	return &RawAttendanceHandler{
		rawAttendanceService: rawAttendanceService,
	}
}

func (h *RawAttendanceHandler) CreateRawAttendance(c *gin.Context) {
	var rawAttendance models.RawAttendance
	if err := c.ShouldBindJSON(&rawAttendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.rawAttendanceService.CreateRawAttendance(c.Request.Context(), &rawAttendance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := transformRawAttendance(&rawAttendance)

	c.JSON(http.StatusCreated, response)
}

func (h *RawAttendanceHandler) GetRawAttendanceByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	rawAttendance, err := h.rawAttendanceService.GetRawAttendanceByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := transformRawAttendance(rawAttendance)
	c.JSON(http.StatusCreated, response)
}

func (h *RawAttendanceHandler) GetRawAttendancesByCompanyAndWorkDay(c *gin.Context) {
	companyID, err := strconv.ParseUint(c.Param("companyId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	workDayID, err := strconv.ParseUint(c.Param("workDayId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workday ID"})
		return
	}

	rawAttendances, err := h.rawAttendanceService.GetRawAttendancesByCompanyIDAndWorkDay(c.Request.Context(), uint(companyID), uint(workDayID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Transform each attendance record before returning
	var responses []types.RawAttendanceResponse
	for _, ra := range rawAttendances {
		responses = append(responses, transformRawAttendance(ra))
	}

	c.JSON(http.StatusOK, responses)
}

func (h *RawAttendanceHandler) UpdateRawAttendance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Struct for only updating the specified fields
	var payload struct {
		StartAt string `json:"start_at"`
		EndAt   string `json:"end_at"`
		Notes   string `json:"notes"`
		Status  string `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Create a RawAttendance instance with only the updated fields
	rawAttendance := models.RawAttendance{
		StartAt: sql.NullString{String: payload.StartAt, Valid: payload.StartAt != ""},
		EndAt:   sql.NullString{String: payload.EndAt, Valid: payload.EndAt != ""},
		Notes:   sql.NullString{String: payload.Notes, Valid: payload.Notes != ""},
		Status:  sql.NullString{String: payload.Status, Valid: payload.Status != ""},
	}

	if err := h.rawAttendanceService.UpdateRawAttendance(c.Request.Context(), &rawAttendance, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := transformRawAttendance(&rawAttendance)
	c.JSON(http.StatusCreated, response)
}

func (h *RawAttendanceHandler) DeleteRawAttendance(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.rawAttendanceService.DeleteRawAttendance(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Raw attendance deleted successfully"})
}

func (h *RawAttendanceHandler) ListRawAttendances(c *gin.Context) {
	rawAttendances, err := h.rawAttendanceService.ListRawAttendances(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rawAttendances)
}
