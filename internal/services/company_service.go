package services

import (
	"context"
	"errors"
	"fmt"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
)

// CompanyService defines the interface for company-related operations.
type CompanyService interface {
	CreateCompany(ctx context.Context, company models.Company) (uint, error)
	GetCompanyByID(ctx context.Context, id uint) (*models.Company, error)
	ListCompanies(ctx context.Context, page, limit int, filters map[string]interface{}) ([]models.Company, int64, error)
	UpdateCompany(ctx context.Context, company models.Company) (bool, error)
	DeleteCompany(ctx context.Context, id uint) (bool, error)
	ListCompaniesForSelect(ctx context.Context) ([]map[string]interface{}, error)
}

// companyService implements the CompanyService interface.
type companyService struct {
	companyRepo repositories.CompanyRepository
}

// NewCompanyService creates a new instance of CompanyService.
func NewCompanyService(companyRepo repositories.CompanyRepository) CompanyService {
	return &companyService{
		companyRepo: companyRepo,
	}
}

// CreateCompany creates a new company in the database.
func (s *companyService) CreateCompany(ctx context.Context, company models.Company) (uint, error) {
	// Validate the company name
	if company.CompanyName == "" {
		return 0, errors.New("company name is required")
	}

	// Check if the company name already exists
	existingCompany, err := s.companyRepo.GetCompanyByName(ctx, company.CompanyName)
	if err != nil {
		return 0, fmt.Errorf("failed to check existing company: %w", err)
	}
	if existingCompany != nil {
		return 0, errors.New("company name already exists")
	}

	// Create the company in the database
	companyID, err := s.companyRepo.CreateCompany(ctx, company)
	if err != nil {
		return 0, fmt.Errorf("failed to create company: %w", err)
	}

	return companyID, nil
}

// GetCompanyByID retrieves a company by its ID.
func (s *companyService) GetCompanyByID(ctx context.Context, id uint) (*models.Company, error) {
	company, err := s.companyRepo.GetCompanyByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve company by ID: %w", err)
	}
	return company, nil
}

// ListCompanies retrieves all companies with pagination, filtering, and search.
func (s *companyService) ListCompanies(ctx context.Context, page, limit int, filters map[string]interface{}) ([]models.Company, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Call the repository to get paginated and filtered results
	companies, total, err := s.companyRepo.ListCompaniesWithFilters(ctx, page, limit, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list companies: %w", err)
	}

	return companies, total, nil
}

// UpdateCompany updates an existing company in the database.
func (s *companyService) UpdateCompany(ctx context.Context, company models.Company) (bool, error) {
	// Validate that the company ID is provided
	if company.ID == 0 {
		return false, errors.New("company ID is required")
	}

	companyDb, err := s.companyRepo.GetCompanyByID(ctx, company.ID)
	if err != nil {
		return false, fmt.Errorf("failed to update company: %w", err)
	}
	companyDb.CompanyName = company.CompanyName
	// Update the company in the database
	success, err := s.companyRepo.UpdateCompany(ctx, *companyDb)
	if err != nil {
		return false, fmt.Errorf("failed to update company: %w", err)
	}

	return success, nil
}

// DeleteCompany deletes a company by its ID.
func (s *companyService) DeleteCompany(ctx context.Context, id uint) (bool, error) {
	// Validate that the company ID is provided
	if id == 0 {
		return false, errors.New("company ID is required")
	}

	// Delete the company from the database
	success, err := s.companyRepo.DeleteCompany(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete company: %w", err)
	}

	return success, nil
}

// ListCompaniesForSelect retrieves all companies for use in select options.
func (s *companyService) ListCompaniesForSelect(ctx context.Context) ([]map[string]interface{}, error) {
	companies, err := s.companyRepo.ListCompanies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}

	// Simplify the response for select options
	var result []map[string]interface{}
	for _, company := range companies {
		result = append(result, map[string]interface{}{
			"id":          company.ID,
			"companyName": company.CompanyName,
		})
	}

	return result, nil
}
