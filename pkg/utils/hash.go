package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the user's password using bcrypt before saving it.
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash for the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// Set the hashed password
	return string(hash), nil
}

// CheckPassword compares the given plain text password with the stored hashed password.
// CheckPassword compares the given plain text password with the stored hashed password.
func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
