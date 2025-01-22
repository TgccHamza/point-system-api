package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
	"point-system-api/internal/types"
)

// EmployeeRepository defines the interface for employee-related database operations.
type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, employee *models.Employee) error
	GetEmployeeByIDWithUser(ctx context.Context, id uint) (*types.EmployeeWithUser, error)
	GetEmployeeByID(ctx context.Context, id uint) (*models.Employee, error)
	GetEmployeesByCompanyID(ctx context.Context, companyID uint) ([]*models.Employee, error)
	UpdateEmployee(ctx context.Context, employee *models.Employee) error
	DeleteEmployee(ctx context.Context, id uint) error
	FetchEmployees(ctx context.Context) ([]*models.Employee, error)
	ListEmployeesWithFilters(ctx context.Context, limit, page int, filters map[string]interface{}, search string) ([]*types.EmployeeWithUser, int64, error) // New method
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

// GetEmployeeByID retrieves an employee by their ID.
func (r *employeeRepository) GetEmployeeByIDWithUser(ctx context.Context, id uint) (*types.EmployeeWithUser, error) {
	var employeeWithUser types.EmployeeWithUser

	// Use raw SQL to join Employee and User tables
	query := `
		SELECT 
			e.id, e.user_id, e.registration_number, e.qualification, e.company_id, 
			e.start_hour, e.end_hour, e.created_at, e.updated_at,
			u.id AS user_id, u.first_name, u.last_name, u.username, u.role
		FROM 
			employees e
		JOIN 
			users u ON e.user_id = u.id
		WHERE 
			e.id = ? && e.deleted_at is null
	`

	// Execute the query
	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&employeeWithUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No employee found
		}
		return nil, fmt.Errorf("failed to retrieve employee by ID: %w", err)
	}

	return &employeeWithUser, nil
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

func (r *employeeRepository) ListEmployeesWithFilters(ctx context.Context, limit, page int, filters map[string]interface{}, search string) ([]*types.EmployeeWithUser, int64, error) {
	var employeesWithUsers []*types.EmployeeWithUser
	var totalCount int64
	offset := (page - 1) * limit

	// Base query
	query := r.db.WithContext(ctx).Table("employees e").
		Select(`
			e.id, e.user_id, e.registration_number, e.qualification, e.company_id, 
			e.start_hour, e.end_hour, e.created_at, e.updated_at,
			u.id AS user_id, u.first_name, u.last_name, u.username, u.role
		`).
		Joins("JOIN users u ON e.user_id = u.id")

	// Apply filters
	for key, value := range filters {
		if key == "search" {
			continue
		}
		query = query.Where(fmt.Sprintf("e.%s = ?", key), value)
	}

	query = query.Where("e.deleted_at is null")

	// Apply search (if provided)
	if search != "" {
		query = query.Where(`
			e.registration_number LIKE ? OR 
			e.qualification LIKE ? OR 
			u.first_name LIKE ? OR 
			u.last_name LIKE ? OR 
			u.username LIKE ?
		`, "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total records (for pagination)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}

	// Apply pagination
	query = query.Offset(offset).Limit(limit)

	// Execute the query
	if err := query.Scan(&employeesWithUsers).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch employees: %w", err)
	}

	return employeesWithUsers, totalCount, nil
}
