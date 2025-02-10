package services

import (
	"context"
	"database/sql"
	"errors"
	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
	"time"
)

type RawAttendanceService interface {
	CreateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error
	GetRawAttendanceByID(ctx context.Context, id uint) (*models.RawAttendance, error)
	GetRawAttendancesByCompanyIDAndWorkDay(ctx context.Context, companyID uint, workDayID uint) ([]*models.RawAttendance, error)
	UpdateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance, id uint) error
	DeleteRawAttendance(ctx context.Context, id uint) error
	ListRawAttendances(ctx context.Context) ([]*models.RawAttendance, error)
}

type rawAttendanceService struct {
	rawAttendanceRepo repositories.RawAttendanceRepository
}

func NewRawAttendanceService(rawAttendanceRepo repositories.RawAttendanceRepository) RawAttendanceService {
	return &rawAttendanceService{
		rawAttendanceRepo: rawAttendanceRepo,
	}
}

func (s *rawAttendanceService) CreateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error {
	if rawAttendance == nil {
		return errors.New("raw attendance is nil")
	}
	return s.rawAttendanceRepo.CreateRawAttendance(ctx, rawAttendance)
}

func (s *rawAttendanceService) GetRawAttendanceByID(ctx context.Context, id uint) (*models.RawAttendance, error) {
	return s.rawAttendanceRepo.GetRawAttendanceByID(ctx, id)
}

func (s *rawAttendanceService) GetRawAttendancesByCompanyIDAndWorkDay(ctx context.Context, companyID uint, workDayID uint) ([]*models.RawAttendance, error) {
	return s.rawAttendanceRepo.GetRawAttendancesByCompanyIDAndWorkDay(ctx, companyID, workDayID)
}

func (s *rawAttendanceService) UpdateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance, id uint) error {
	if rawAttendance == nil {
		return errors.New("raw attendance is nil")
	}

	if rawAttendance.StartAt.String != "" && rawAttendance.EndAt.String != "" {
		const layoutWithoutSeconds = "15:04"
		// Use only first 5 characters to ignore seconds
		if len(rawAttendance.StartAt.String) < 5 || len(rawAttendance.EndAt.String) < 5 {
			rawAttendance.TotalHours = sql.NullFloat64{Valid: false}
		} else {
			start, errStart := time.Parse(layoutWithoutSeconds, rawAttendance.StartAt.String[:5])
			end, errEnd := time.Parse(layoutWithoutSeconds, rawAttendance.EndAt.String[:5])
			if errStart == nil && errEnd == nil {
				// Adjust for overnight shift
				if end.Before(start) {
					end = end.Add(24 * time.Hour)
				}
				hours := end.Sub(start).Hours()
				rawAttendance.TotalHours = sql.NullFloat64{Float64: hours, Valid: true}
			} else {
				rawAttendance.TotalHours = sql.NullFloat64{Valid: false}
			}
		}
	}

	return s.rawAttendanceRepo.UpdateRawAttendance(ctx, rawAttendance, id)
}

func (s *rawAttendanceService) DeleteRawAttendance(ctx context.Context, id uint) error {
	return s.rawAttendanceRepo.DeleteRawAttendance(ctx, id)
}

func (s *rawAttendanceService) ListRawAttendances(ctx context.Context) ([]*models.RawAttendance, error) {
	return s.rawAttendanceRepo.ListRawAttendances(ctx)
}
