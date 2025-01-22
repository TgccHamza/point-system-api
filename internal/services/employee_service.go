package services

import (
	"context"
	"errors"
	"fmt"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
	"point-system-api/internal/types"
	"point-system-api/pkg/utils"
)

// EmployeeService defines the interface for employee-related operations.
type EmployeeService interface {
	CreateEmployee(ctx context.Context, employee models.Employee, user models.User) (*models.Employee, error)
	GetEmployeeByID(ctx context.Context, id uint) (*types.EmployeeWithUser, error)
	GetEmployeesByCompanyID(ctx context.Context, companyID uint) ([]*models.Employee, error)
	UpdateEmployee(ctx context.Context, employee models.Employee, userUpdates *models.User) (bool, error)
	DeleteEmployee(ctx context.Context, id uint) (bool, error)
	FetchEmployees(ctx context.Context, page, limit int, filters map[string]interface{}, search string) ([]*types.EmployeeWithUser, int64, error)
}

// employeeService implements the EmployeeService interface.
type employeeService struct {
	employeeRepo repositories.EmployeeRepository
	userService  UserService
}

// NewEmployeeService creates a new instance of EmployeeService.
func NewEmployeeService(employeeRepo repositories.EmployeeRepository, userService UserService) EmployeeService {
	return &employeeService{
		employeeRepo: employeeRepo,
		userService:  userService,
	}
}

// CreateEmployee creates a new employee and a corresponding user.
func (s *employeeService) CreateEmployee(ctx context.Context, employee models.Employee, user models.User) (*models.Employee, error) {
	// Validate that the required fields are provided
	if user.FirstName == "" || user.LastName == "" || user.Username == "" || user.Password == "" {
		return nil, errors.New("first name, last name, username, and password are required")
	}

	// Set the default role for the user
	user.Role = "employee" // Default role for employees

	// Create the user first
	userID, err := s.userService.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Set the UserID for the employee
	employee.UserID = userID

	// Create the employee in the database
	err = s.employeeRepo.CreateEmployee(ctx, &employee)
	if err != nil {
		// Rollback: Delete the user if employee creation fails
		_, _ = s.userService.DeleteUser(ctx, userID)
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	return &employee, nil
}

// GetEmployeeByID retrieves an employee by their ID.
func (s *employeeService) GetEmployeeByID(ctx context.Context, id uint) (*types.EmployeeWithUser, error) {
	employee, err := s.employeeRepo.GetEmployeeByIDWithUser(ctx, id)
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

// UpdateEmployee updates an existing employee and their associated user.
func (s *employeeService) UpdateEmployee(ctx context.Context, employee models.Employee, userUpdates *models.User) (bool, error) {
	// Validate that the employee ID is provided
	if employee.ID == 0 {
		return false, errors.New("employee ID is required")
	}

	// Retrieve the existing employee
	existingEmployee, err := s.employeeRepo.GetEmployeeByID(ctx, employee.ID)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve employee: %w", err)
	}

	if existingEmployee == nil {
		return false, errors.New("employee not found")
	}

	// Update the employee details
	if employee.RegistrationNumber != "" {
		existingEmployee.RegistrationNumber = employee.RegistrationNumber
	}

	if employee.Qualification != "" {
		existingEmployee.Qualification = employee.Qualification
	}

	if employee.CompanyID != 0 {
		existingEmployee.CompanyID = employee.CompanyID
	}

	if employee.StartHour != "" {
		existingEmployee.StartHour = employee.StartHour
	}

	if employee.EndHour != "" {
		existingEmployee.EndHour = employee.EndHour
	}

	// Save the updated employee
	err = s.employeeRepo.UpdateEmployee(ctx, existingEmployee)
	if err != nil {
		return false, fmt.Errorf("failed to update employee: %w", err)
	}

	// Update the associated user (if userUpdates is provided)
	if userUpdates != nil {
		// Retrieve the existing user
		existingUser, err := s.userService.GetUserByID(ctx, existingEmployee.UserID)
		if err != nil {
			return false, fmt.Errorf("failed to retrieve user: %w", err)
		}
		if existingUser == nil {
			return false, errors.New("user not found")
		}

		// Update the user details
		if userUpdates.FirstName != "" {
			existingUser.FirstName = userUpdates.FirstName
		}
		if userUpdates.LastName != "" {
			existingUser.LastName = userUpdates.LastName
		}
		if userUpdates.Username != "" {
			existingUser.Username = userUpdates.Username
		}
		if userUpdates.Password != "" {
			// Hash the new password before saving
			hashedPassword, err := utils.HashPassword(userUpdates.Password)
			if err != nil {
				return false, fmt.Errorf("failed to hash password: %w", err)
			}
			existingUser.Password = hashedPassword
		}

		// Save the updated user
		_, err = s.userService.UpdateUser(ctx, *existingUser)
		if err != nil {
			return false, fmt.Errorf("failed to update user: %w", err)
		}
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

// FetchEmployees retrieves all employees with pagination, filtering, and search.
func (s *employeeService) FetchEmployees(ctx context.Context, page, limit int, filters map[string]interface{}, search string) ([]*types.EmployeeWithUser, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Call the repository to get paginated and filtered results
	employees, total, err := s.employeeRepo.ListEmployeesWithFilters(ctx, limit, page, filters, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch employees: %w", err)
	}

	return employees, total, nil
}
