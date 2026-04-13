package parking_session

import (
	"errors"
	"strings"
	"time"

	appErrors "backend/internal/common/errors"
	"backend/internal/modules/rfid_card"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) toResponse(session *ParkingSession) *ParkingSessionResponse {
	return &ParkingSessionResponse{
		ID:          session.ID,
		LotID:       session.LotID,
		SlotID:      session.SlotID,
		CardUID:     session.CardUID,
		CardType:    string(session.CardType),
		PlateNumber: session.PlateNumber,
		EntryTime:   session.EntryTime,
		ExitTime:    session.ExitTime,
		Fee:         session.Fee,
		IsActive:    session.IsActive,
	}
}

// Tạo phiên gửi xe mới
func (s *Service) Create(input CreateParkingSessionInput) (*ParkingSession, error) {
	input.CardUID = strings.TrimSpace(input.CardUID)
	input.PlateNumber = strings.TrimSpace(input.PlateNumber)

	if input.CardUID == "" {
		return nil, appErrors.NewBadRequest("CardUID không được để trống")
	}
	if input.PlateNumber == "" {
		return nil, appErrors.NewBadRequest("PlateNumber không được để trống")
	}
	if input.CardType != string(rfid_card.CardTypeRegistered) && input.CardType != string(rfid_card.CardTypeGuest) {
		return nil, appErrors.NewBadRequest("CardType không hợp lệ")
	}

	session := &ParkingSession{
		LotID:       input.LotID,
		CardUID:     input.CardUID,
		CardType:    rfid_card.CardType(input.CardType),
		PlateNumber: input.PlateNumber,
		IsActive:    true,
	}

	if err := s.repo.Create(session); err != nil {
		return nil, appErrors.NewInternal("Tạo phiên gửi xe thất bại")
	}

	return session, nil
}

// Gán slot cho phiên gửi xe
func (s *Service) AssignSlot(input AssignSlotInput) (*ParkingSession, error) {
	session, err := s.repo.FindByID(input.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không tìm thấy phiên gửi xe")
		}
		return nil, appErrors.NewInternal("Lấy phiên gửi xe thất bại")
	}

	if !session.IsActive {
		return nil, appErrors.NewBadRequest("Phiên gửi xe đã đóng")
	}

	if err := s.repo.UpdateByID(session.ID, map[string]any{
		"slot_id": input.SlotID,
	}); err != nil {
		return nil, appErrors.NewInternal("Gán slot cho phiên gửi xe thất bại")
	}

	return s.repo.FindByID(session.ID)
}

// Kết thúc phiên gửi xe
func (s *Service) FinishSession(input FinishParkingSessionInput) (*ParkingSession, error) {
	session, err := s.repo.FindByID(input.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không tìm thấy phiên gửi xe")
		}
		return nil, appErrors.NewInternal("Lấy phiên gửi xe thất bại")
	}

	if !session.IsActive {
		return nil, appErrors.NewBadRequest("Phiên gửi xe đã đóng")
	}

	now := time.Now()

	if err := s.repo.UpdateByID(session.ID, map[string]any{
		"exit_time": now,
		"fee":       input.Fee,
		"is_active": false,
	}); err != nil {
		return nil, appErrors.NewInternal("Kết thúc phiên gửi xe thất bại")
	}

	return s.repo.FindByID(session.ID)
}

// Lấy danh sách phiên gửi xe theo ngày
func (s *Service) FindAll(
	date time.Time,
	page int,
	pageSize int,
	search string,
	lotId uint64,
) (*ManageParkingSessionListResponse, error) {
	sessions, total, err := s.repo.FindAll(date, page, pageSize, search, lotId)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy danh sách phiên gửi xe thất bại")
	}

	items := make([]ManageParkingSessionResponse, 0, len(sessions))
	for _, session := range sessions {
		items = append(items, session)
	}

	meta := buildParkingSessionListMeta(total, page, pageSize)

	return &ManageParkingSessionListResponse{
		Data: items,
		Meta: meta,
	}, nil
}

// Lấy phiên gửi xe đang hoạt động theo CardUID
func (s *Service) FindActiveByCardUID(cardUID string) (*ParkingSession, error) {
	session, err := s.repo.FindActiveByCardUID(strings.TrimSpace(cardUID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không có phiên gửi xe đang hoạt động")
		}
		return nil, appErrors.NewInternal("Lấy phiên gửi xe đang hoạt động thất bại")
	}
	return session, nil
}

// Lấy phiên gửi xe đang hoạt động theo PlateNumber
func (s *Service) FindActiveByPlateNumber(plateNumber string) (*ParkingSession, error) {
	session, err := s.repo.FindActiveByPlateNumber(strings.TrimSpace(plateNumber))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("No active session with this plate")
		}
		return nil, appErrors.NewInternal("Failed to check plate session")
	}
	return session, nil
}

// Lấy danh sách phiên gửi xe theo ngày và userID
func (s *Service) GetByDate(
	date time.Time,
	userID uint64,
	page int,
	pageSize int,
) (*ParkingSessionListResponse, error) {
	sessions, total, err := s.repo.FindByDate(date, userID, page, pageSize)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy danh sách phiên gửi xe thất bại")
	}

	items := make([]ParkingSessionResponse, 0, len(sessions))

	for _, session := range sessions {
		items = append(items, *s.toResponse(&session))
	}

	meta := buildParkingSessionListMeta(total, page, pageSize)

	return &ParkingSessionListResponse{
		Data: items,
		Meta: meta,
	}, nil
}

func buildParkingSessionListMeta(
	totalElements int64,
	page int,
	pageSize int,
) ParkingSessionListMeta {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	totalPages := 0

	if totalElements > 0 {
		totalPages = int((totalElements + int64(pageSize) - 1) / int64(pageSize))
	}

	return ParkingSessionListMeta{
		TotalElements: totalElements,
		TotalPages:    totalPages,
		CurrentPage:   page,
		PageSize:      pageSize,
	}
}
