package wallet

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	appErrors "backend/internal/common/errors"
	"backend/pkg/response"
)

type Handler struct {
	service   *Service
	cancelURL string
}

func NewHandler(service *Service, cancelURL string) *Handler {
	return &Handler{service: service, cancelURL: cancelURL}
}

// CreateDeposit tạo link nạp tiền
// @Summary Tạo link nạp tiền
// @Description Tạo link nạp tiền qua PayOS
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateDepositRequest true "Thông tin nạp tiền"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /wallet/deposit [post]
func (h *Handler) CreateDeposit(c *gin.Context) {
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

	var req CreateDepositRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Dữ liệu không hợp lệ"))
		return
	}

	data, err := h.service.CreateDeposit(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}
	response.Success(c, http.StatusCreated, "Tạo link nạp tiền thành công", data)
}

// PayOSWebhook xử lý webhook từ PayOS
// @Summary Xử lý webhook từ PayOS
// @Description Xử lý webhook từ PayOS để cập nhật trạng thái giao dịch
// @Tags wallet
// @Accept json
// @Produce json
// @Param request body PayOSWebhookRequest true "Thông tin webhook"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /wallet/webhook [post]
func (h *Handler) PayOSWebhook(c *gin.Context) {
	var req PayOSWebhookRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(appErrors.NewBadRequest("Webhook không hợp lệ"))
		return
	}

	err := h.service.HandlePayOSWebhook(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidWebhook) {
			c.Error(appErrors.NewBadRequest("Webhook không hợp lệ"))
			return
		}
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Webhook processed", nil)
}

// GetMyTransactions godoc
// @Summary Lấy lịch sử giao dịch của người dùng
// @Description Lấy lịch sử giao dịch của người dùng với pagination kiểu cursor
// @Tags wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Số lượng giao dịch trả về, mặc định 20"
// @Param cursorCreatedAt query string false "Thời gian tạo của giao dịch cuối cùng ở trang trước, dùng để phân trang kiểu cursor"
// @Param cursorId query uint64 false "ID của giao dịch cuối cùng ở trang trước, dùng để phân trang kiểu cursor"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /wallet/transactions [get]
func (h *Handler) GetMyTransactions(c *gin.Context) {
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

	limit := 20

	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsed, err := strconv.Atoi(rawLimit)
		if err == nil {
			limit = parsed
		}
	}

	var cursorCreatedAt *time.Time
	var cursorID *uint64

	if rawCreatedAt := c.Query("cursorCreatedAt"); rawCreatedAt != "" {
		parsed, err := time.Parse(time.RFC3339, rawCreatedAt)
		if err != nil {
			c.Error(appErrors.NewBadRequest("cursorCreatedAt không hợp lệ, cần dạng RFC3339"))
			return
		}

		cursorCreatedAt = &parsed
	}

	if rawID := c.Query("cursorId"); rawID != "" {
		parsed, err := strconv.ParseUint(rawID, 10, 64)
		if err != nil {
			c.Error(appErrors.NewBadRequest("cursorId không hợp lệ"))
			return
		}

		cursorID = &parsed
	}

	if (cursorCreatedAt == nil && cursorID != nil) || (cursorCreatedAt != nil && cursorID == nil) {
		c.Error(appErrors.NewBadRequest("Cần truyền đủ cursorCreatedAt và cursorId"))
		return
	}

	data, err := h.service.GetMyTransactions(
		c.Request.Context(),
		userID,
		cursorCreatedAt,
		cursorID,
		limit,
	)

	if err != nil {
		c.Error(err)
		return
	}

	response.Success(c, http.StatusOK, "Lấy lịch sử ví thành công", data)
}

// UpdateWalletCancel xử lý callback khi người dùng hủy giao dịch nạp tiền
// @Summary Xử lý callback khi người dùng hủy giao dịch nạp tiền
// @Description Xử lý callback khi người dùng hủy giao dịch nạp tiền để cập nhật trạng thái giao dịch
// @Tags wallet
// @Accept json
// @Produce json
// @Param orderCode query uint64 true "Mã đơn hàng cần hủy"
// @Success 302 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /wallet/deposit/cancel-return [get]
func (h *Handler) UpdateWalletCancel(c *gin.Context) {
	orderCodeStr := c.Query("orderCode")
	if orderCodeStr == "" {
		orderCodeStr = c.Query("order_code")
	}

	if orderCodeStr == "" {
		c.Redirect(
			http.StatusFound,
			h.cancelURL+"?message=missing_order_code",
		)
		return
	}

	orderCode, err := strconv.ParseUint(orderCodeStr, 10, 64)
	if err != nil || orderCode == 0 {
		c.Redirect(
			http.StatusFound,
			h.cancelURL+"?message=invalid_order_code",
		)
		return
	}

	status := c.Query("status")
	cancel := c.Query("cancel")

	isCanceled := cancel == "true" ||
		status == "CANCELLED" ||
		status == "CANCELED"

	if isCanceled {
		if err := h.service.UpdateWalletCancel(c.Request.Context(), orderCode); err != nil {
			c.Redirect(
				http.StatusFound,
				h.cancelURL+"?orderCode="+orderCodeStr+"&message=cancel_failed",
			)
			return
		}
	}

	c.Redirect(
		http.StatusFound,
		h.cancelURL+"?orderCode="+orderCodeStr+"&status=CANCELED",
	)
}
