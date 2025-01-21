// internal/repositories/device_repo.go
package repositories

import (
	"context"
	"errors"
	"fmt"
	"point-system-api/internal/models"

	"gorm.io/gorm"
)

// DeviceRepository defines the interface for device-related database operations.
type DeviceRepository interface {
	// FindDeviceBySerial checks if a device exists by its SerialNumber.
	FindDeviceBySerial(serialNumber string) (*models.Device, error)

	// CreateDevice adds a new device to the database.
	CreateDevice(device *models.Device) error

	// GetDeviceByID retrieves a device by its ID.
	GetDeviceByID(ctx context.Context, id uint) (*models.Device, error)

	// FilterDevices retrieves devices based on filters.
	FilterDevices(filters map[string]interface{}) ([]*models.Device, error)

	// UpdateDevice updates the details of an existing device.
	UpdateDevice(device *models.Device) error

	// DeleteDevice removes a device from the database by its ID.
	DeleteDevice(id uint) error
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) DeviceRepository {
	return &deviceRepository{db: db}
}

// FindDeviceBySerial checks if a device exists by its SerialNumber
func (r *deviceRepository) FindDeviceBySerial(serialNumber string) (*models.Device, error) {
	var device models.Device
	if err := r.db.Where("serial_number = ?", serialNumber).First(&device).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Device not found
		}
		return nil, fmt.Errorf("failed to find device by serial number: %w", err)
	}
	return &device, nil
}

// CreateDevice adds a new device to the database
func (r *deviceRepository) CreateDevice(device *models.Device) error {
	return r.db.Create(device).Error
}

// GetDeviceByID retrieves a device by its ID
func (r *deviceRepository) GetDeviceByID(ctx context.Context, id uint) (*models.Device, error) {
	var device models.Device
	if err := r.db.WithContext(ctx).First(&device, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Device not found
		}
		return nil, fmt.Errorf("failed to retrieve device by ID: %w", err)
	}
	return &device, nil
}

// FilterDevices retrieves devices based on filters
func (r *deviceRepository) FilterDevices(filters map[string]interface{}) ([]*models.Device, error) {
	var devices []*models.Device
	query := r.db.Model(&models.Device{})
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	if err := query.Find(&devices).Error; err != nil {
		return nil, fmt.Errorf("failed to filter devices: %w", err)
	}
	return devices, nil
}

// UpdateDevice updates the details of an existing device
func (r *deviceRepository) UpdateDevice(device *models.Device) error {
	return r.db.Save(device).Error
}

// DeleteDevice removes a device from the database by its ID
func (r *deviceRepository) DeleteDevice(id uint) error {
	return r.db.Delete(&models.Device{}, id).Error
}
