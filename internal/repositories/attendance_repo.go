package repositories

import (
	"context"

	"point-system-api/internal/models"

	"gorm.io/gorm"
)

// AttendanceRepository defines the interface for attendance-related database operations.
type AttendanceRepository interface {
	// CreateAttendanceLog adds a new attendance log to the database.
	CreateAttendanceLog(ctx context.Context, attendanceLog *models.AttendanceLog) error

	// GetAttendanceByID retrieves an attendance log by its ID.
	GetAttendanceByID(ctx context.Context, id uint) (*models.AttendanceLog, error)

	// GetAllAttendanceLogs retrieves all attendance logs with optional filters.
	GetAllAttendanceLogs(ctx context.Context, filters map[string]interface{}) ([]models.AttendanceLog, error)

	// UpdateAttendanceLog updates an existing attendance log.
	UpdateAttendanceLog(ctx context.Context, attendanceLog *models.AttendanceLog) error

	// DeleteAttendanceLog deletes an attendance log by its ID.
	DeleteAttendanceLog(ctx context.Context, id uint) error
}

type attendanceRepository struct {
	db *gorm.DB
}

// NewAttendanceRepository creates a new instance of AttendanceRepository
func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db: db}
}

// CreateAttendanceLog adds a new attendance log to the database
func (r *attendanceRepository) CreateAttendanceLog(ctx context.Context, attendanceLog *models.AttendanceLog) error {
	return r.db.WithContext(ctx).Create(attendanceLog).Error
}

// GetAttendanceByID retrieves an attendance log by its ID
func (r *attendanceRepository) GetAttendanceByID(ctx context.Context, id uint) (*models.AttendanceLog, error) {
	var attendanceLog models.AttendanceLog
	if err := r.db.WithContext(ctx).First(&attendanceLog, id).Error; err != nil {
		return nil, err
	}
	return &attendanceLog, nil
}

// GetAllAttendanceLogs retrieves all attendance logs with optional filters
func (r *attendanceRepository) GetAllAttendanceLogs(ctx context.Context, filters map[string]interface{}) ([]models.AttendanceLog, error) {
	var attendanceLogs []models.AttendanceLog
	query := r.db.WithContext(ctx).Model(&models.AttendanceLog{})
	for key, value := range filters {
		query = query.Where(key+" = ?", value)
	}
	if err := query.Find(&attendanceLogs).Error; err != nil {
		return nil, err
	}
	return attendanceLogs, nil
}

// UpdateAttendanceLog updates an existing attendance log
func (r *attendanceRepository) UpdateAttendanceLog(ctx context.Context, attendanceLog *models.AttendanceLog) error {
	return r.db.WithContext(ctx).Save(attendanceLog).Error
}

// DeleteAttendanceLog deletes an attendance log by its ID
func (r *attendanceRepository) DeleteAttendanceLog(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.AttendanceLog{}, id).Error
}
