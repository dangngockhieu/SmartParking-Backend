package parking_lot

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
// @Summary Tạo bãi đỗ xe
// @Description Tạo mới một bãi đỗ xe trong hệ thống
// @Tags parking_lot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateParkingLotRequest true "Thông tin bãi đỗ xe"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /parking-lots [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateParkingLotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	lot, err := h.service.Create(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "Tạo bãi đỗ thành công", lot)
}

// FindAll godoc
// @Summary Lấy danh sách bãi đỗ xe
// @Description Trả về toàn bộ danh sách bãi đỗ xe
// @Tags parking_lot
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /parking-lots [get]
func (h *Handler) FindAll(c *gin.Context) {
	lots, err := h.service.FindAll()
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy danh sách bãi đỗ thành công", lots)
}

// FindByID godoc
// @Summary Lấy chi tiết bãi đỗ xe
// @Description Lấy thông tin chi tiết của một bãi đỗ xe theo ID
// @Tags parking_lot
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID bãi đỗ xe"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /parking-lots/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequest("id không hợp lệ"))
		return
	}

	lot, err := h.service.FindByID(uint64(id64))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy thông tin bãi đỗ thành công", lot)
}

// GetGatesByLotID godoc
// @Summary Lấy danh sách cổng theo bãi đỗ
// @Tags parking_lot
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID bãi đỗ xe"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /parking-lots/{id}/gates [get]
func (h *Handler) GetGatesByLotID(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequest("id không hợp lệ"))
		return
	}

	gates, err := h.service.FindGatesByLotID(uint64(id64))
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy danh sách cổng thành công", gates)
}

// Update godoc
// @Summary Cập nhật bãi đỗ xe
// @Description Cập nhật thông tin bãi đỗ xe theo ID
// @Tags parking_lot
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID bãi đỗ xe"
// @Param request body UpdateParkingLotRequest true "Thông tin cập nhật bãi đỗ xe"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /parking-lots/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequest("id không hợp lệ"))
		return
	}

	var req UpdateParkingLotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	lot, err := h.service.Update(uint64(id64), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Cập nhật bãi đỗ thành công", lot)
}
