package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/models"
	"point-system-api/internal/services"
)

// EmployeeWorkdayHandler handles HTTP requests for employee workday-related operations.
type EmployeeWorkdayHandler struct {
	workdayService services.EmployeeWorkdayService
}

// NewEmployeeWorkdayHandler creates a new instance of EmployeeWorkdayHandler.
func NewEmployeeWorkdayHandler(workdayService services.EmployeeWorkdayService) *EmployeeWorkdayHandler {
	return &EmployeeWorkdayHandler{
		workdayService: workdayService,
	}
}

// CreateEmployeeWorkday handles the creation of a new employee workday.
func (h *EmployeeWorkdayHandler) CreateEmployeeWorkday(c *gin.Context) {
	var workday models.EmployeeWorkday
	if err := c.ShouldBindJSON(&workday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	workdayID, err := h.workdayService.CreateEmployeeWorkday(c.Request.Context(), workday)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": workdayID, "message": "Employee workday created successfully"})
}

// GetEmployeeWorkdayByID retrieves an employee workday by its ID.
func (h *EmployeeWorkdayHandler) GetEmployeeWorkdayByID(c *gin.Context) {
	workdayID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workday ID"})
		return
	}

	workday, err := h.workdayService.GetEmployeeWorkdayByID(c.Request.Context(), uint(workdayID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if workday == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee workday not found"})
		return
	}

	c.JSON(http.StatusOK, workday)
}

// GetEmployeeWorkdaysByEmployeeID retrieves all workdays for a specific employee.
func (h *EmployeeWorkdayHandler) GetEmployeeWorkdaysByEmployeeID(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	workdays, err := h.workdayService.GetEmployeeWorkdaysByEmployeeID(c.Request.Context(), uint(employeeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, workdays)
}

// UpdateEmployeeWorkday handles updating an employee workday by its ID.
func (h *EmployeeWorkdayHandler) UpdateEmployeeWorkday(c *gin.Context) {
	workdayID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workday ID"})
		return
	}

	var workday models.EmployeeWorkday
	if err := c.ShouldBindJSON(&workday); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	workday.ID = uint(workdayID)
	success, err := h.workdayService.UpdateEmployeeWorkday(c.Request.Context(), workday)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee workday not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee workday updated successfully"})
}

// DeleteEmployeeWorkday handles deleting an employee workday by its ID.
func (h *EmployeeWorkdayHandler) DeleteEmployeeWorkday(c *gin.Context) {
	workdayID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid workday ID"})
		return
	}

	success, err := h.workdayService.DeleteEmployeeWorkday(c.Request.Context(), uint(workdayID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee workday not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee workday deleted successfully"})
}
