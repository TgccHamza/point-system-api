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
	ListUsers(ctx context.Context) ([]models.User, error)
	UpdateUser(ctx context.Context, user models.User) (bool, error)
	DeleteUser(ctx context.Context, id uint) (bool, error)
	AuthenticateUser(ctx context.Context, username, password string) (*models.User, error)
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

// ListUsers retrieves all users from the database.
func (s *userService) ListUsers(ctx context.Context) ([]models.User, error) {
	users, err := s.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	return users, nil
}

// UpdateUser updates an existing user in the database.
func (s *userService) UpdateUser(ctx context.Context, user models.User) (bool, error) {
	// Hash the user's password if it's being updated
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return false, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	// Update the user in the database
	success, err := s.userRepo.UpdateUser(ctx, user)
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
