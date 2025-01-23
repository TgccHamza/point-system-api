// internal/services/device_service.go
package services

import (
	"context"
	"errors"
	"fmt"

	"point-system-api/internal/models"
	"point-system-api/internal/repositories"
)

type DeviceService interface {
	CreateDevice(ctx context.Context, device *models.Device) error
	GetDeviceByID(ctx context.Context, id uint) (*models.Device, error)
	GetAllDevices(ctx context.Context, page, limit int, filters map[string]interface{}, search string) ([]models.Device, int64, error)
	UpdateDevice(ctx context.Context, device *models.Device) error
	DeleteDevice(ctx context.Context, id uint) error
}

type deviceService struct {
	deviceRepo repositories.DeviceRepository
}

func NewDeviceService(deviceRepo repositories.DeviceRepository) DeviceService {
	return &deviceService{
		deviceRepo: deviceRepo,
	}
}

// CreateDevice creates a new device in the database after checking if it already exists
func (s *deviceService) CreateDevice(ctx context.Context, device *models.Device) error {
	// Check if the device already exists by serial number
	existingDevice, err := s.deviceRepo.FindDeviceBySerial(device.SerialNumber)
	if err != nil {
		return fmt.Errorf("failed to check existing device: %w", err)
	}
	if existingDevice != nil {
		return errors.New("device with this serial number already exists")
	}

	// Create the device in the database
	if err := s.deviceRepo.CreateDevice(device); err != nil {
		return fmt.Errorf("failed to create device: %w", err)
	}

	return nil
}

// GetDeviceByID retrieves a device by its ID
func (s *deviceService) GetDeviceByID(ctx context.Context, id uint) (*models.Device, error) {
	device, err := s.deviceRepo.GetDeviceByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve device by ID: %w", err)
	}
	return device, nil
}

// GetAllDevices retrieves all devices with pagination, filtering, and search.
func (s *deviceService) GetAllDevices(ctx context.Context, page, limit int, filters map[string]interface{}, search string) ([]models.Device, int64, error) {
	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	// Call the repository to get paginated and filtered results
	devices, total, err := s.deviceRepo.ListDevicesWithFilters(ctx, page, limit, filters, search)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch devices: %w", err)
	}

	return devices, total, nil
}

// UpdateDevice updates an existing device in the database
func (s *deviceService) UpdateDevice(ctx context.Context, device *models.Device) error {
	// Validate that the device ID is provided
	if device.ID == 0 {
		return errors.New("device ID is required")
	}

	// Update the device in the database
	if err := s.deviceRepo.UpdateDevice(device); err != nil {
		return fmt.Errorf("failed to update device: %w", err)
	}

	return nil
}

// DeleteDevice deletes a device by its ID
func (s *deviceService) DeleteDevice(ctx context.Context, id uint) error {
	// Validate that the device ID is provided
	if id == 0 {
		return errors.New("device ID is required")
	}

	// Delete the device from the database
	if err := s.deviceRepo.DeleteDevice(id); err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	return nil
}
