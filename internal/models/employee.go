package models

import "gorm.io/gorm"

type Employee struct {
	ID                 uint   `gorm:"primaryKey"`
	UserID             uint   `gorm:"not null"` // Foreign key to User
	RegistrationNumber string `gorm:"size:255;not null;unique"`
	Qualification      string `gorm:"size:255;not null"`
	CompanyID          uint   `gorm:"not null"` // Foreign key to Company
	gorm.Model
}
