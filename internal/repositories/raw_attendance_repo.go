package repositories

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
)

// RawAttendanceRepository defines the interface for raw attendance-related database operations.
type RawAttendanceRepository interface {
	CreateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error
	CreateManyRawAttendances(ctx context.Context, rawAttendances []*models.RawAttendance) error
	GetRawAttendanceByID(ctx context.Context, id uint) (*models.RawAttendance, error)
	GetRawAttendancesByWorkDayID(ctx context.Context, workDayID uint) ([]*models.RawAttendance, error)
	UpdateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error
	DeleteRawAttendance(ctx context.Context, id uint) error
}

// rawAttendanceRepository implements the RawAttendanceRepository interface.
type rawAttendanceRepository struct {
	db *gorm.DB
}

// NewRawAttendanceRepository creates a new instance of RawAttendanceRepository.
func NewRawAttendanceRepository(db *gorm.DB) RawAttendanceRepository {
	return &rawAttendanceRepository{
		db: db,
	}
}

// CreateRawAttendance inserts a new raw attendance record into the database.
func (r *rawAttendanceRepository) CreateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error {
	if rawAttendance == nil {
		return errors.New("raw attendance is nil")
	}

	if rawAttendance.WorkDayID == 0 {
		return errors.New("work day ID is required")
	}

	if rawAttendance.UserID == 0 {
		return errors.New("user ID is required")
	}

	if rawAttendance.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}

	if err := r.db.WithContext(ctx).Create(rawAttendance).Error; err != nil {
		return fmt.Errorf("failed to create raw attendance: %w", err)
	}

	return nil
}

// CreateManyRawAttendances inserts multiple raw attendance records into the database.
func (r *rawAttendanceRepository) CreateManyRawAttendances(ctx context.Context, rawAttendances []*models.RawAttendance) error {
	if len(rawAttendances) == 0 {
		return errors.New("no raw attendances provided")
	}

	for _, rawAttendance := range rawAttendances {
		if rawAttendance.WorkDayID == 0 {
			return errors.New("work day ID is required")
		}
		if rawAttendance.UserID == 0 {
			return errors.New("user ID is required")
		}
		if rawAttendance.Timestamp.IsZero() {
			return errors.New("timestamp is required")
		}
	}

	if err := r.db.WithContext(ctx).Create(rawAttendances).Error; err != nil {
		return fmt.Errorf("failed to create raw attendances: %w", err)
	}

	return nil
}

// GetRawAttendanceByID retrieves a raw attendance record by its ID.
func (r *rawAttendanceRepository) GetRawAttendanceByID(ctx context.Context, id uint) (*models.RawAttendance, error) {
	if id == 0 {
		return nil, errors.New("invalid raw attendance ID")
	}

	var rawAttendance models.RawAttendance
	if err := r.db.WithContext(ctx).First(&rawAttendance, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No raw attendance found
		}
		return nil, fmt.Errorf("failed to retrieve raw attendance by ID: %w", err)
	}
	return &rawAttendance, nil
}

// GetRawAttendancesByWorkDayID retrieves all raw attendance records for a specific work day.
func (r *rawAttendanceRepository) GetRawAttendancesByWorkDayID(ctx context.Context, workDayID uint) ([]*models.RawAttendance, error) {
	if workDayID == 0 {
		return nil, errors.New("invalid work day ID")
	}

	var rawAttendances []*models.RawAttendance
	if err := r.db.WithContext(ctx).Where("work_day_id = ?", workDayID).Find(&rawAttendances).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve raw attendances by work day ID: %w", err)
	}
	return rawAttendances, nil
}

// UpdateRawAttendance updates an existing raw attendance record in the database.
func (r *rawAttendanceRepository) UpdateRawAttendance(ctx context.Context, rawAttendance *models.RawAttendance) error {
	if rawAttendance == nil || rawAttendance.ID == 0 {
		return errors.New("invalid raw attendance data")
	}

	if rawAttendance.WorkDayID == 0 {
		return errors.New("work day ID is required")
	}

	if rawAttendance.UserID == 0 {
		return errors.New("user ID is required")
	}

	if rawAttendance.Timestamp.IsZero() {
		return errors.New("timestamp is required")
	}

	if err := r.db.WithContext(ctx).Save(rawAttendance).Error; err != nil {
		return fmt.Errorf("failed to update raw attendance: %w", err)
	}

	return nil
}

// DeleteRawAttendance deletes a raw attendance record by its ID.
func (r *rawAttendanceRepository) DeleteRawAttendance(ctx context.Context, id uint) error {
	if id == 0 {
		return errors.New("invalid raw attendance ID")
	}

	if err := r.db.WithContext(ctx).Delete(&models.RawAttendance{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete raw attendance: %w", err)
	}

	return nil
}