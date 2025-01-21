package repositories

import (
	"fmt"
	"point-system-api/internal/models"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{db: db}
}

// FindDeviceBySerial checks if a device exists by its SerialNumber
func (r *DeviceRepository) FindDeviceBySerial(serialNumber string, device *models.Device) error {
	return r.db.Where("serial_number = ?", serialNumber).First(device).Error
}

// CreateDevice adds a new device to the database
func (r *DeviceRepository) CreateDevice(device *models.Device) error {
	return r.db.Create(device).Error
}

// FilterDevices retrieves devices based on filters
func (r *DeviceRepository) FilterDevices(filters map[string]interface{}) ([]models.Device, error) {
	var devices []models.Device
	query := r.db.Model(&models.Device{})
	for key, value := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	if err := query.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

// UpdateDevice updates the details of an existing device
func (r *DeviceRepository) UpdateDevice(device *models.Device) error {
	return r.db.Save(device).Error
}

// DeleteDevice removes a device from the database by its ID
func (r *DeviceRepository) DeleteDevice(id uint) error {
	return r.db.Delete(&models.Device{}, id).Error
}
