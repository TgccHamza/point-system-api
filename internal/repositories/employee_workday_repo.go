package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
)

// EmployeeWorkDayRepository defines the interface for employee workday-related database operations.
type EmployeeWorkDayRepository interface {
	CreateEmployeeWorkDay(ctx context.Context, employeeWorkDay *models.EmployeeWorkday) error
	GetEmployeeWorkDayByID(ctx context.Context, id uint) (*models.EmployeeWorkday, error)
	UpdateEmployeeWorkDay(ctx context.Context, employeeWorkDay *models.EmployeeWorkday) error
}

// employeeWorkDayRepository implements the EmployeeWorkDayRepository interface.
type employeeWorkDayRepository struct {
	db *gorm.DB
}

// NewEmployeeWorkDayRepository creates a new instance of EmployeeWorkDayRepository.
func NewEmployeeWorkDayRepository(db *gorm.DB) EmployeeWorkDayRepository {
	return &employeeWorkDayRepository{
		db: db,
	}
}

// CreateEmployeeWorkDay inserts a new employee workday record into the database.
func (r *employeeWorkDayRepository) CreateEmployeeWorkDay(ctx context.Context, employeeWorkDay *models.EmployeeWorkday) error {
	if employeeWorkDay == nil {
		return errors.New("employee workday is nil")
	}

	if employeeWorkDay.WorkDayID == 0 {
		return errors.New("work day ID is required")
	}

	if employeeWorkDay.EmployeeID == 0 {
		return errors.New("employee ID is required")
	}

	if err := r.db.WithContext(ctx).Create(employeeWorkDay).Error; err != nil {
		return fmt.Errorf("failed to create employee workday: %w", err)
	}

	return nil
}

// GetEmployeeWorkDayByID retrieves an employee workday record by its ID.
func (r *employeeWorkDayRepository) GetEmployeeWorkDayByID(ctx context.Context, id uint) (*models.EmployeeWorkday, error) {
	if id == 0 {
		return nil, errors.New("invalid employee workday ID")
	}

	var employeeWorkDay models.EmployeeWorkday
	if err := r.db.WithContext(ctx).First(&employeeWorkDay, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No employee workday found
		}
		return nil, fmt.Errorf("failed to retrieve employee workday by ID: %w", err)
	}
	return &employeeWorkDay, nil
}

// UpdateEmployeeWorkDay updates an existing employee workday record in the database.
func (r *employeeWorkDayRepository) UpdateEmployeeWorkDay(ctx context.Context, employeeWorkDay *models.EmployeeWorkday) error {
	if employeeWorkDay == nil || employeeWorkDay.ID == 0 {
		return errors.New("invalid employee workday data")
	}

	if employeeWorkDay.WorkDayID == 0 {
		return errors.New("work day ID is required")
	}

	if employeeWorkDay.EmployeeID == 0 {
		return errors.New("employee ID is required")
	}

	if err := r.db.WithContext(ctx).Save(employeeWorkDay).Error; err != nil {
		return fmt.Errorf("failed to update employee workday: %w", err)
	}

	return nil
}