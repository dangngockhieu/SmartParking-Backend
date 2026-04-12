package parking_slot

import (
	"errors"

	appErrors "backend/internal/common/errors"
	"backend/internal/realtime/parking"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
	hub  *parking.Hub
}

func NewService(repo *Repository, hub *parking.Hub) *Service {
	return &Service{
		repo: repo,
		hub:  hub,
	}
}

func (s *Service) updateStatus(slot *ParkingSlot, newStatus SlotStatus) (*UpdateParkingSlotResponse, error) {
	if slot.Status == newStatus {
		return &UpdateParkingSlotResponse{
			Changed:   false,
			ID:        slot.ID,
			LotID:     slot.LotID,
			Name:      slot.Name,
			Message:   "Trạng thái không thay đổi",
			OldStatus: slot.Status,
			NewStatus: slot.Status,
		}, nil
	}

	oldStatus := slot.Status

	if err := s.repo.UpdateStatus(slot.ID, newStatus); err != nil {
		return nil, appErrors.NewInternal("Cập nhật trạng thái thất bại")
	}

	result := &UpdateParkingSlotResponse{
		Changed:   true,
		ID:        slot.ID,
		LotID:     slot.LotID,
		Name:      slot.Name,
		Message:   "Cập nhật trạng thái thành công",
		OldStatus: oldStatus,
		NewStatus: newStatus,
	}

	s.hub.BroadcastToLot(slot.LotID, "SLOT_STATUS_CHANGE", result)

	return result, nil
}

func (s *Service) Create(req CreateParkingSlotRequest) (*ParkingSlot, error) {
	slot := &ParkingSlot{
		Name:       req.Name,
		LotID:      req.LotID,
		DeviceMac:  req.DeviceMac,
		PortNumber: req.PortNumber,
		Status:     SlotStatusAvailable,
	}

	if err := s.repo.Create(slot); err != nil {
		return nil, appErrors.NewInternal("Tạo vị trí đỗ thất bại")
	}
	return slot, nil
}

func (s *Service) FindByID(id uint64) (*ParkingSlot, error) {
	slot, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không tìm thấy vị trí đỗ")
		}
		return nil, appErrors.NewInternal("Lấy thông tin vị trí đỗ thất bại")
	}
	return slot, nil
}

func (s *Service) AdminUpdateStatus(id uint64, req AdminUpdateParkingSlotRequest) (*UpdateParkingSlotResponse, error) {
	slot, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không tìm thấy vị trí đỗ")
		}
		return nil, appErrors.NewInternal("Lấy vị trí đỗ thất bại")
	}
	return s.updateStatus(slot, req.Status)
}

func (s *Service) SensorUpdateStatus(req SensorUpdateParkingSlotRequest) (*UpdateParkingSlotResponse, error) {
	slot, err := s.repo.FindByDeviceMacAndPort(req.Mac, req.Port)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không tìm thấy vị trí đỗ")
		}
		return nil, appErrors.NewInternal("Lấy vị trí đỗ thất bại")
	}

	if slot.Status == SlotStatusMaintain {
		return &UpdateParkingSlotResponse{
			Changed:   false,
			ID:        slot.ID,
			LotID:     slot.LotID,
			Name:      slot.Name,
			Message:   "Ô đang bảo trì",
			OldStatus: slot.Status,
			NewStatus: slot.Status,
		}, nil
	}

	newStatus := SlotStatusAvailable
	if req.IsOccupied != nil && *req.IsOccupied {
		newStatus = SlotStatusOccupied
	}

	return s.updateStatus(slot, newStatus)
}

func (s *Service) ChangeDevice(id uint64, req ChangeSlotDeviceRequest) (*ParkingSlot, error) {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không tìm thấy vị trí đỗ")
		}
		return nil, appErrors.NewInternal("Lấy vị trí đỗ thất bại")
	}

	conflict, err := s.repo.FindConflictSlot(req.DeviceMac, req.PortNumber, id)
	if err == nil && conflict != nil {
		return nil, appErrors.NewBadRequest("Device + port already assigned to another slot")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, appErrors.NewInternal("Kiểm tra trùng thiết bị thất bại")
	}

	updated, err := s.repo.ChangeDevice(id, req.DeviceMac, req.PortNumber)
	if err != nil {
		return nil, appErrors.NewInternal("Cập nhật thiết bị thất bại")
	}

	return updated, nil
}

func (s *Service) IsAvailable(lotID uint64) (bool, error) {
	exists, err := s.repo.IsAvailable(lotID)
	if err != nil {
		return false, appErrors.NewInternal("Kiểm tra vị trí đỗ thất bại")
	}
	return exists, nil
}
