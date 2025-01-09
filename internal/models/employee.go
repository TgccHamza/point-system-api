package models

import (
	"time"

	"gorm.io/gorm"
)

type Employee struct {
	ID                 uint      `gorm:"primaryKey"`
	UserID             uint      `gorm:"not null"` // Foreign key to User
	RegistrationNumber string    `gorm:"size:255;not null;unique"`
	Qualification      string    `gorm:"size:255;not null"`
	CompanyID          uint      `gorm:"not null"` // Foreign key to Company
	StartHour          time.Time `gorm:"not null"` // Daily start hour (e.g., 08:00)
	EndHour            time.Time `gorm:"not null"` // Daily end hour (e.g., 18:00)
	gorm.Model
}
