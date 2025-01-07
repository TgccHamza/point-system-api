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
	ListCompanies(ctx context.Context) ([]models.Company, error)
	UpdateCompany(ctx context.Context, company models.Company) (bool, error)
	DeleteCompany(ctx context.Context, id uint) (bool, error)
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

// ListCompanies retrieves all companies from the database.
func (s *companyService) ListCompanies(ctx context.Context) ([]models.Company, error) {
	companies, err := s.companyRepo.ListCompanies(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}
	return companies, nil
}

// UpdateCompany updates an existing company in the database.
func (s *companyService) UpdateCompany(ctx context.Context, company models.Company) (bool, error) {
	// Validate that the company ID is provided
	if company.ID == 0 {
		return false, errors.New("company ID is required")
	}

	// Update the company in the database
	success, err := s.companyRepo.UpdateCompany(ctx, company)
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
