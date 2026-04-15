package gate

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
// @Summary Tạo cổng mới
// @Tags gate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateGateRequest true "Thông tin cổng"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /gates [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateGateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest(err.Error()))
		return
	}

	gate, err := h.service.Create(&req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "Tạo cổng thành công", ToGateResponse(gate))
}

// Update godoc
// @Summary Cập nhật thông tin cổng
// @Tags gate
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param gateId path int true "ID cổng"
// @Param body body UpdateGateRequest true "Thông tin cập nhật"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /gates/{gateId} [put]
func (h *Handler) Update(c *gin.Context) {
	id, err := parseUintParam(c, "id")
	if err != nil {
		c.Error(appErrors.NewBadRequest("gateId không hợp lệ"))
		return
	}

	var req UpdateGateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest(err.Error()))
		return
	}

	gate, err := h.service.Update(id, &req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Cập nhật cổng thành công", ToGateResponse(gate))
}

// ─── helper ──────────────────────────────────────────────────────────────────

func parseUintParam(c *gin.Context, key string) (uint64, error) {
	v, err := strconv.ParseUint(c.Param(key), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint64(v), nil
}
