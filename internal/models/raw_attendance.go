package models

import (
	"time"

	"gorm.io/gorm"
)

// RawAttendance represents the raw data from the biometric system.
type RawAttendance struct {
	gorm.Model
	WorkDayID uint      `gorm:"not null"` // Foreign key to WorkDay
	UserID    uint      `gorm:"not null"` // Status: punch-in, punch-out
	Timestamp time.Time `gorm:"not null"` // Timestamp of the punch
	Status    uint      `gorm:"not null"` // Status: always 1
	Punch     uint      `gorm:"not null"` // Status: punch-in (1), punch-out (0)
}
