package services

import (
	"context"
	"errors"
	"fmt"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
)

// EmployeeWorkdayService defines the interface for employee workday-related operations.
type EmployeeWorkdayService interface {
	CreateEmployeeWorkday(ctx context.Context, workday models.EmployeeWorkday) (uint, error)
	GetEmployeeWorkdayByID(ctx context.Context, id uint) (*models.EmployeeWorkday, error)
	GetEmployeeWorkdaysByEmployeeID(ctx context.Context, employeeID uint) ([]*models.EmployeeWorkday, error)
	UpdateEmployeeWorkday(ctx context.Context, workday models.EmployeeWorkday) (bool, error)
	DeleteEmployeeWorkday(ctx context.Context, id uint) (bool, error)
}

// employeeWorkdayService implements the EmployeeWorkdayService interface.
type employeeWorkdayService struct {
	workdayRepo repositories.EmployeeWorkdayRepository
}

// NewEmployeeWorkdayService creates a new instance of EmployeeWorkdayService.
func NewEmployeeWorkdayService(workdayRepo repositories.EmployeeWorkdayRepository) EmployeeWorkdayService {
	return &employeeWorkdayService{
		workdayRepo: workdayRepo,
	}
}

// CreateEmployeeWorkday creates a new employee workday in the database.
func (s *employeeWorkdayService) CreateEmployeeWorkday(ctx context.Context, workday models.EmployeeWorkday) (uint, error) {
	// Validate that the associated employee exists
	if workday.EmployeeID == 0 {
		return 0, errors.New("employee ID is required")
	}

	// Create the employee workday in the database
	err := s.workdayRepo.CreateEmployeeWorkday(ctx, &workday)
	if err != nil {
		return 0, fmt.Errorf("failed to create employee workday: %w", err)
	}

	return workday.ID, nil
}

// GetEmployeeWorkdayByID retrieves an employee workday by its ID.
func (s *employeeWorkdayService) GetEmployeeWorkdayByID(ctx context.Context, id uint) (*models.EmployeeWorkday, error) {
	workday, err := s.workdayRepo.GetEmployeeWorkdayByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employee workday by ID: %w", err)
	}
	return workday, nil
}

// GetEmployeeWorkdaysByEmployeeID retrieves all workdays for a specific employee.
func (s *employeeWorkdayService) GetEmployeeWorkdaysByEmployeeID(ctx context.Context, employeeID uint) ([]*models.EmployeeWorkday, error) {
	workdays, err := s.workdayRepo.GetEmployeeWorkdaysByEmployeeID(ctx, employeeID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employee workdays by employee ID: %w", err)
	}
	return workdays, nil
}

// UpdateEmployeeWorkday updates an existing employee workday in the database.
func (s *employeeWorkdayService) UpdateEmployeeWorkday(ctx context.Context, workday models.EmployeeWorkday) (bool, error) {
	// Validate that the workday ID is provided
	if workday.ID == 0 {
		return false, errors.New("workday ID is required")
	}

	// Update the employee workday in the database
	err := s.workdayRepo.UpdateEmployeeWorkday(ctx, &workday)
	if err != nil {
		return false, fmt.Errorf("failed to update employee workday: %w", err)
	}

	return true, nil
}

// DeleteEmployeeWorkday deletes an employee workday by its ID.
func (s *employeeWorkdayService) DeleteEmployeeWorkday(ctx context.Context, id uint) (bool, error) {
	// Validate that the workday ID is provided
	if id == 0 {
		return false, errors.New("workday ID is required")
	}

	// Delete the employee workday from the database
	err := s.workdayRepo.DeleteEmployeeWorkday(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete employee workday: %w", err)
	}

	return true, nil
}
