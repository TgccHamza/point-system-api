package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
	"point-system-api/pkg/utils"
)

// UserRepository defines the interface for user-related database operations.
type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) (uint, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	ListUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, user models.User) (bool, error)
	DeleteUser(ctx context.Context, id uint) (bool, error)
	ListUsersWithFilters(ctx context.Context, page, limit int, filters map[string]interface{}) ([]models.User, int64, error)
}

// userRepository implements the UserRepository interface.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// CreateUser inserts a new user into the database after hashing the password.
func (r *userRepository) CreateUser(ctx context.Context, user models.User) (uint, error) {
	// Hash the user's password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Create the user in the database
	if err := r.db.WithContext(ctx).Create(&user).Error; err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	// Return the ID of the newly created user
	return user.ID, nil
}

// GetUserByID retrieves a user by their ID.
func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("failed to retrieve user by ID: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("failed to retrieve user by username: %w", err)
	}
	return &user, nil
}

// ListUsers retrieves all users from the database.
func (r *userRepository) ListUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.db.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

// UpdateUser updates an existing user in the database.
func (r *userRepository) UpdateUser(ctx context.Context, user models.User) (bool, error) {
	// Hash the user's password before updating
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return false, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	// Update the user in the database
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return false, fmt.Errorf("failed to update user: %w", err)
	}

	return true, nil
}

// DeleteUser deletes a user by their ID.
func (r *userRepository) DeleteUser(ctx context.Context, id uint) (bool, error) {
	if err := r.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}
	return true, nil
}

// ListUsersWithFilters retrieves users with pagination, filtering, and search.
func (r *userRepository) ListUsersWithFilters(ctx context.Context, page, limit int, filters map[string]interface{}) ([]models.User, int64, error) {
	offset := (page - 1) * limit

	// Build the query
	query := r.db.Model(&models.User{})

	// Apply search filter
	if search, ok := filters["search"]; ok {
		query = query.Where("username LIKE ? OR first_name LIKE ? OR last_name LIKE ?", "%"+search.(string)+"%", "%"+search.(string)+"%", "%"+search.(string)+"%")
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
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Apply pagination
	var users []models.User
	if err := query.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}
