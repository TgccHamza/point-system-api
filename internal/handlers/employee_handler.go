package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/models"
	"point-system-api/internal/services"
)

// EmployeeHandler handles HTTP requests for employee-related operations.
type EmployeeHandler struct {
	employeeService services.EmployeeService
}

// NewEmployeeHandler creates a new instance of EmployeeHandler.
func NewEmployeeHandler(employeeService services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
	}
}

// CreateEmployee handles the creation of a new employee.
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Ensure the start and end hours are provided
	if employee.StartHour.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start hour is required"})
		return
	}
	if employee.EndHour.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end hour is required"})
		return
	}

	// Ensure the end hour is after the start hour
	if !employee.EndHour.After(employee.StartHour) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end hour must be after start hour"})
		return
	}

	employeeID, err := h.employeeService.CreateEmployee(c.Request.Context(), employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": employeeID, "message": "Employee created successfully"})
}

// GetEmployeeByID retrieves an employee by their ID.
func (h *EmployeeHandler) GetEmployeeByID(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.employeeService.GetEmployeeByID(c.Request.Context(), uint(employeeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if employee == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

// GetEmployeesByCompanyID retrieves all employees for a specific company.
func (h *EmployeeHandler) GetEmployeesByCompanyID(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	employees, err := h.employeeService.GetEmployeesByCompanyID(c.Request.Context(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, employees)
}

// UpdateEmployee handles updating an employee by their ID.
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Ensure the start and end hours are provided
	if employee.StartHour.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start hour is required"})
		return
	}
	if employee.EndHour.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end hour is required"})
		return
	}

	// Ensure the end hour is after the start hour
	if !employee.EndHour.After(employee.StartHour) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end hour must be after start hour"})
		return
	}

	employee.ID = uint(employeeID)
	success, err := h.employeeService.UpdateEmployee(c.Request.Context(), employee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee updated successfully"})
}

// DeleteEmployee handles deleting an employee by their ID.
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	employeeID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	success, err := h.employeeService.DeleteEmployee(c.Request.Context(), uint(employeeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}
