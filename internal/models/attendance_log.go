package models

import (
	"time"

	"gorm.io/gorm"
)

// Define the database model for attendance logs
type AttendanceLog struct {
	gorm.Model
	SerialNumber string    `json:"serial_number"` // Serial number of the device
	UID          uint16    `json:"uid"`           // User ID (unsigned short)
	UserID       int       `json:"user_id"`       // User ID as an integer
	Status       uint8     `json:"status"`        // Status of the attendance record
	Punch        uint8     `json:"punch"`         // Punch type (e.g., check-in, check-out)
	Timestamp    time.Time `json:"timestamp"`     // Timestamp of the attendance record
}
