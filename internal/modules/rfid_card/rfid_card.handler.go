package rfid_card

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
// @Summary Tạo thẻ RFID
// @Description Tạo mới thẻ RFID trong hệ thống. Thẻ GUEST không được gán user_id.
// @Tags rfid_card
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRfidCardRequest true "Thông tin thẻ RFID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /rfid-cards [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateRfidCardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	card, err := h.service.Create(req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusCreated, "Tạo thẻ RFID thành công", card)
}

// Update godoc
// @Summary Cập nhật thẻ RFID
// @Description Cập nhật thông tin thẻ RFID theo ID
// @Tags rfid_card
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID thẻ RFID"
// @Param request body UpdateRfidCardRequest true "Thông tin cập nhật"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /rfid-cards/{id} [patch]
func (h *Handler) Update(c *gin.Context) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.Error(appErrors.NewBadRequest("id không hợp lệ"))
		return
	}

	var req UpdateRfidCardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	card, err := h.service.Update(uint64(id64), req)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Cập nhật thẻ RFID thành công", card)
}

// GetStatistics godoc
// @Summary Thống kê thẻ RFID
// @Description Lấy thống kê tổng số thẻ, đã đăng ký, chưa đăng ký, đang hoạt động
// @Tags rfid_card
// @Produce json
// @Security BearerAuth
// @Param lotId query int false "ID bãi xe"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /rfid-cards/statistics [get]
func (h *Handler) GetStatistics(c *gin.Context) {

	result, err := h.service.GetStatistics()
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy thống kê thẻ RFID thành công", result)
}

// FindWithFilters godoc
// @Summary Lấy danh sách thẻ RFID
// @Description Lấy danh sách thẻ RFID có filter và phân trang
// @Tags rfid_card
// @Produce json
// @Security BearerAuth
// @Param lotId query int false "ID bãi xe, bỏ trống để lấy toàn bộ"
// @Param status query string false "REGISTERED hoặc GUEST"
// @Param keyword query string false "Từ khóa tìm kiếm theo UID, biển số hoặc tên chủ thẻ"
// @Param page query int false "Trang hiện tại" default(1)
// @Param pageSize query int false "Số lượng mỗi trang" default(10)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /rfid-cards [get]
func (h *Handler) FindWithFilters(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	keyword := c.Query("keyword")

	status := c.Query("status")

	result, err := h.service.FindWithFilters(status, keyword, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy danh sách thẻ RFID thành công", result)
}

// GetMyRfidCard godoc
// @Summary Lấy thẻ RFID của tôi
// @Description Lấy thông tin thẻ RFID đang gắn với tài khoản của người dùng hiện tại
// @Tags rfid_card
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /rfid-cards/my-card [get]
func (h *Handler) GetMyRfidCard(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.Error(appErrors.NewUnauthorized("Chưa đăng nhập"))
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		c.Error(appErrors.NewUnauthorized("userID không hợp lệ"))
		return
	}

	card, err := h.service.GetByUserID(userID)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy thẻ RFID của tôi thành công", card)
}

func (h *Handler) GetUserIDByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.Error(appErrors.NewBadRequest("Email là bắt buộc"))
		return
	}
	userID, err := h.service.GetUserIDByEmail(email)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "Lấy ID người dùng thành công", map[string]interface{}{"user_id": userID})
}
