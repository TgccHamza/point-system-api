package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
)

// EmployeeRepository defines the interface for employee-related database operations.
type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, employee *models.Employee) error
	GetEmployeeByID(ctx context.Context, id uint) (*models.Employee, error)
	GetEmployeesByCompanyID(ctx context.Context, companyID uint) ([]*models.Employee, error)
	UpdateEmployee(ctx context.Context, employee *models.Employee) error
	DeleteEmployee(ctx context.Context, id uint) error
	FetchEmployees(ctx context.Context) ([]*models.Employee, error)
	ListEmployeesWithFilters(ctx context.Context, page, limit int, filters map[string]interface{}) ([]*models.Employee, int64, error) // New method
}

// employeeRepository implements the EmployeeRepository interface.
type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository creates a new instance of EmployeeRepository.
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{
		db: db,
	}
}

// CreateEmployee inserts a new employee into the database.
func (r *employeeRepository) CreateEmployee(ctx context.Context, employee *models.Employee) error {
	// Validate that the associated user and company exist
	if employee.UserID == 0 {
		return errors.New("user ID is required")
	}
	if employee.CompanyID == 0 {
		return errors.New("company ID is required")
	}

	// Create the employee in the database
	if err := r.db.WithContext(ctx).Create(employee).Error; err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}

	return nil
}

// GetEmployeeByID retrieves an employee by their ID.
func (r *employeeRepository) GetEmployeeByID(ctx context.Context, id uint) (*models.Employee, error) {
	var employee models.Employee
	if err := r.db.WithContext(ctx).First(&employee, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No employee found
		}
		return nil, fmt.Errorf("failed to retrieve employee by ID: %w", err)
	}
	return &employee, nil
}

// GetEmployeesByCompanyID retrieves all employees for a specific company.
func (r *employeeRepository) GetEmployeesByCompanyID(ctx context.Context, companyID uint) ([]*models.Employee, error) {
	var employees []*models.Employee
	if err := r.db.WithContext(ctx).Where("company_id = ?", companyID).Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve employees by company ID: %w", err)
	}
	return employees, nil
}

// UpdateEmployee updates an existing employee in the database.
func (r *employeeRepository) UpdateEmployee(ctx context.Context, employee *models.Employee) error {
	// Validate that the employee ID is provided
	if employee.ID == 0 {
		return errors.New("employee ID is required")
	}

	// Update the employee in the database
	if err := r.db.WithContext(ctx).Save(employee).Error; err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	return nil
}

// DeleteEmployee deletes an employee by their ID.
func (r *employeeRepository) DeleteEmployee(ctx context.Context, id uint) error {
	// Validate that the employee ID is provided
	if id == 0 {
		return errors.New("employee ID is required")
	}

	// Delete the employee from the database
	if err := r.db.WithContext(ctx).Delete(&models.Employee{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	return nil
}

// FetchEmployees retrieves all employees from the database.
func (r *employeeRepository) FetchEmployees(ctx context.Context) ([]*models.Employee, error) {
	var employees []*models.Employee
	if err := r.db.WithContext(ctx).Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch employees: %w", err)
	}
	return employees, nil
}

// ListEmployeesWithFilters retrieves employees with pagination, filtering, and search.
func (r *employeeRepository) ListEmployeesWithFilters(ctx context.Context, page, limit int, filters map[string]interface{}) ([]*models.Employee, int64, error) {
	offset := (page - 1) * limit

	// Build the query
	query := r.db.Model(&models.Employee{})

	// Apply search filter
	if search, ok := filters["search"]; ok {
		query = query.Where("registration_number LIKE ? OR qualification LIKE ?", "%"+search.(string)+"%", "%"+search.(string)+"%")
	}

	// Apply other filters
	for key, value := range filters {
		if key != "search" {
			query = query.Where(key+" = ?", value)
		}
	}

	// Count total records
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}

	// Apply pagination
	var employees []*models.Employee
	if err := query.Offset(offset).Limit(limit).Find(&employees).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list employees: %w", err)
	}

	return employees, total, nil
}
