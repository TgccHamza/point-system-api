package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
)

// WorkDayService defines the interface for workday-related operations.
type WorkDayService interface {
	CreateWorkDay(ctx context.Context, workday *models.WorkDay) error
	GetWorkDayByID(ctx context.Context, id uint) (*models.WorkDay, error)
	ListWorkDays(ctx context.Context) ([]*models.WorkDay, error)
	UpdateWorkDay(ctx context.Context, workday *models.WorkDay) error
	DeleteWorkDay(ctx context.Context, id uint) error
}

// workDayService implements the WorkDayService interface.
type workDayService struct {
	workDayRepo       repositories.WorkDayRepository
	rawAttendanceRepo repositories.RawAttendanceRepository
	attendanceRepo    repositories.AttendanceRepository
}

// NewWorkDayService creates a new instance of WorkDayService.
func NewWorkDayService(workDayRepo repositories.WorkDayRepository, rawAttendanceRepo repositories.RawAttendanceRepository, attendanceRepo repositories.AttendanceRepository) *workDayService {
	return &workDayService{
		workDayRepo:       workDayRepo,
		rawAttendanceRepo: rawAttendanceRepo,
		attendanceRepo:    attendanceRepo,
	}
}

// CreateWorkDay creates a new workday in the database.
func (s *workDayService) CreateWorkDay(ctx context.Context, workday *models.WorkDay) error {
	if workday == nil {
		return errors.New("workday is nil")
	}

	if workday.Date.ToTime().IsZero() {
		return errors.New("date is required")
	}

	if workday.DayType == "" {
		return errors.New("day type is required")
	}

	// Extra validation: cannot create workday for the current day
	currentDate := time.Now().Truncate(24 * time.Hour)
	if !workday.Date.ToTime().Before(currentDate) {
		return errors.New("cannot create workday for the current day or future dates")
	}

	if err := s.workDayRepo.CreateWorkDay(ctx, workday); err != nil {
		return err
	}

	// create raw attendance for all employees that existing in this workday existing view
	employeeAttendances, err := s.workDayRepo.GetEmployeesWithAttendance(ctx, workday.Date.ToTime())
	if err != nil {
		return err
	}

	for _, ea := range employeeAttendances {
		status := determineAttendanceStatus(
			ea.Checkin,
			ea.Checkout)

		rawAttendance := models.RawAttendance{
			WorkDayID: workday.ID,
			CompanyID: ea.CompanyID,
			UserID:    ea.UserID,
			EmployeeName: sql.NullString{
				String: ea.FirstName + " " + ea.LastName,
				Valid:  true,
			},
			Position: sql.NullString{
				String: ea.Qualification,
				Valid:  ea.Qualification != "",
			},
			StartAt: sql.NullString{
				String: ea.Checkin.Time.Format("15:04:05"),
				Valid:  !ea.Checkin.Time.IsZero(),
			},
			EndAt: sql.NullString{
				String: ea.Checkout.Time.Format("15:04:05"),
				Valid:  !ea.Checkout.Time.IsZero(),
			},
			TotalHours: sql.NullFloat64{
				Float64: calculateTotalHours(ea.Checkin.Time, ea.Checkout.Time),
				Valid:   !ea.Checkin.Time.IsZero() && !ea.Checkout.Time.IsZero(),
			},
			Status: status,
			Notes: sql.NullString{
				String: "",
				Valid:  false,
			},
			CalculateOverTime:  false,
			CalculateLunchHour: true,
		}

		// Calculate TotalHourOut based on attendance logs between checkin and checkout if both are not zero
		if !ea.Checkin.Time.IsZero() && !ea.Checkout.Time.IsZero() {
			totalHourOut, err := s.attendanceRepo.GetTotalHourOutByUserAndTimeRange(ctx, ea.RegisterNumber, ea.Checkin.Time, ea.Checkout.Time)
			if err != nil {
				return err
			}
			rawAttendance.TotalHourOut = sql.NullFloat64{
				Float64: totalHourOut,
				Valid:   true,
			}
		} else {
			rawAttendance.TotalHourOut = sql.NullFloat64{Valid: false}
		}

		if err := s.rawAttendanceRepo.CreateRawAttendance(ctx, &rawAttendance); err != nil {
			return err
		}
	}

	return nil
}

// Add the following helper function in the same file

func calculateTotalHours(checkin, checkout time.Time) float64 {
	if checkin.IsZero() || checkout.IsZero() {
		return 0
	}

	// Extract time components
	checkInTime := time.Date(2000, 1, 1, checkin.Hour(), checkin.Minute(), 0, 0, time.UTC)
	checkOutTime := time.Date(2000, 1, 1, checkout.Hour(), checkout.Minute(), 0, 0, time.UTC)

	// If checkout is before checkin, add 24 hours to checkout (next day)
	if checkOutTime.Before(checkInTime) {
		checkOutTime = checkOutTime.Add(24 * time.Hour)
	}

	return checkOutTime.Sub(checkInTime).Hours()
}

func determineAttendanceStatus(checkin, checkout sql.NullTime) sql.NullString {
	if !checkin.Valid || !checkout.Valid {
		return sql.NullString{String: "absent", Valid: true}
	}

	totalHours := calculateTotalHours(checkin.Time, checkout.Time)

	if totalHours > 0 {
		return sql.NullString{String: "present", Valid: true}
	}

	return sql.NullString{String: "absent", Valid: true}
}

// GetWorkDayByID retrieves a workday by its ID.
func (s *workDayService) GetWorkDayByID(ctx context.Context, id uint) (*models.WorkDay, error) {
	if id == 0 {
		return nil, errors.New("invalid workday ID")
	}

	return s.workDayRepo.GetWorkDayByID(ctx, id)
}

// ListWorkDays retrieves all workdays from the database.
func (s *workDayService) ListWorkDays(ctx context.Context) ([]*models.WorkDay, error) {
	return s.workDayRepo.ListWorkDays(ctx)
}

// UpdateWorkDay updates an existing workday in the database.
func (s *workDayService) UpdateWorkDay(ctx context.Context, workday *models.WorkDay) error {
	if workday == nil || workday.ID == 0 {
		return errors.New("invalid workday data")
	}

	return s.workDayRepo.UpdateWorkDay(ctx, workday)
}

// DeleteWorkDay deletes a workday by its ID.
func (s *workDayService) DeleteWorkDay(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid workday ID")
	}

	return s.workDayRepo.DeleteWorkDay(ctx, id)
}
