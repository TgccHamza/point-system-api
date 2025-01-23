package repositories

import (
	"context"
	"fmt"

	"point-system-api/internal/models"
	"point-system-api/internal/types"

	"gorm.io/gorm"
)

// AttendanceRepository defines the interface for attendance-related database operations.
type AttendanceRepository interface {
	// CreateAttendanceLog adds a new attendance log to the database.
	CreateAttendanceLog(ctx context.Context, attendanceLog *models.AttendanceLog) error

	// GetAttendanceByID retrieves an attendance log by its ID.
	GetAttendanceByID(ctx context.Context, id uint) (*models.AttendanceLog, error)

	// GetAllAttendanceLogs retrieves all attendance logs with optional filters.
	GetAllAttendanceLogsWithFilters(ctx context.Context, page, limit int, filters map[string]interface{}, search string) ([]types.AttendanceLogResponse, int64, error)

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

func (r *attendanceRepository) GetAllAttendanceLogsWithFilters(ctx context.Context, page, limit int, filters map[string]interface{}, search string) ([]types.AttendanceLogResponse, int64, error) {
	var attendanceLogs []types.AttendanceLogResponse
	var total int64
	offset := (page - 1) * limit

	query := r.db.WithContext(ctx).Table("attendance_logs al").
		Select(`
            al.*,
            e.registration_number as employee_registration,
            e.qualification as employee_qualification,
            e.company_id as employee_company_id,
            e.start_hour as employee_start_hour,
            e.end_hour as employee_end_hour,
            u.first_name as employee_first_name,
            u.last_name as employee_last_name,
            u.username as employee_username,
            u.role as employee_role
        `).
		Joins("JOIN employees e ON CAST(al.user_id AS CHAR) = e.registration_number").
		Joins("JOIN users u ON e.user_id = u.id").
		Order("al.timestamp DESC")

	// Apply filters
	for key, value := range filters {
		if key == "search" {
			continue
		}
		if key == "company_id" {
			query = query.Where("e.company_id = ?", value)
		} else if key == "start_date" {
			query = query.Where("al.timestamp >= ?", value)
		} else if key == "end_date" {
			query = query.Where("al.timestamp <= ?", value)
		} else {
			query = query.Where(fmt.Sprintf("al.%s = ?", key), value)
		}
	}

	// Apply search if provided
	if search != "" {
		query = query.Where(`
        al.serial_number LIKE ? OR 
        CAST(al.user_id AS CHAR) LIKE ? OR 
        e.registration_number LIKE ? OR 
        e.qualification LIKE ? OR 
        e.company_id LIKE ? OR 
        u.first_name LIKE ? OR 
        u.last_name LIKE ? OR 
        u.username LIKE ? OR 
        u.role LIKE ?`,
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count attendance logs: %w", err)
	}

	// Get paginated records
	if err := query.Offset(offset).Limit(limit).Find(&attendanceLogs).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch attendance logs: %w", err)
	}

	return attendanceLogs, total, nil
}
