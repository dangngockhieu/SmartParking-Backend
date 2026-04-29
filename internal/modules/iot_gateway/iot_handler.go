package iot_gateway

import (
	"net/http"

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

// CameraPlate godoc
// @Summary Camera AI gửi biển số
// @Description Lưu biển số tạm thời theo gate để chờ RFID scan consume
// @Tags iot
// @Accept json
// @Produce json
// @Param request body CameraPlateRequest true "Dữ liệu từ camera"
// @Success 200 {object} CameraPlateResponse
// @Failure 400 {object} map[string]interface{}
// @Router /iot/camera [post]
func (h *Handler) CameraPlate(c *gin.Context) {
	var req CameraPlateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	result, err := h.service.HandleCameraPlate(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, result.Message, result)
}

// RfidScan godoc
// @Summary ESP32 quẹt thẻ RFID
// @Description Xử lý sự kiện quẹt thẻ từ ESP32+RC522, quyết định mở/đóng barie
// @Tags iot
// @Accept json
// @Produce json
// @Param request body RfidScanRequest true "Dữ liệu từ ESP32"
// @Success 200 {object} RfidScanResponse
// @Failure 400 {object} map[string]interface{}
// @Router /iot/rfid [post]
func (h *Handler) RfidScan(c *gin.Context) {
	var req RfidScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	result, err := h.service.HandleRfidScan(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, result.Message, result)
}
