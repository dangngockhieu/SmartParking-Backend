package iot_device

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	appErrors "backend/internal/common/errors"
	"backend/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateDevice godoc
// @Summary Create IoT device
// @Description Create a new IoT device
// @Tags iot_device
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateIoTDeviceRequest true "IoT device payload"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /iot-devices [post]
func (h *Handler) CreateDevice(c *gin.Context) {
	var req CreateIoTDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Du lieu khong hop le"))
		return
	}

	device, err := h.service.CreateDevice(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "Tao thiet bi IoT thanh cong", device)
}

// FindAllDevices godoc
// @Summary Get IoT devices
// @Description Return all IoT devices
// @Tags iot_device
// @Produce json
// @Security BearerAuth
// @Param lot_id query int false "Filter by parking lot id"
// @Param status query string false "Filter by status (ACTIVE|INACTIVE|ERROR)"
// @Param keyword query string false "Prefix search by mac_address or device_name"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /iot-devices [get]
func (h *Handler) FindAllDevices(c *gin.Context) {
	var query GetIoTDevicesQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.Error(appErrors.NewBadRequest("Query khong hop le"))
		return
	}

	devices, err := h.service.FindAllDevices(query)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lay danh sach thiet bi IoT thanh cong", devices)
}

// UpdateDevice godoc
// @Summary Update IoT device
// @Description Update an IoT device by MAC address
// @Tags iot_device
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param mac_address path string true "Device MAC address"
// @Param request body UpdateIoTDeviceRequest true "IoT device update payload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /iot-devices/{mac_address} [patch]
func (h *Handler) UpdateDevice(c *gin.Context) {
	rawMacAddress := c.Param("mac_address")
	macAddress, err := url.PathUnescape(rawMacAddress)
	if err != nil {
		c.Error(appErrors.NewBadRequest("mac_address khong hop le"))
		return
	}

	var req UpdateIoTDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Du lieu khong hop le"))
		return
	}

	device, err := h.service.UpdateDevice(macAddress, req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Cap nhat thiet bi IoT thanh cong", device)
}
