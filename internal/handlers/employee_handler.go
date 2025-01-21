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
	var request struct {
		RegistrationNumber string `json:"registration_number"`
		Qualification      string `json:"qualification"`
		CompanyID          uint   `json:"company_id"`
		StartHour          string `json:"start_hour"`
		EndHour            string `json:"end_hour"`
		FirstName          string `json:"first_name"`
		LastName           string `json:"last_name"`
		Username           string `json:"username"`
		Password           string `json:"password"`
	}

	// Bind the JSON payload into the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Map the request to the Employee and User models
	employee := models.Employee{
		RegistrationNumber: request.RegistrationNumber,
		Qualification:      request.Qualification,
		CompanyID:          request.CompanyID,
		StartHour:          request.StartHour,
		EndHour:            request.EndHour,
	}

	user := models.User{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Username:  request.Username,
		Password:  request.Password,
	}

	// Call the service to create the employee
	employeeID, err := h.employeeService.CreateEmployee(c.Request.Context(), employee, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("CREATE_EMPLOYEE")
	// Respond with success
	c.JSON(http.StatusCreated, gin.H{
		"id":      employeeID,
		"data":    employee,
		"message": "Employee created successfully",
	})
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

// UpdateEmployee handles updating an existing employee.
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	var request struct {
		ID                 uint   `json:"id"`
		RegistrationNumber string `json:"registration_number"`
		Qualification      string `json:"qualification"`
		CompanyID          uint   `json:"company_id"`
		StartHour          string `json:"start_hour"`
		EndHour            string `json:"end_hour"`
		FirstName          string `json:"first_name"`
		LastName           string `json:"last_name"`
		Username           string `json:"username"`
		Password           string `json:"password"`
	}

	// Bind the JSON payload into the request struct
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Map the request to the Employee and User models
	employee := models.Employee{
		ID:                 request.ID,
		RegistrationNumber: request.RegistrationNumber,
		Qualification:      request.Qualification,
		CompanyID:          request.CompanyID,
		StartHour:          request.StartHour,
		EndHour:            request.EndHour,
	}

	var userUpdates *models.User
	if request.FirstName != "" || request.LastName != "" || request.Username != "" || request.Password != "" {
		userUpdates = &models.User{
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Username:  request.Username,
			Password:  request.Password,
		}
	}

	// Call the service to update the employee
	success, err := h.employeeService.UpdateEmployee(c.Request.Context(), employee, userUpdates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	manager.broadcast <- []byte("UPDATE_EMPLOYEE")
	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message": "Employee updated successfully",
	})
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

	manager.broadcast <- []byte("DELETE_EMPLOYEE")
	c.JSON(http.StatusOK, gin.H{"message": "Employee deleted successfully"})
}

// FetchEmployees retrieves all employees with optional filters, pagination, and search.
func (h *EmployeeHandler) FetchEmployees(c *gin.Context) {
	// Extract query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	filters := map[string]interface{}{}

	// Add search filter if provided
	if search != "" {
		filters["search"] = search
	}

	// Add other filters from query parameters
	for key, values := range c.Request.URL.Query() {
		if key != "page" && key != "limit" && key != "search" {
			filters[key] = values[0]
		}
	}

	// Call the service to get paginated and filtered results
	employees, total, err := h.employeeService.FetchEmployees(c.Request.Context(), page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the response with pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"data":  employees,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}
