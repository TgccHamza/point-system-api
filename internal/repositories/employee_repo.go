package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
)

func CreateEmployee(ctx context.Context, db *gorm.DB, employee *models.Employee) error {
	if err := db.WithContext(ctx).Create(employee).Error; err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}
	return nil
}

func GetEmployeeByID(ctx context.Context, db *gorm.DB, id uint) (*models.Employee, error) {
	employee := &models.Employee{}
	if err := db.WithContext(ctx).First(employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("employee not found: %w", err)
		}
		return nil, fmt.Errorf("failed to retrieve employee: %w", err)
	}
	return employee, nil
}

func GetEmployeesByCompanyID(ctx context.Context, db *gorm.DB, companyID uint) ([]*models.Employee, error) {
	var employees []*models.Employee
	if err := db.WithContext(ctx).Where("company_id = ?", companyID).Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve employees for company %d: %w", companyID, err)
	}
	return employees, nil
}

func UpdateEmployee(ctx context.Context, db *gorm.DB, employee *models.Employee) error {
	if err := db.WithContext(ctx).Save(employee).Error; err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}
	return nil
}

func DeleteEmployee(ctx context.Context, db *gorm.DB, id uint) error {
	if err := db.WithContext(ctx).Delete(&models.Employee{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}
	return nil
}
