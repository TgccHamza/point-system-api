package types

import "time"

type AttendanceLogResponse struct {
	ID                    uint      `json:"id"`
	SerialNumber          string    `json:"serial_number"` // Serial number of the device
	UID                   uint16    `json:"uid"`           // User ID (unsigned short)
	UserID                int       `json:"user_id"`       // User ID as an integer
	Status                uint8     `json:"status"`        // Status of the attendance record
	Punch                 uint8     `json:"punch"`         // Punch type (e.g., check-in, check-out)
	Timestamp             time.Time `json:"timestamp"`     // Timestamp of the attendance record
	EmployeeRegistration  string    `json:"employee_registration"`
	EmployeeQualification string    `json:"employee_qualification"`
	EmployeeCompanyID     int       `json:"employee_company_id"`
	EmployeeStartHour     string    `json:"employee_start_hour"`
	EmployeeEndHour       string    `json:"employee_end_hour"`
	EmployeeFirstName     string    `json:"employee_first_name"`
	EmployeeLastName      string    `json:"employee_last_name"`
	EmployeeUsername      string    `json:"employee_username"`
	EmployeeRole          string    `json:"employee_role"`
}
