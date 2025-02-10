package models

import (
	"database/sql"

	"gorm.io/gorm"
)

// RawAttendance represents the raw data from the biometric system.
type RawAttendance struct {
	gorm.Model
	WorkDayID    uint           `gorm:"not null"`
	CompanyID    uint           `gorm:"not null"`
	UserID       uint           `gorm:"not null"`
	EmployeeName sql.NullString `gorm:"type:varchar(255)"`
	Position     sql.NullString `gorm:"type:varchar(255)"`
	StartAt      sql.NullString `gorm:"type:varchar(8)"`
	EndAt        sql.NullString `gorm:"type:varchar(8)"`
	TotalHours   sql.NullFloat64
	Status       sql.NullString `gorm:"type:varchar(50)"`
	Notes        sql.NullString `gorm:"type:varchar(500)"`
	// New field: TotalHourOut in company calculated from AttendanceLog between checkin and checkout, if both not null; otherwise null.
	TotalHourOut sql.NullFloat64
	// New field: CalculateOverTime always false until user modifies to confirm calculation over time.
	CalculateOverTime  bool `gorm:"default:false"`
	CalculateLunchHour bool `gorm:"default:true"`
}
