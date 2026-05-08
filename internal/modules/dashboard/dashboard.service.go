package dashboard

import (
	"errors"
	"fmt"
	"math"
	"time"

	appErrors "backend/internal/common/errors"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// GetParkingFlow lấy thông tin lưu lượng xe vào/ra, số xe hiện tại
// Tỉ lệ lấp đầy và giờ cao điểm cho bãi xe hoặc toàn bộ bãi nếu lotId không được cung cấp
func (s *Service) GetParkingFlow(query ParkingFlowQuery) (*ParkingFlowResponse, error) {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		location = time.Local
	}

	if query.Date == "" {
		query.Date = time.Now().In(location).Format("2006-01-02")
	}

	parsedDate, err := time.ParseInLocation("2006-01-02", query.Date, location)
	if err != nil {
		return nil, appErrors.NewBadRequest("Date không hợp lệ, định dạng đúng là yyyy-mm-dd")
	}

	startOfDay := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		0,
		0,
		0,
		0,
		location,
	)

	endOfDay := startOfDay.AddDate(0, 0, 1)

	lotName := "Toàn bộ bãi"

	if query.LotID != nil {
		name, err := s.repo.GetLotName(*query.LotID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, appErrors.NewNotFound("Không tìm thấy bãi xe")
			}

			return nil, appErrors.NewInternal("Lấy thông tin bãi xe thất bại")
		}

		if name == "" {
			return nil, appErrors.NewNotFound("Không tìm thấy bãi xe")
		}

		lotName = name
	}

	todayIn, err := s.repo.CountTodayIn(startOfDay, endOfDay, query.LotID)
	if err != nil {
		return nil, appErrors.NewInternal("Thống kê số xe vào thất bại")
	}

	todayOut, err := s.repo.CountTodayOut(startOfDay, endOfDay, query.LotID)
	if err != nil {
		return nil, appErrors.NewInternal("Thống kê số xe ra thất bại")
	}

	currentVehicles, err := s.repo.CountCurrentVehicles(query.LotID)
	if err != nil {
		return nil, appErrors.NewInternal("Thống kê số xe hiện tại thất bại")
	}

	capacity, err := s.repo.CountCapacity(query.LotID)
	if err != nil {
		return nil, appErrors.NewInternal("Thống kê sức chứa bãi xe thất bại")
	}

	availableSlots, err := s.repo.CountAvailableSlots(query.LotID)
	if err != nil {
		return nil, appErrors.NewInternal("Thống kê số chỗ trống thất bại")
	}

	hourlyInRows, err := s.repo.GetHourlyIn(startOfDay, endOfDay, query.LotID)
	if err != nil {
		return nil, appErrors.NewInternal("Thống kê xe vào theo giờ thất bại")
	}

	hourlyOutRows, err := s.repo.GetHourlyOut(startOfDay, endOfDay, query.LotID)
	if err != nil {
		return nil, appErrors.NewInternal("Thống kê xe ra theo giờ thất bại")
	}

	hourlyFlow, peakTime := s.buildHourlyFlow(hourlyInRows, hourlyOutRows)

	occupancyRate := s.calculateOccupancyRate(currentVehicles, capacity)
	occupancyStatus := s.getOccupancyStatus(occupancyRate)
	statusMessage := s.getOccupancyStatusMessage(occupancyStatus)

	return &ParkingFlowResponse{
		Date:    query.Date,
		LotID:   query.LotID,
		LotName: lotName,
		Summary: DashboardSummaryResponse{
			TodayIn:         todayIn,
			TodayOut:        todayOut,
			CurrentVehicles: currentVehicles,
			Capacity:        capacity,
			AvailableSlots:  availableSlots,
			OccupancyRate:   occupancyRate,
		},
		HourlyFlow: hourlyFlow,
		Insights: DashboardInsightsResponse{
			PeakTime:        peakTime,
			OccupancyStatus: occupancyStatus,
			StatusMessage:   statusMessage,
		},
	}, nil
}

func (s *Service) buildHourlyFlow(
	hourlyInRows []HourlyCountRow,
	hourlyOutRows []HourlyCountRow,
) ([]HourlyFlowResponse, PeakTimeResponse) {
	inMap := make(map[int]int64)
	outMap := make(map[int]int64)

	for _, row := range hourlyInRows {
		inMap[row.Hour] = row.Count
	}

	for _, row := range hourlyOutRows {
		outMap[row.Hour] = row.Count
	}

	result := make([]HourlyFlowResponse, 0, 24)

	peakHour := 0
	var peakTotal int64 = -1

	for hour := 0; hour < 24; hour++ {
		inCount := inMap[hour]
		outCount := outMap[hour]
		total := inCount + outCount

		if total > peakTotal {
			peakTotal = total
			peakHour = hour
		}

		result = append(result, HourlyFlowResponse{
			Hour: fmt.Sprintf("%02d:00", hour),
			In:   inCount,
			Out:  outCount,
		})
	}

	return result, PeakTimeResponse{
		From: fmt.Sprintf("%02d:00", peakHour),
		To:   fmt.Sprintf("%02d:00", (peakHour+1)%24),
	}
}

func (s *Service) calculateOccupancyRate(currentVehicles int64, capacity int64) float64 {
	if capacity <= 0 {
		return 0
	}

	rate := (float64(currentVehicles) / float64(capacity)) * 100

	return math.Round(rate*100) / 100
}

func (s *Service) getOccupancyStatus(rate float64) string {
	if rate >= 100 {
		return "FULL"
	}

	if rate >= 80 {
		return "HIGH"
	}

	if rate >= 50 {
		return "MEDIUM"
	}

	return "LOW"
}

func (s *Service) getOccupancyStatusMessage(status string) string {
	switch status {
	case "FULL":
		return "Bãi xe đã đầy"
	case "HIGH":
		return "Bãi đang gần đầy"
	case "MEDIUM":
		return "Bãi đang ở mức trung bình"
	default:
		return "Bãi còn nhiều chỗ trống"
	}
}

// Đếm doanh thu theo ngày cho bãi xe
func (s *Service) CountDailyRevenue(lotID *uint64, date string) (int64, error) {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		location = time.Local
	}

	parsedDate, err := time.ParseInLocation("2006-01-02", date, location)
	if err != nil {
		return 0, appErrors.NewBadRequest("Date không hợp lệ, định dạng đúng là yyyy-mm-dd")
	}

	startOfDay := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		parsedDate.Day(),
		0,
		0,
		0,
		0,
		location,
	)
	endOfDay := startOfDay.AddDate(0, 0, 1)

	total, err := s.repo.GetTotalRevenue(lotID, startOfDay, endOfDay)
	if err != nil {
		return 0, appErrors.NewInternal("Thống kê doanh thu thất bại")
	}
	return total, nil
}

// Đếm doanh thu theo tháng cho bãi xe
func (s *Service) CountMonthlyRevenue(lotID *uint64, date string) (int64, error) {
	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		location = time.Local
	}

	parsedDate, err := time.ParseInLocation("2006-01-02", date, location)
	if err != nil {
		return 0, appErrors.NewBadRequest("Date không hợp lệ, định dạng đúng là yyyy-mm-dd")
	}

	startOfMonth := time.Date(
		parsedDate.Year(),
		parsedDate.Month(),
		1,
		0,
		0,
		0,
		0,
		location,
	)

	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	total, err := s.repo.GetTotalRevenue(lotID, startOfMonth, endOfMonth)
	if err != nil {
		return 0, appErrors.NewInternal("Thống kê doanh thu thất bại")
	}

	return total, nil
}
