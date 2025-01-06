package handlers

import (
	"net/http"
	"point-system-api/internal/database"
	"point-system-api/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CompanyHandler handles HTTP requests for company-related operations.
type CompanyHandler struct {
	service database.Service
}

// NewCompanyHandler creates a new CompanyHandler instance.
func NewCompanyHandler(service database.Service) *CompanyHandler {
	return &CompanyHandler{service: service}
}

// CreateCompany handles the creation of a new company.
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	err := h.service.CreateCompany(company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Company created successfully"})
}

// GetCompany retrieves a company by its ID.
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	company, err := h.service.GetCompanyByID(companyID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}

	c.JSON(http.StatusOK, company)
}

// ListCompanies lists all companies.
func (h *CompanyHandler) ListCompanies(c *gin.Context) {
	companies, err := h.service.ListCompanies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve companies"})
		return
	}

	c.JSON(http.StatusOK, companies)
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

	company.ID = companyID
	err = h.service.UpdateCompany(company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company updated successfully"})
}

// DeleteCompany handles deleting a company by its ID.
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	companyID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	err = h.service.DeleteCompany(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company deleted successfully"})
}
