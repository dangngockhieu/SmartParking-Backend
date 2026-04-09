package parking_lot

import (
	"errors"
	"strings"

	appErrors "backend/internal/common/errors"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Tạo mới thẻ đỗ xe
func (s *Service) Create(req CreateParkingLotRequest) (*ParkingLotResponse, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Location = strings.TrimSpace(req.Location)

	if req.Name == "" {
		return nil, appErrors.NewBadRequest("Name không được để trống")
	}
	if req.Location == "" {
		return nil, appErrors.NewBadRequest("Location không được để trống")
	}

	lot := &ParkingLot{
		Name:     req.Name,
		Location: &req.Location,
	}

	if err := s.repo.Create(lot); err != nil {
		return nil, appErrors.NewInternal("Tạo bãi đỗ thất bại")
	}

	return &ParkingLotResponse{
		ID:       lot.ID,
		Name:     lot.Name,
		Location: lot.Location,
	}, nil
}

// Lấy danh sách tất cả thẻ đỗ xe
func (s *Service) FindAll() ([]ParkingLotResponse, error) {
	lots, err := s.repo.FindAll()
	if err != nil {
		return nil, appErrors.NewInternal("Lấy danh sách bãi đỗ thất bại")
	}

	result := make([]ParkingLotResponse, 0, len(lots))
	for _, lot := range lots {
		result = append(result, ParkingLotResponse{
			ID:       lot.ID,
			Name:     lot.Name,
			Location: lot.Location,
		})
	}

	return result, nil
}

// Lấy thông tin chi tiết bãi đỗ theo ID, bao gồm danh sách slot và thống kê
func (s *Service) FindByID(id uint64) (*ParkingLotDetailResponse, error) {
	lot, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Parking lot not found")
		}
		return nil, appErrors.NewInternal("Lấy thông tin bãi đỗ thất bại")
	}

	slots, err := s.repo.FindSlotsByLotID(id)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy danh sách slot thất bại")
	}

	statRows, err := s.repo.CountStatsByLotID(id)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy thống kê bãi đỗ thất bại")
	}

	stats := ParkingLotStatsResponse{}
	for _, row := range statRows {
		stats.Total += row.Count

		switch row.Status {
		case "AVAILABLE":
			stats.Available = row.Count
		case "OCCUPIED":
			stats.Occupied = row.Count
		case "MAINTAIN":
			stats.Maintain = row.Count
		}
	}

	return &ParkingLotDetailResponse{
		ID:       lot.ID,
		Name:     lot.Name,
		Location: lot.Location,
		Slots:    slots,
		Stats:    stats,
	}, nil
}

// Lấy danh sách cổng của bãi đỗ theo ID
func (s *Service) FindGatesByLotID(lotID uint64) ([]ParkingLotGateResponse, error) {
	gates, err := s.repo.FindGatesByLotID(lotID)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy danh sách cổng thất bại")
	}
	return gates, nil
}

func (s *Service) Update(id uint64, req UpdateParkingLotRequest) (*ParkingLotResponse, error) {
	data := map[string]interface{}{}

	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return nil, appErrors.NewBadRequest("Name không được để trống")
		}
		data["name"] = name
	}

	if req.Location != nil {
		location := strings.TrimSpace(*req.Location)
		if location == "" {
			return nil, appErrors.NewBadRequest("Location không được để trống")
		}
		data["location"] = location
	}

	if len(data) == 0 {
		return nil, appErrors.NewBadRequest("Không có dữ liệu để cập nhật")
	}

	if err := s.repo.UpdateByID(id, data); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Parking lot not found")
		}
		return nil, appErrors.NewInternal("Cập nhật bãi đỗ thất bại")
	}

	lot, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Parking lot not found")
		}
		return nil, appErrors.NewInternal("Lấy thông tin bãi đỗ thất bại")
	}

	return &ParkingLotResponse{
		ID:       lot.ID,
		Name:     lot.Name,
		Location: lot.Location,
	}, nil
}
