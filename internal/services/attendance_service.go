package services

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"point-system-api/internal/models"
	"point-system-api/pkg/utils"
)

// AttendanceService handles business logic for attendance logs.
type AttendanceService struct {
	db *gorm.DB
}

// NewAttendanceService creates a new instance of AttendanceService.
func NewAttendanceService(db *gorm.DB) *AttendanceService {
	return &AttendanceService{
		db: db,
	}
}

// CreateAttendanceLog processes hex data and saves the attendance log to the database.
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

	// Create the attendance log model
	attendanceLog := models.AttendanceLog{
		SerialNumber: serialNumber,
		UID:          record.UID,
		UserID:       userIDNumber,
		Status:       record.Status,
		Punch:        record.Punch,
		Timestamp:    decodedTime,
	}

	// Save the model to the database
	result := s.db.WithContext(ctx).Create(&attendanceLog)
	if result.Error != nil {
		return nil, result.Error
	}

	return &attendanceLog, nil
}
