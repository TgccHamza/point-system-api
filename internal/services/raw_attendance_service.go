package services

import (
	"context"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
)

// RawAttendanceService defines the interface for raw attendance-related operations.
type RawAttendanceService interface {
	CreateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error
	CreateManyRawAttendances(ctx context.Context, rawAttendances []*models.RawAttendance) error
	GetRawAttendanceByID(ctx context.Context, id uint) (*models.RawAttendance, error)
	GetRawAttendancesByWorkDayID(ctx context.Context, workDayID uint) ([]*models.RawAttendance, error)
	UpdateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error
	DeleteRawAttendance(ctx context.Context, id uint) error
}

// rawAttendanceService implements the RawAttendanceService interface.
type rawAttendanceService struct {
	rawAttendanceRepo repositories.RawAttendanceRepository
}

// NewRawAttendanceService creates a new instance of RawAttendanceService.
func NewRawAttendanceService(rawAttendanceRepo repositories.RawAttendanceRepository) RawAttendanceService {
	return &rawAttendanceService{
		rawAttendanceRepo: rawAttendanceRepo,
	}
}

// CreateRawAttendance creates a new raw attendance record in the database.
func (s *rawAttendanceService) CreateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error {
	return s.rawAttendanceRepo.CreateRawAttendance(ctx, rawAttendance)
}

// CreateManyRawAttendances creates multiple raw attendance records in the database.
func (s *rawAttendanceService) CreateManyRawAttendances(ctx context.Context, rawAttendances []*models.RawAttendance) error {
	return s.rawAttendanceRepo.CreateManyRawAttendances(ctx, rawAttendances)
}

// GetRawAttendanceByID retrieves a raw attendance record by its ID.
func (s *rawAttendanceService) GetRawAttendanceByID(ctx context.Context, id uint) (*models.RawAttendance, error) {
	return s.rawAttendanceRepo.GetRawAttendanceByID(ctx, id)
}

// GetRawAttendancesByWorkDayID retrieves all raw attendance records for a specific work day.
func (s *rawAttendanceService) GetRawAttendancesByWorkDayID(ctx context.Context, workDayID uint) ([]*models.RawAttendance, error) {
	return s.rawAttendanceRepo.GetRawAttendancesByWorkDayID(ctx, workDayID)
}

// UpdateRawAttendance updates an existing raw attendance record in the database.
func (s *rawAttendanceService) UpdateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error {
	return s.rawAttendanceRepo.UpdateRawAttendance(ctx, rawAttendance)
}

// DeleteRawAttendance deletes a raw attendance record by its ID.
func (s *rawAttendanceService) DeleteRawAttendance(ctx context.Context, id uint) error {
	return s.rawAttendanceRepo.DeleteRawAttendance(ctx, id)
}
