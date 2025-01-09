package models

import (
	"time"

	"gorm.io/gorm"
)

// WorkDay represents a single workday.
type WorkDay struct {
	gorm.Model
	Date    time.Time `gorm:"not null"`         // Date of the workday
	DayType string    `gorm:"size:50;not null"` // Type of day: workday, free, holiday
}
