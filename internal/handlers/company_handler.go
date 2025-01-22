package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/models"
	"point-system-api/internal/services"
)

// CompanyHandler handles HTTP requests for company-related operations.
type CompanyHandler struct {
	companyService services.CompanyService
}

// NewCompanyHandler creates a new instance of CompanyHandler.
func NewCompanyHandler(companyService services.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

// CreateCompany handles the creation of a new company.
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	companyID, err := h.companyService.CreateCompany(c.Request.Context(), company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("CREATE_COMPANY")

	company_res, err := h.companyService.GetCompanyByID(c.Request.Context(), companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": companyID, "data": company_res, "message": "Company created successfully"})
}

// GetCompanyByID retrieves a company by its ID.
func (h *CompanyHandler) GetCompanyByID(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	company, err := h.companyService.GetCompanyByID(c.Request.Context(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if company == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

// ListCompanies retrieves all companies with optional filters, pagination, and search.
func (h *CompanyHandler) ListCompanies(c *gin.Context) {
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
	companies, total, err := h.companyService.ListCompanies(c.Request.Context(), page, limit, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the response with pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"data":  companies,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

// UpdateCompany handles updating a company by its ID.
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	company.ID = uint(companyID)
	success, err := h.companyService.UpdateCompany(c.Request.Context(), company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	company_res, err := h.companyService.GetCompanyByID(c.Request.Context(), company.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	manager.broadcast <- []byte("UPDATE_COMPANY")
	c.JSON(http.StatusOK, gin.H{"message": "Company updated successfully", "data": company_res})
}

// DeleteCompany handles deleting a company by its ID.
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	success, err := h.companyService.DeleteCompany(c.Request.Context(), uint(companyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !success {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	manager.broadcast <- []byte("DELETE_COMPANY")
	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}

// ListCompaniesForSelect retrieves all companies for use in select options.
func (h *CompanyHandler) ListCompaniesForSelect(c *gin.Context) {
	companies, err := h.companyService.ListCompaniesForSelect(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, companies)
}
