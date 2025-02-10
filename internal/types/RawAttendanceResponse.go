package types

// RawAttendanceResponse defines the response format for raw attendance.
type RawAttendanceResponse struct {
	ID           uint     `json:"id"`
	CreatedAt    *string  `json:"created_at"`
	UpdatedAt    *string  `json:"updated_at"`
	WorkDayID    *uint    `json:"work_day_id"`
	CompanyID    *uint    `json:"company_id"`
	UserID       *uint    `json:"user_id"`
	EmployeeName *string  `json:"employee_name"`
	Position     *string  `json:"position"`
	StartAt      *string  `json:"start_at"`
	EndAt        *string  `json:"end_at"`
	TotalHours   *float64 `json:"total_hours"`
	Status       *string  `json:"status"`
	Notes        *string  `json:"notes"`
}
