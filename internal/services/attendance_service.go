package services

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
	"point-system-api/pkg/utils"
)

// AttendanceService defines the interface for attendance-related business logic.
type AttendanceService interface {
	// CreateAttendanceLog creates a new attendance log in the database.
	CreateAttendanceLog(ctx context.Context, serialNumber string, hexData string) (*models.AttendanceLog, error)

	// GetAttendanceByID retrieves an attendance log by its ID.
	GetAttendanceByID(ctx context.Context, id uint) (*models.AttendanceLog, error)

	// GetAllAttendanceLogs retrieves all attendance logs with optional filters.
	GetAllAttendanceLogs(ctx context.Context, filters map[string]interface{}) ([]models.AttendanceLog, error)

	// UpdateAttendanceLog updates an existing attendance log.
	UpdateAttendanceLog(ctx context.Context, attendanceLog *models.AttendanceLog) error

	// DeleteAttendanceLog deletes an attendance log by its ID.
	DeleteAttendanceLog(ctx context.Context, id uint) error
}

// AttendanceService handles business logic for attendance logs.
type attendanceService struct {
	deviceRepo     repositories.DeviceRepository
	attendanceRepo repositories.AttendanceRepository
}

// NewAttendanceService creates a new instance of AttendanceService.
func NewAttendanceService(deviceRepo repositories.DeviceRepository,
	attendanceRepo repositories.AttendanceRepository) AttendanceService {
	return &attendanceService{
		deviceRepo:     deviceRepo,
		attendanceRepo: attendanceRepo,
	}
}

// CreateAttendanceLog processes hex data, checks/creates the device, and saves the attendance log to the database.
func (s *attendanceService) CreateAttendanceLog(ctx context.Context, serialNumber string, hexData string) (*models.AttendanceLog, error) {
	// Convert the hex string to bytes
	byteData, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, errors.New("invalid hex data")
	}

	// Unpack the byte data into the struct
	var record struct {
		UID       uint16
		UserID    [24]byte
		Status    uint8
		Timestamp [4]byte
		Punch     uint8
		Space     [8]byte
	}

	buf := bytes.NewReader(byteData)
	err = binary.Read(buf, binary.LittleEndian, &record)
	if err != nil {
		return nil, errors.New("failed to unpack data")
	}

	// Clean the UserID by removing null bytes and converting to a string
	userIDClean := string(bytes.TrimRight(record.UserID[:], "\x00"))

	// Convert UserID to an integer
	var userIDNumber int
	_, err = fmt.Sscanf(userIDClean, "%d", &userIDNumber)
	if err != nil {
		return nil, errors.New("failed to parse user ID")
	}

	// Decode the timestamp
	timestamp := binary.LittleEndian.Uint32(record.Timestamp[:])
	decodedTime := utils.DecodeTime(timestamp)

	// Check if the device exists, or create it if not
	device, err := s.deviceRepo.FindDeviceBySerial(serialNumber)
	if device == nil {
		// Device does not exist; create it
		device = &models.Device{
			SerialNumber: serialNumber,
			Name:         "Unknown Device",   // Placeholder, modify as needed
			Location:     "Unknown Location", // Placeholder, modify as needed
			CompanyID:    0,
		}
		if err := s.deviceRepo.CreateDevice(device); err != nil {
			return nil, fmt.Errorf("failed to create device: %w", err)
		}
	}

	// Create the attendance log model
	attendanceLog := models.AttendanceLog{
		SerialNumber: serialNumber,
		UID:          record.UID,
		UserID:       userIDNumber,
		Status:       record.Status,
		Punch:        record.Punch,
		Timestamp:    decodedTime,
	}

	// Save the attendance log to the database
	result := s.attendanceRepo.CreateAttendanceLog(ctx, &attendanceLog)
	if result != nil {
		return nil, result
	}

	return &attendanceLog, nil
}

// GetAttendanceByID retrieves an attendance log by its ID.
func (s *attendanceService) GetAttendanceByID(ctx context.Context, id uint) (*models.AttendanceLog, error) {
	// Validate the attendance log ID
	if id == 0 {
		return nil, errors.New("attendance log ID is required")
	}

	// Retrieve the attendance log from the database
	attendanceLog, err := s.attendanceRepo.GetAttendanceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attendance log: %w", err)
	}

	// Check if the attendance log exists
	if attendanceLog == nil {
		return nil, errors.New("attendance log not found")
	}

	return attendanceLog, nil
}

// GetAllAttendanceLogs retrieves all attendance logs with optional filters.
func (s *attendanceService) GetAllAttendanceLogs(ctx context.Context, filters map[string]interface{}) ([]models.AttendanceLog, error) {
	// Retrieve all attendance logs from the database
	attendanceLogs, err := s.attendanceRepo.GetAllAttendanceLogs(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve attendance logs: %w", err)
	}

	// Return an empty slice if no attendance logs are found
	if len(attendanceLogs) == 0 {
		return []models.AttendanceLog{}, nil
	}

	return attendanceLogs, nil
}

// UpdateAttendanceLog updates an existing attendance log.
func (s *attendanceService) UpdateAttendanceLog(ctx context.Context, attendanceLog *models.AttendanceLog) error {
	// Validate the attendance log ID
	if attendanceLog.ID == 0 {
		return errors.New("attendance log ID is required")
	}

	// Check if the attendance log exists
	existingLog, err := s.attendanceRepo.GetAttendanceByID(ctx, attendanceLog.ID)
	if err != nil {
		return fmt.Errorf("failed to check existing attendance log: %w", err)
	}
	if existingLog == nil {
		return errors.New("attendance log not found")
	}

	// Update the attendance log in the database
	if err := s.attendanceRepo.UpdateAttendanceLog(ctx, attendanceLog); err != nil {
		return fmt.Errorf("failed to update attendance log: %w", err)
	}

	return nil
}

// DeleteAttendanceLog deletes an attendance log by its ID.
func (s *attendanceService) DeleteAttendanceLog(ctx context.Context, id uint) error {
	// Validate the attendance log ID
	if id == 0 {
		return errors.New("attendance log ID is required")
	}

	// Check if the attendance log exists
	existingLog, err := s.attendanceRepo.GetAttendanceByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check existing attendance log: %w", err)
	}
	if existingLog == nil {
		return errors.New("attendance log not found")
	}

	// Delete the attendance log from the database
	if err := s.attendanceRepo.DeleteAttendanceLog(ctx, id); err != nil {
		return fmt.Errorf("failed to delete attendance log: %w", err)
	}

	return nil
}
