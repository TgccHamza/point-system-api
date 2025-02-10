package types

import (
	"database/sql"
	"time"
)

type EmployeeAttendance struct {
	UserID         uint
	CompanyID      uint
	FirstName      string
	LastName       string
	RegisterNumber uint
	Qualification  string
	Date           time.Time
	Checkin        sql.NullTime
	Checkout       sql.NullTime
}
