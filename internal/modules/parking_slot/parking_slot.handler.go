package parking_slot

import (
	"net/http"
	"strconv"

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

// Create godoc
// @Summary Tạo vị trí đỗ xe
// @Description Tạo mới một vị trí đỗ xe trong hệ thống
// @Tags parking_slot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateParkingSlotRequest true "Thông tin vị trí đỗ xe"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /parking-slots [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateParkingSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	slot, err := h.service.Create(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "Tạo vị trí đỗ thành công", slot)
}

// FindByID godoc
// @Summary Lấy chi tiết vị trí đỗ xe
// @Description Lấy thông tin chi tiết của một vị trí đỗ xe theo ID
// @Tags parking_slot
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID vị trí đỗ xe"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /parking-slots/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequest("id không hợp lệ"))
		return
	}

	slot, err := h.service.FindByID(uint64(id64))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy thông tin vị trí đỗ thành công", slot)
}

// AdminUpdateStatus godoc
// @Summary Quản trị viên cập nhật trạng thái vị trí đỗ
// @Description Cập nhật trạng thái vị trí đỗ xe theo ID bởi quản trị viên hoặc quản lý
// @Tags parking_slot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID vị trí đỗ xe"
// @Param request body AdminUpdateParkingSlotRequest true "Thông tin cập nhật trạng thái"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /parking-slots/admin/{id} [patch]
func (h *Handler) AdminUpdateStatus(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequest("id không hợp lệ"))
		return
	}

	var req AdminUpdateParkingSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	result, err := h.service.AdminUpdateStatus(uint64(id64), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, result.Message, result)
}

// SensorUpdateStatus godoc
// @Summary Thiết bị cảm biến cập nhật trạng thái vị trí đỗ
// @Description Cập nhật trạng thái vị trí đỗ xe từ cảm biến hoặc thiết bị IoT
// @Tags parking_slot
// @Accept json
// @Produce json
// @Param request body SensorUpdateParkingSlotRequest true "Dữ liệu trạng thái từ cảm biến"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /parking-slots/sensor [post]
func (h *Handler) SensorUpdateStatus(c *gin.Context) {
	var req SensorUpdateParkingSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	result, err := h.service.SensorUpdateStatus(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, result.Message, result)
}

// ChangeDevice godoc
// @Summary Đổi thiết bị cho vị trí đỗ
// @Description Cập nhật hoặc thay đổi thiết bị IoT gắn với vị trí đỗ xe theo ID
// @Tags parking_slot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID vị trí đỗ xe"
// @Param request body ChangeSlotDeviceRequest true "Thông tin thiết bị mới"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /parking-slots/{id}/device [patch]
func (h *Handler) ChangeDevice(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequest("id không hợp lệ"))
		return
	}

	var req ChangeSlotDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	_, err = h.service.ChangeDevice(uint64(id64), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Cập nhật thiết bị thành công", nil)
}
