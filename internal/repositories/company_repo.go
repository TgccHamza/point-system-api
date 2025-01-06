package repository

import (
	"errors"
	"point-system-api/internal/models"

	"gorm.io/gorm"
)

// CreateCompany adds a new company to the database.
func CreateCompany(db *gorm.DB, company models.Company) error {
	// Create the company in the database
	if err := db.Create(&company).Error; err != nil {
		return err
	}
	return nil
}

// GetCompanyByID retrieves a company from the database by its ID.
func GetCompanyByID(db *gorm.DB, id uint) (models.Company, error) {
	var company models.Company
	// Retrieve the company by ID
	if err := db.First(&company, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Company{}, errors.New("company not found")
		}
		return models.Company{}, err
	}
	return company, nil
}

// ListCompanies retrieves all companies from the database.
func ListCompanies(db *gorm.DB) ([]models.Company, error) {
	var companies []models.Company
	// Retrieve all companies
	if err := db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

// UpdateCompany updates the details of an existing company in the database.
func UpdateCompany(db *gorm.DB, company models.Company) error {
	// Update the company in the database
	if err := db.Save(&company).Error; err != nil {
		return err
	}
	return nil
}

// DeleteCompany removes a company from the database by its ID.
func DeleteCompany(db *gorm.DB, id uint) error {
	var company models.Company
	// Find the company by ID and delete it
	if err := db.Delete(&company, id).Error; err != nil {
		return err
	}
	return nil
}
