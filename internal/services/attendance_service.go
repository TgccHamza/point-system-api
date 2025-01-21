package services

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
	"point-system-api/pkg/utils"
)

// AttendanceService handles business logic for attendance logs.
type AttendanceService struct {
	db             *gorm.DB
	deviceRepo     *repositories.DeviceRepository
	attendanceRepo *repositories.AttendanceRepository
}

// NewAttendanceService creates a new instance of AttendanceService.
func NewAttendanceService(db *gorm.DB, deviceRepo *repositories.DeviceRepository, attendanceRepo *repositories.AttendanceRepository) *AttendanceService {
	return &AttendanceService{
		db:             db,
		deviceRepo:     deviceRepo,
		attendanceRepo: attendanceRepo,
	}
}

// CreateAttendanceLog processes hex data, checks/creates the device, and saves the attendance log to the database.
func (s *AttendanceService) CreateAttendanceLog(ctx context.Context, serialNumber string, hexData string) (*models.AttendanceLog, error) {
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
	var device models.DeviceModel
	err = s.deviceRepo.FindDeviceBySerial(serialNumber, &device)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Device does not exist; create it
			device = models.DeviceModel{
				SerialNumber: serialNumber,
				Name:         "Unknown Device",   // Placeholder, modify as needed
				Location:     "Unknown Location", // Placeholder, modify as needed
				CompanyID:    0,                  // Set default or retrieve dynamically if needed
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			if err := s.deviceRepo.CreateDevice(&device); err != nil {
				return nil, fmt.Errorf("failed to create device: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to check device: %w", err)
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
