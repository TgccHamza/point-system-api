package models

import "gorm.io/gorm"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	FirstName string `gorm:"size:255;not null"`
	LastName  string `gorm:"size:255;not null"`
	Username  string `gorm:"size:255;unique;not null"`
	Password  string `gorm:"not null"`         // Hashed password
	Role      string `gorm:"size:50;not null"` // super-admin, manager, employee
	gorm.Model
}
