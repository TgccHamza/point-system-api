package models

import "gorm.io/gorm"

type DeviceModel struct {
	gorm.Model
	Name         string `gorm:"size:255;null"`
	SerialNumber string `gorm:"size:255;not null;unique"`
	CompanyID    uint   `gorm:"null"`
	Location     string `gorm:"size:255;null"` // Embedded Location struct
}
