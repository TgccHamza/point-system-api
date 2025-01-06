package models

import "gorm.io/gorm"

type Company struct {
	ID          uint   `gorm:"primaryKey"`
	CompanyName string `gorm:"size:255;not null;unique"`
	gorm.Model
}
