package types

type EmployeeWithUser struct {
	ID                 uint   `json:"id"`
	UserID             uint   `json:"user_id"`
	RegistrationNumber string `json:"registration_number"`
	Qualification      string `json:"qualification"`
	CompanyID          uint   `json:"company_id"`
	StartHour          string `json:"start_hour"`
	EndHour            string `json:"end_hour"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	Username           string `json:"username"`
	Role               string `json:"role"`
}
