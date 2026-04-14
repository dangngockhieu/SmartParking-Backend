package gate

import (
	appErrors "backend/internal/common/errors"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Create tạo mới một cổng
func (s *Service) Create(req *CreateGateRequest) (*Gate, error) {
	gate := &Gate{
		Name:       req.Name,
		Type:       req.Type,
		MacAddress: req.MacAddress,
		LotID:      req.LotID,
	}
	if err := s.repo.Create(gate); err != nil {
		return nil, appErrors.NewInternal("Không thể tạo cổng")
	}
	return gate, nil
}

// FindByID lấy thông tin cổng theo ID
func (s *Service) FindByID(id uint64) (*Gate, error) {
	gate, err := s.repo.FindByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, appErrors.NewNotFound("Không tìm thấy cổng")
		}
		return nil, appErrors.NewInternal("Không thể lấy thông tin cổng")
	}
	return gate, nil
}

// Update cập nhật thông tin cổng
func (s *Service) Update(id uint64, req *UpdateGateRequest) (*Gate, error) {
	gate, err := s.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.MacAddress != "" {
		gate.MacAddress = req.MacAddress
	}

	if err := s.repo.Update(gate); err != nil {
		return nil, appErrors.NewInternal("Không thể cập nhật cổng")
	}
	return gate, nil
}
