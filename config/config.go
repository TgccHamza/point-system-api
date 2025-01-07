package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the configuration values for the application.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort int
}

// LoadConfig loads the configuration from environment variables.
func LoadConfig() *Config {
	// Load environment variables from .env file (if it exists)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Parse server port from environment variable
	serverPort, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		log.Fatalf("Invalid server port: %v", err)
	}

	return &Config{
		DBHost:     getEnv("BLUEPRINT_DB_HOST", "localhost"),
		DBPort:     getEnv("BLUEPRINT_DB_PORT", "3306"),
		DBUser:     getEnv("BLUEPRINT_DB_USERNAME", "root"),
		DBPassword: getEnv("BLUEPRINT_DB_PASSWORD", "password"),
		DBName:     getEnv("BLUEPRINT_DB_DATABASE", "point_system_db"),
		ServerPort: serverPort,
	}
}

// getEnv retrieves the value of an environment variable or returns a default value.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
