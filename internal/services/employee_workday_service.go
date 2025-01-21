package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
)

// EmployeeWorkDayService defines the interface for employee workday-related operations.
type EmployeeWorkDayService interface {
	GenerateEmployeeWorkDay(ctx context.Context, workDayID uint) error
	UpdateEmployeeWorkDay(ctx context.Context, employeeWorkDayID uint, notes string, status string) error
}

// employeeWorkDayService implements the EmployeeWorkDayService interface.
type employeeWorkDayService struct {
	workDayRepo         repositories.WorkDayRepository
	rawAttendanceRepo   repositories.RawAttendanceRepository
	employeeWorkDayRepo repositories.EmployeeWorkDayRepository
}

// NewEmployeeWorkDayService creates a new instance of EmployeeWorkDayService.
func NewEmployeeWorkDayService(
	workDayRepo repositories.WorkDayRepository,
	rawAttendanceRepo repositories.RawAttendanceRepository,
	employeeWorkDayRepo repositories.EmployeeWorkDayRepository,
) EmployeeWorkDayService {
	return &employeeWorkDayService{
		workDayRepo:         workDayRepo,
		rawAttendanceRepo:   rawAttendanceRepo,
		employeeWorkDayRepo: employeeWorkDayRepo,
	}
}

// GenerateEmployeeWorkDay generates EmployeeWorkDay records based on WorkDay and RawAttendance.
func (s *employeeWorkDayService) GenerateEmployeeWorkDay(ctx context.Context, workDayID uint) error {
	// Retrieve the work day
	workDay, err := s.workDayRepo.GetWorkDayByID(ctx, workDayID)
	if err != nil {
		return fmt.Errorf("failed to retrieve work day: %w", err)
	}
	if workDay == nil {
		return errors.New("work day not found")
	}

	// Retrieve all raw attendances for the work day
	rawAttendances, err := s.rawAttendanceRepo.GetRawAttendancesByWorkDayID(ctx, workDayID)
	if err != nil {
		return fmt.Errorf("failed to retrieve raw attendances: %w", err)
	}

	// Group raw attendances by employee ID
	attendanceMap := make(map[uint][]*models.RawAttendance)
	for _, rawAttendance := range rawAttendances {
		attendanceMap[rawAttendance.UserID] = append(attendanceMap[rawAttendance.UserID], rawAttendance)
	}

	// Generate EmployeeWorkDay records for each employee
	for employeeID, attendances := range attendanceMap {
		// Calculate start time, end time, and work hours
		var startTime, endTime time.Time
		var workHours float64

		// Assuming the first punch is the start time and the last punch is the end time
		if len(attendances) > 0 {
			startTime = attendances[0].Timestamp
			endTime = attendances[len(attendances)-1].Timestamp
			workHours = endTime.Sub(startTime).Hours()
		}

		// Create the EmployeeWorkDay record
		employeeWorkDay := &models.EmployeeWorkDay{
			WorkDayID:  workDayID,
			EmployeeID: employeeID,
			StartTime:  startTime,
			EndTime:    endTime,
			WorkHours:  workHours,
			Status:     "present", // Default status
			Notes:      "",        // Default empty notes
		}

		// Save the EmployeeWorkDay record
		if err := s.employeeWorkDayRepo.CreateEmployeeWorkDay(ctx, employeeWorkDay); err != nil {
			return fmt.Errorf("failed to create employee work day: %w", err)
		}
	}

	return nil
}

// UpdateEmployeeWorkDay updates an existing EmployeeWorkDay record with notes and status.
func (s *employeeWorkDayService) UpdateEmployeeWorkDay(ctx context.Context, employeeWorkDayID uint, notes string, status string) error {
	if employeeWorkDayID == 0 {
		return errors.New("invalid employee work day ID")
	}

	// Retrieve the EmployeeWorkDay record
	employeeWorkDay, err := s.employeeWorkDayRepo.GetEmployeeWorkDayByID(ctx, employeeWorkDayID)
	if err != nil {
		return fmt.Errorf("failed to retrieve employee work day: %w", err)
	}
	if employeeWorkDay == nil {
		return errors.New("employee work day not found")
	}

	// Update the notes and status
	employeeWorkDay.Notes = notes
	employeeWorkDay.Status = status

	// Save the updated EmployeeWorkDay record
	if err := s.employeeWorkDayRepo.UpdateEmployeeWorkDay(ctx, employeeWorkDay); err != nil {
		return fmt.Errorf("failed to update employee work day: %w", err)
	}

	return nil
}
