package models

import (
	"point-system-api/internal/types"

	"gorm.io/gorm"
)

// WorkDay represents a single workday.
type WorkDay struct {
	gorm.Model
	Date    types.DateOnly `gorm:"type:date;not null" json:"date"`  // Date of the workday
	DayType string         `gorm:"size:50;not null" json:"dayType"` // Type of day: workday, free, holiday
}
