package dashboard

import (
	"net/http"

	appErrors "backend/internal/common/errors"
	"backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// GetParkingFlow godoc
// @Summary Thống kê lưu lượng xe dashboard
// @Description Lấy thống kê xe vào/ra theo ngày, số xe hiện tại, tỉ lệ lấp đầy và giờ cao điểm
// @Tags dashboard
// @Security BearerAuth
// @Param date query string false "Ngày thống kê, format yyyy-mm-dd"
// @Param lotId query int false "ID bãi xe, bỏ trống để lấy toàn bộ bãi"
// @Success 200 {object} ParkingFlowResponse
// @Router /dashboard/parking-flow [get]
func (h *Handler) GetParkingFlow(c *gin.Context) {
	// Lấy dashboard cho bãi xe hoặc toàn bộ bãi nếu lotId không được cung cấp
	var query ParkingFlowQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(appErrors.NewBadRequest("Query không hợp lệ"))
		return
	}

	result, err := h.service.GetParkingFlow(query)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "Lấy thông tin lưu lượng xe thành công", result)
}

func (h *Handler) GetRevenueByMonth(c *gin.Context) {
	var query RevenueDateQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(appErrors.NewBadRequest("Query không hợp lệ"))
		return
	}

	result, err := h.service.CountMonthlyRevenue(query.LotID, query.Date)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "Lấy thông tin doanh thu theo tháng thành công", gin.H{"revenue": result})
}

func (h *Handler) GetRevenueByDay(c *gin.Context) {
	var query RevenueDateQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		_ = c.Error(appErrors.NewBadRequest("Query không hợp lệ"))
		return
	}

	result, err := h.service.CountDailyRevenue(query.LotID, query.Date)
	if err != nil {
		_ = c.Error(err)
		return
	}
	response.Success(c, http.StatusOK, "Lấy thông tin doanh thu theo ngày thành công", gin.H{"revenue": result})
}
