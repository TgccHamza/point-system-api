package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
)

// EmployeeWorkdayRepository defines the interface for employee workday-related database operations.
type EmployeeWorkdayRepository interface {
	CreateEmployeeWorkday(ctx context.Context, workday *models.EmployeeWorkday) error
	GetEmployeeWorkdayByID(ctx context.Context, id uint) (*models.EmployeeWorkday, error)
	GetEmployeeWorkdaysByEmployeeID(ctx context.Context, employeeID uint) ([]*models.EmployeeWorkday, error)
	UpdateEmployeeWorkday(ctx context.Context, workday *models.EmployeeWorkday) error
	DeleteEmployeeWorkday(ctx context.Context, id uint) error
}

// employeeWorkdayRepository implements the EmployeeWorkdayRepository interface.
type employeeWorkdayRepository struct {
	db *gorm.DB
}

// NewEmployeeWorkdayRepository creates a new instance of EmployeeWorkdayRepository.
func NewEmployeeWorkdayRepository(db *gorm.DB) EmployeeWorkdayRepository {
	return &employeeWorkdayRepository{
		db: db,
	}
}

// CreateEmployeeWorkday inserts a new employee workday into the database.
func (r *employeeWorkdayRepository) CreateEmployeeWorkday(ctx context.Context, workday *models.EmployeeWorkday) error {
	// Validate that the associated employee exists
	if workday.EmployeeID == 0 {
		return errors.New("employee ID is required")
	}

	// Create the employee workday in the database
	if err := r.db.WithContext(ctx).Create(workday).Error; err != nil {
		return fmt.Errorf("failed to create employee workday: %w", err)
	}

	return nil
}

// GetEmployeeWorkdayByID retrieves an employee workday by its ID.
func (r *employeeWorkdayRepository) GetEmployeeWorkdayByID(ctx context.Context, id uint) (*models.EmployeeWorkday, error) {
	var workday models.EmployeeWorkday
	if err := r.db.WithContext(ctx).First(&workday, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No workday found
		}
		return nil, fmt.Errorf("failed to retrieve employee workday by ID: %w", err)
	}
	return &workday, nil
}

// GetEmployeeWorkdaysByEmployeeID retrieves all workdays for a specific employee.
func (r *employeeWorkdayRepository) GetEmployeeWorkdaysByEmployeeID(ctx context.Context, employeeID uint) ([]*models.EmployeeWorkday, error) {
	var workdays []*models.EmployeeWorkday
	if err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Find(&workdays).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve employee workdays by employee ID: %w", err)
	}
	return workdays, nil
}

// UpdateEmployeeWorkday updates an existing employee workday in the database.
func (r *employeeWorkdayRepository) UpdateEmployeeWorkday(ctx context.Context, workday *models.EmployeeWorkday) error {
	// Validate that the workday ID is provided
	if workday.ID == 0 {
		return errors.New("workday ID is required")
	}

	// Update the employee workday in the database
	if err := r.db.WithContext(ctx).Save(workday).Error; err != nil {
		return fmt.Errorf("failed to update employee workday: %w", err)
	}

	return nil
}

// DeleteEmployeeWorkday deletes an employee workday by its ID.
func (r *employeeWorkdayRepository) DeleteEmployeeWorkday(ctx context.Context, id uint) error {
	// Validate that the workday ID is provided
	if id == 0 {
		return errors.New("workday ID is required")
	}

	// Delete the employee workday from the database
	if err := r.db.WithContext(ctx).Delete(&models.EmployeeWorkday{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete employee workday: %w", err)
	}

	return nil
}
