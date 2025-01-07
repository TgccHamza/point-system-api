package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
)

// CompanyRepository defines the interface for company-related database operations.
type CompanyRepository interface {
	CreateCompany(ctx context.Context, company models.Company) (uint, error)
	GetCompanyByID(ctx context.Context, id uint) (*models.Company, error)
	GetCompanyByName(ctx context.Context, name string) (*models.Company, error)
	ListCompanies(ctx context.Context) ([]models.Company, error)
	UpdateCompany(ctx context.Context, company models.Company) (bool, error)
	DeleteCompany(ctx context.Context, id uint) (bool, error)
}

// companyRepository implements the CompanyRepository interface.
type companyRepository struct {
	db *gorm.DB
}

// NewCompanyRepository creates a new instance of CompanyRepository.
func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{
		db: db,
	}
}

// CreateCompany inserts a new company into the database.
func (r *companyRepository) CreateCompany(ctx context.Context, company models.Company) (uint, error) {
	// Validate the company name
	if company.CompanyName == "" {
		return 0, errors.New("company name is required")
	}

	// Create the company in the database
	if err := r.db.WithContext(ctx).Create(&company).Error; err != nil {
		return 0, fmt.Errorf("failed to create company: %w", err)
	}

	// Return the ID of the newly created company
	return company.ID, nil
}

// GetCompanyByID retrieves a company by its ID.
func (r *companyRepository) GetCompanyByID(ctx context.Context, id uint) (*models.Company, error) {
	var company models.Company
	if err := r.db.WithContext(ctx).First(&company, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No company found
		}
		return nil, fmt.Errorf("failed to retrieve company by ID: %w", err)
	}
	return &company, nil
}

// GetCompanyByName retrieves a company by its name.
func (r *companyRepository) GetCompanyByName(ctx context.Context, name string) (*models.Company, error) {
	var company models.Company
	if err := r.db.WithContext(ctx).Where("company_name = ?", name).First(&company).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No company found
		}
		return nil, fmt.Errorf("failed to retrieve company by name: %w", err)
	}
	return &company, nil
}

// ListCompanies retrieves all companies from the database.
func (r *companyRepository) ListCompanies(ctx context.Context) ([]models.Company, error) {
	var companies []models.Company
	if err := r.db.WithContext(ctx).Find(&companies).Error; err != nil {
		return nil, fmt.Errorf("failed to list companies: %w", err)
	}
	return companies, nil
}

// UpdateCompany updates an existing company in the database.
func (r *companyRepository) UpdateCompany(ctx context.Context, company models.Company) (bool, error) {
	// Validate that the company ID is provided
	if company.ID == 0 {
		return false, errors.New("company ID is required")
	}

	// Update the company in the database
	if err := r.db.WithContext(ctx).Save(&company).Error; err != nil {
		return false, fmt.Errorf("failed to update company: %w", err)
	}

	return true, nil
}

// DeleteCompany deletes a company by its ID.
func (r *companyRepository) DeleteCompany(ctx context.Context, id uint) (bool, error) {
	// Validate that the company ID is provided
	if id == 0 {
		return false, errors.New("company ID is required")
	}

	// Delete the company from the database
	if err := r.db.WithContext(ctx).Delete(&models.Company{}, id).Error; err != nil {
		return false, fmt.Errorf("failed to delete company: %w", err)
	}

	return true, nil
}
