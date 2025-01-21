package models

import (
	"time"

	"gorm.io/gorm"
)

// EmployeeWorkday represents a single entry in the table
type EmployeeWorkDay struct {
	ID         uint      `gorm:"primaryKey"`
	WorkDayID  uint      `gorm:"not null"`
	EmployeeID uint      `gorm:"not null"`
	StartTime  time.Time // Start time of work
	EndTime    time.Time // End time of work
	WorkHours  float64   `gorm:"not null"` // Number of hours worked
	Overtime   float64   // Overtime hours
	Breaks     float64   // Break time in hours
	IsFreeDay  bool      `gorm:"not null;default:false"` // Whether the day is a free day (true/false)
	Status     string    `gorm:"size:50;not null"`       // Status: present, absent, late, early-leave
	Notes      string    `gorm:"size:500;not null"`      // Additional notes or comments (increased size)
	gorm.Model           // Adds fields like ID, CreatedAt, UpdatedAt, DeletedAt
}
