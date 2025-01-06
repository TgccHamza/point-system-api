package repository

import (
	"fmt"
	"point-system-api/internal/models"

	"gorm.io/gorm"

	"point-system-api/pkg/utils"
)

// CreateUser inserts a new user into the database after hashing the password.
func CreateUser(db *gorm.DB, user models.User) (uint, error) {
	// Hash the user's password before saving
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Create the user in the database
	if err := db.Create(&user).Error; err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}

	// Return the ID of the newly created user
	return user.ID, nil
}

// GetUserByID retrieves a user by their ID.
func GetUserByID(db *gorm.DB, id uint) (*models.User, error) {
	var user models.User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("failed to retrieve user by ID: %w", err)
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No user found
		}
		return nil, fmt.Errorf("failed to retrieve user by username: %w", err)
	}
	return &user, nil
}

// ListUsers retrieves all users from the database.
func ListUsers(db *gorm.DB) ([]models.User, error) {
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}

// UpdateUser updates an existing user in the database.
func UpdateUser(db *gorm.DB, user models.User) (bool, error) {
	// Hash the user's password before updating
	if user.Password != "" {
		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			return false, fmt.Errorf("failed to hash password: %w", err)
		}
		user.Password = hashedPassword
	}

	// Update the user in the database
	if err := db.Save(&user).Error; err != nil {
		return false, fmt.Errorf("failed to update user: %w", err)
	}

	return true, nil
}

// DeleteUser deletes a user by their ID.
func DeleteUser(db *gorm.DB, id uint) (bool, error) {
	if err := db.Delete(&models.User{}, id).Error; err != nil {
		return false, fmt.Errorf("failed to delete user: %w", err)
	}
	return true, nil
}
