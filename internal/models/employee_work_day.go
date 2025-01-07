package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeWorkday represents a single entry in the table
type EmployeeWorkday struct {
	EmployeeID uint      `gorm:"not null"`               // Foreign key referencing the Employee table
	Date       time.Time `gorm:"not null"`               // Date of the workday
	WorkHours  float64   `gorm:"not null"`               // Number of hours worked
	IsFreeDay  bool      `gorm:"not null;default:false"` // Whether the day is a free day (true/false)
	Notes      string    `gorm:"size:500;not null"`      // Additional notes or comments (increased size)
	gorm.Model           // Adds fields like ID, CreatedAt, UpdatedAt, DeletedAt
}
