package parking_session

import (
	"net/http"
	"strconv"
	"time"

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

// FindAll godoc
// @Summary Lấy danh sách phiên gửi xe
// @Description Trả về danh sách tất cả phiên gửi xe
// @Tags parking_session
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /parking-sessions [get]
func (h *Handler) FindAll(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.Error(appErrors.NewBadRequest("date không được để trống"))
		return
	}

	search := c.Query("search")

	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		c.Error(appErrors.NewBadRequest("date không hợp lệ, định dạng đúng là YYYY-MM-DD"))
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.Error(appErrors.NewBadRequest("page không hợp lệ"))
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		c.Error(appErrors.NewBadRequest("pageSize không hợp lệ"))
		return
	}

	lotId, err := strconv.ParseUint(c.Query("lotId"), 10, 64)
	if err != nil && c.Query("lotId") != "" {
		c.Error(appErrors.NewBadRequest("lotId không hợp lệ"))
		return
	}

	result, err := h.service.FindAll(date, page, pageSize, search, lotId)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy danh sách phiên gửi xe theo ngày thành công", result)
}

// GetByDate godoc
// @Summary Lấy danh sách phiên gửi xe theo ngày
// @Description Lấy danh sách phiên gửi xe theo ngày, có phân trang
// @Tags parking_session
// @Produce json
// @Param date query string true "Ngày theo định dạng YYYY-MM-DD"
// @Param page query int false "Số trang, mặc định 1"
// @Param pageSize query int false "Số lượng mỗi trang, mặc định 10"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /parking-sessions/by-date [get]
func (h *Handler) GetByDate(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.Error(appErrors.NewBadRequest("date không được để trống"))
		return
	}

	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}

	userID, ok := userIDValue.(uint64)
	if !ok {
		c.Error(appErrors.NewUnauthorized("Unauthorized"))
		return
	}

	date, err := time.ParseInLocation("2006-01-02", dateStr, time.Local)
	if err != nil {
		c.Error(appErrors.NewBadRequest("date không hợp lệ, định dạng đúng là YYYY-MM-DD"))
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.Error(appErrors.NewBadRequest("page không hợp lệ"))
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		c.Error(appErrors.NewBadRequest("pageSize không hợp lệ"))
		return
	}

	result, err := h.service.GetByDate(date, userID, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy danh sách phiên gửi xe theo ngày thành công", result)
}
