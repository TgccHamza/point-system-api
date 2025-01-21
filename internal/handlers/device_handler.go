// internal/handlers/device_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"point-system-api/internal/models"
	"point-system-api/internal/services"
)

// DeviceHandler handles HTTP requests for device-related operations.
type DeviceHandler struct {
	deviceService services.DeviceService
}

// NewDeviceHandler creates a new instance of DeviceHandler.
func NewDeviceHandler(deviceService services.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		deviceService: deviceService,
	}
}

// GetAllDevices retrieves all devices with optional filters.
func (h *DeviceHandler) GetAllDevices(c *gin.Context) {
	// Extract filters from query parameters
	filters := map[string]interface{}{}
	if serial := c.Query("serial_number"); serial != "" {
		filters["serial_number"] = serial
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if location := c.Query("location"); location != "" {
		filters["location"] = location
	}
	if companyID := c.Query("company_id"); companyID != "" {
		companyIDInt, err := strconv.Atoi(companyID)
		if err == nil {
			filters["company_id"] = companyIDInt
		}
	}

	// Retrieve devices from the service
	devices, err := h.deviceService.GetAllDevices(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve devices"})
		return
	}

	c.JSON(http.StatusOK, devices)
}

// GetDeviceByID retrieves a single device by its ID.
func (h *DeviceHandler) GetDeviceByID(c *gin.Context) {
	id := c.Param("id")
	deviceID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	device, err := h.deviceService.GetDeviceByID(c.Request.Context(), uint(deviceID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// CreateDevice handles creating a new device.
func (h *DeviceHandler) CreateDevice(c *gin.Context) {
	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Create the device using the service
	err := h.deviceService.CreateDevice(c.Request.Context(), &device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create device"})
		return
	}

	c.JSON(http.StatusCreated, device)
}

// UpdateDevice handles updating an existing device.
func (h *DeviceHandler) UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	deviceID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	var device models.Device
	if err := c.ShouldBindJSON(&device); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	device.ID = uint(deviceID)

	// Update the device using the service
	err = h.deviceService.UpdateDevice(c.Request.Context(), &device)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device updated successfully"})
}

// DeleteDevice handles deleting a device by its ID.
func (h *DeviceHandler) DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	deviceID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid device ID"})
		return
	}

	// Delete the device using the service
	err = h.deviceService.DeleteDevice(c.Request.Context(), uint(deviceID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete device"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Device deleted successfully"})
}
