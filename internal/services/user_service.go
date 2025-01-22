package services

import (
	"context"
	"errors"
	"fmt"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
	"point-system-api/pkg/utils"
)

// UserService defines the interface for user-related operations.
type UserService interface {
	CreateUser(ctx context.Context, user models.User) (uint, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user models.User) (bool, error)
	DeleteUser(ctx context.Context, id uint) (bool, error)
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
	ListUsersForSelect(ctx context.Context) ([]map[string]interface{}, error)
	ListUsers(ctx context.Context, page, limit int, filters map[string]interface{}) ([]models.User, int64, error)
}

// userService implements the UserService interface.
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new instance of UserService.
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user in the database.
func (s *userService) CreateUser(ctx context.Context, user models.User) (uint, error) {
	// Check if the username already exists
	existingUser, err := s.userRepo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		return 0, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existingUser != nil {
		return 0, errors.New("username already exists")
	}

	// Hash the user's password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Create the user in the database
	userID, err := s.userRepo.CreateUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

// GetUserByID retrieves a user by their ID.
func (s *userService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by ID: %w", err)
	}
	return user, nil
}

// GetUserByUsername retrieves a user by their username.
func (s *userService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user by username: %w", err)
	}
	return user, nil
}

// UpdateUser updates an existing user in the database.
func (s *userService) UpdateUser(ctx context.Context, user models.User) (bool, error) {
	existingUser, err := s.userRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return false, fmt.Errorf("failed to check existing user: %w", err)
	}
	existingUser.Username = user.Username
	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Role = user.Role
	// Hash the user's password if it's being updated
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return false, fmt.Errorf("failed to hash password: %w", err)
		}
		existingUser.Password = hashedPassword
	}

	// Update the user in the database
	success, err := s.userRepo.UpdateUser(ctx, *existingUser)
	if err != nil {
		return false, fmt.Errorf("failed to update user: %w", err)
	}

	return success, nil
}

// DeleteUser deletes a user by their ID.
func (s *userService) DeleteUser(ctx context.Context, id uint) (bool, error) {
	success, err := s.userRepo.DeleteUser(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}
	return success, nil
}

// AuthenticateUser authenticates a user by their username and password.
func (s *userService) AuthenticateUser(ctx context.Context, username, password string) (*models.User, error) {
	// Retrieve the user by username
	user, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// Check if the provided password matches the stored hashed password
	if !utils.CheckPassword(user.Password, password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// ListUsersForSelect retrieves all users for use in select options.
func (s *userService) ListUsersForSelect(ctx context.Context) ([]map[string]interface{}, error) {
	users, err := s.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Simplify the response for select options
	var result []map[string]interface{}
	for _, user := range users {
		result = append(result, map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
		})
	}

	return result, nil
}

// ListUsers retrieves all users with pagination, filtering, and search.
func (s *userService) ListUsers(ctx context.Context, page, limit int, filters map[string]interface{}) ([]models.User, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	// Call the repository to get paginated and filtered results
	users, total, err := s.userRepo.ListUsersWithFilters(ctx, page, limit, filters)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	return users, total, nil
}
