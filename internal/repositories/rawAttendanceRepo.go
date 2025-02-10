package repositories

import (
	"context"
	"errors"
	"point-system-api/internal/models"

	"gorm.io/gorm"
)

// RawAttendanceRepository defines the interface for raw attendance-related database operations.
type RawAttendanceRepository interface {
	CreateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error
	GetRawAttendancesByCompanyIDAndWorkDay(ctx context.Context, companyID uint, workDayID uint) ([]*models.RawAttendance, error)
	GetRawAttendanceByID(ctx context.Context, id uint) (*models.RawAttendance, error)
	UpdateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance, id uint) error
	DeleteRawAttendance(ctx context.Context, id uint) error
	ListRawAttendances(ctx context.Context) ([]*models.RawAttendance, error)
}

type rawAttendanceRepo struct {
	db *gorm.DB
}

func NewRawAttendanceRepo(db *gorm.DB) RawAttendanceRepository {
	return &rawAttendanceRepo{db: db}
}

func (r *rawAttendanceRepo) CreateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error {
	return r.db.WithContext(ctx).Create(rawAttendance).Error
}

func (r *rawAttendanceRepo) GetRawAttendancesByCompanyIDAndWorkDay(ctx context.Context, companyID uint, workDayID uint) ([]*models.RawAttendance, error) {
	var rawAttendances []*models.RawAttendance

	err := r.db.WithContext(ctx).
		Where("company_id = ? AND work_day_id = ?", companyID, workDayID).
		Find(&rawAttendances).Error

	if err != nil {
		return nil, err
	}

	return rawAttendances, nil
}

func (r *rawAttendanceRepo) GetRawAttendanceByID(ctx context.Context, id uint) (*models.RawAttendance, error) {
	var rawAttendance models.RawAttendance
	if err := r.db.WithContext(ctx).First(&rawAttendance, id).Error; err != nil {
		return nil, err
	}

	return &rawAttendance, nil
}

func (s *rawAttendanceRepo) UpdateRawAttendance(ctx context.Context, rawAtt *models.RawAttendance, id uint) error {
	return s.db.WithContext(ctx).
		Model(&models.RawAttendance{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"start_at":    rawAtt.StartAt,
			"end_at":      rawAtt.EndAt,
			"notes":       rawAtt.Notes,
			"status":      rawAtt.Status,
			"total_hours": rawAtt.TotalHours,
		}).Error
}

func (r *rawAttendanceRepo) DeleteRawAttendance(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.RawAttendance{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("raw attendance not found")
	}
	return nil
}

func (r *rawAttendanceRepo) ListRawAttendances(ctx context.Context) ([]*models.RawAttendance, error) {
	var rawAttendances []*models.RawAttendance
	if err := r.db.WithContext(ctx).Find(&rawAttendances).Error; err != nil {
		return nil, err
	}
	return rawAttendances, nil
}
