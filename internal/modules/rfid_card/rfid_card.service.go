package rfid_card

import (
	"errors"
	"log"
	"strings"
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

func (s *Service) validateCardType(cardType CardType) bool {
	return cardType == CardTypeRegistered || cardType == CardTypeGuest
}

func (s *Service) toResponse(card *RfidCard) *RfidCardResponse {
	return &RfidCardResponse{
		ID:        card.ID,
		UID:       card.UID,
		CardType:  card.CardType,
		UserID:    card.UserID,
		IsActive:  card.IsActive,
		CreatedAt: card.CreatedAt.Format(time.RFC3339),
		UpdatedAt: card.UpdatedAt.Format(time.RFC3339),
	}
}

func (s *Service) Create(req CreateRfidCardRequest) (*RfidCardResponse, error) {
	req.UID = strings.TrimSpace(req.UID)

	if req.UID == "" {
		return nil, appErrors.NewBadRequest("UID không được để trống")
	}

	if !s.validateCardType(req.CardType) {
		return nil, appErrors.NewBadRequest("CardType phải là REGISTERED hoặc GUEST")
	}

	if req.CardType == CardTypeGuest && req.UserID != nil {
		return nil, appErrors.NewBadRequest("Thẻ GUEST không được gán user")
	}

	card := &RfidCard{
		UID:      req.UID,
		CardType: req.CardType,
		UserID:   req.UserID,
	}

	log.Println("FE gửi lên:", card)

	if err := s.repo.Create(card); err != nil {
		return nil, appErrors.NewInternal("Tạo thẻ RFID thất bại")
	}

	return s.toResponse(card), nil
}

func (s *Service) Update(id uint64, req UpdateRfidCardRequest) (*RfidCardResponse, error) {
	currentCard, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không tìm thấy thẻ RFID")
		}

		return nil, appErrors.NewInternal("Lấy thông tin thẻ RFID thất bại")
	}

	data := map[string]any{}

	nextCardType := currentCard.CardType

	if req.CardType != nil {
		if !s.validateCardType(*req.CardType) {
			return nil, appErrors.NewBadRequest("CardType phải là REGISTERED hoặc GUEST")
		}

		nextCardType = *req.CardType
		data["card_type"] = *req.CardType
	}

	if nextCardType == CardTypeGuest {
		data["user_id"] = nil
	} else {
		if req.UserID == nil {
			return nil, appErrors.NewBadRequest("Thẻ REGISTERED bắt buộc phải có user_id")
		}

		data["user_id"] = *req.UserID
	}

	if req.IsActive != nil {
		data["is_active"] = *req.IsActive
	}

	if len(data) == 0 {
		return nil, appErrors.NewBadRequest("Không có dữ liệu để cập nhật")
	}

	if err := s.repo.UpdateByID(id, data); err != nil {
		return nil, appErrors.NewInternal("Cập nhật thẻ RFID thất bại")
	}

	card, err := s.repo.FindByID(id)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy thông tin thẻ RFID thất bại")
	}

	return s.toResponse(card), nil
}

func (s *Service) FindByUID(uid string) (*RfidCard, error) {
	uid = strings.TrimSpace(uid)

	if uid == "" {
		return nil, appErrors.NewBadRequest("UID không được để trống")
	}

	card, err := s.repo.FindByUID(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không tìm thấy thẻ RFID")
		}

		return nil, appErrors.NewInternal("Lấy thông tin thẻ RFID thất bại")
	}

	return card, nil
}

func (s *Service) GetStatistics() (*RfidCardStatisticsResponse, error) {
	total, registered, unregistered, active, err := s.repo.CountStatistics()
	if err != nil {
		return nil, appErrors.NewInternal("Thống kê thẻ RFID thất bại")
	}

	return &RfidCardStatisticsResponse{
		TotalCards:        total,
		RegisteredCards:   registered,
		UnregisteredCards: unregistered,
		ActiveCards:       active,
	}, nil
}

func (s *Service) GetByUserID(userID uint64) (*MyRfidCardResponse, error) {
	row, err := s.repo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, appErrors.NewNotFound("Không tìm thấy thẻ RFID của user")
		}

		return nil, appErrors.NewInternal("Lấy thẻ RFID của user thất bại")
	}

	var registeredAt *string
	if row.CardType == CardTypeRegistered {
		value := row.CreatedAt.Format("2006-01-02")
		registeredAt = &value
	}

	return &MyRfidCardResponse{
		ID:           row.ID,
		CardUID:      row.UID,
		UserID:       row.UserID,
		OwnerName:    row.OwnerName,
		Status:       row.CardType,
		IsActive:     row.IsActive,
		RegisteredAt: registeredAt,
	}, nil
}

func (s *Service) GetUserIDByEmail(email string) (*uint64, error) {
	userID, err := s.repo.GetUserIDByEmail(email)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy ID người dùng thất bại")
	}
	return userID, nil
}

func (s *Service) FindWithFilters(
	status string,
	keyword string,
	page int,
	pageSize int,
) (*RfidCardListResponse, error) {

	rows, total, err := s.repo.FindWithFilters(status, keyword, page, pageSize)
	if err != nil {
		return nil, appErrors.NewInternal("Lấy danh sách thẻ RFID thất bại")
	}

	items := make([]RfidCardListItem, 0, len(rows))

	for _, row := range rows {
		var registeredAt *string

		if row.CardType == CardTypeRegistered {
			value := row.CreatedAt.Format("2006-01-02")
			registeredAt = &value
		}

		items = append(items, RfidCardListItem{
			ID:           row.ID,
			CardUID:      row.UID,
			UserID:       row.UserID,
			OwnerName:    row.OwnerName,
			Status:       row.CardType,
			IsActive:     row.IsActive,
			RegisteredAt: registeredAt,
		})
	}

	meta := buildRfidCardListMeta(total, page, pageSize)

	return &RfidCardListResponse{
		Data: items,
		Meta: meta,
	}, nil
}

func buildRfidCardListMeta(totalElements int64, page int, pageSize int) RfidCardListMeta {
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

	return RfidCardListMeta{
		TotalElements: totalElements,
		TotalPages:    totalPages,
		CurrentPage:   page,
		PageSize:      pageSize,
	}
}
