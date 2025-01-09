package services

import (
	"context"
	"errors"
	"fmt"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
)

// EmployeeService defines the interface for employee-related operations.
type EmployeeService interface {
	CreateEmployee(ctx context.Context, employee models.Employee) (uint, error)
	GetEmployeeByID(ctx context.Context, id uint) (*models.Employee, error)
	GetEmployeesByCompanyID(ctx context.Context, companyID uint) ([]*models.Employee, error)
	UpdateEmployee(ctx context.Context, employee models.Employee) (bool, error)
	DeleteEmployee(ctx context.Context, id uint) (bool, error)
}

// employeeService implements the EmployeeService interface.
type employeeService struct {
	employeeRepo repositories.EmployeeRepository
}

// NewEmployeeService creates a new instance of EmployeeService.
func NewEmployeeService(employeeRepo repositories.EmployeeRepository) EmployeeService {
	return &employeeService{
		employeeRepo: employeeRepo,
	}
}

// CreateEmployee creates a new employee in the database.
func (s *employeeService) CreateEmployee(ctx context.Context, employee models.Employee) (uint, error) {
	// Validate that the associated user and company exist
	if employee.UserID == 0 {
		return 0, errors.New("user ID is required")
	}
	if employee.CompanyID == 0 {
		return 0, errors.New("company ID is required")
	}

	// Validate that the start and end hours are provided
	if employee.StartHour.IsZero() {
		return 0, errors.New("start hour is required")
	}
	if employee.EndHour.IsZero() {
		return 0, errors.New("end hour is required")
	}

	// Ensure the end hour is after the start hour
	if !employee.EndHour.After(employee.StartHour) {
		return 0, errors.New("end hour must be after start hour")
	}

	// Create the employee in the database
	err := s.employeeRepo.CreateEmployee(ctx, &employee)
	if err != nil {
		return 0, fmt.Errorf("failed to create employee: %w", err)
	}

	return employee.ID, nil
}

// GetEmployeeByID retrieves an employee by their ID.
func (s *employeeService) GetEmployeeByID(ctx context.Context, id uint) (*models.Employee, error) {
	employee, err := s.employeeRepo.GetEmployeeByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employee by ID: %w", err)
	}
	return employee, nil
}

// GetEmployeesByCompanyID retrieves all employees for a specific company.
func (s *employeeService) GetEmployeesByCompanyID(ctx context.Context, companyID uint) ([]*models.Employee, error) {
	employees, err := s.employeeRepo.GetEmployeesByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve employees by company ID: %w", err)
	}
	return employees, nil
}

// UpdateEmployee updates an existing employee in the database.
func (s *employeeService) UpdateEmployee(ctx context.Context, employee models.Employee) (bool, error) {
	// Validate that the employee ID is provided
	if employee.ID == 0 {
		return false, errors.New("employee ID is required")
	}

	// Validate that the start and end hours are provided
	if employee.StartHour.IsZero() {
		return false, errors.New("start hour is required")
	}
	if employee.EndHour.IsZero() {
		return false, errors.New("end hour is required")
	}

	// Ensure the end hour is after the start hour
	if !employee.EndHour.After(employee.StartHour) {
		return false, errors.New("end hour must be after start hour")
	}

	// Update the employee in the database
	err := s.employeeRepo.UpdateEmployee(ctx, &employee)
	if err != nil {
		return false, fmt.Errorf("failed to update employee: %w", err)
	}

	return true, nil
}

// DeleteEmployee deletes an employee by their ID.
func (s *employeeService) DeleteEmployee(ctx context.Context, id uint) (bool, error) {
	// Validate that the employee ID is provided
	if id == 0 {
		return false, errors.New("employee ID is required")
	}

	// Delete the employee from the database
	err := s.employeeRepo.DeleteEmployee(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete employee: %w", err)
	}

	return true, nil
}