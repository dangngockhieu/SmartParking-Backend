package parking_session

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Tạo phiên gửi xe mới
func (r *Repository) Create(session *ParkingSession) error {
	return r.db.Create(session).Error
}

// Lấy phiên gửi xe theo ID
func (r *Repository) FindByID(id uint64) (*ParkingSession, error) {
	var session ParkingSession
	err := r.db.First(&session, id).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Lấy danh sách phiên gửi xe theo ngày, phân trang và tìm kiếm
func (r *Repository) FindAll(
	date time.Time,
	page int,
	pageSize int,
	search string,
	lotId uint64,
) ([]ManageParkingSessionResponse, int64, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	start := time.Date(
		date.Year(), date.Month(), date.Day(),
		0, 0, 0, 0,
		date.Location(),
	)
	end := start.Add(24 * time.Hour)

	if search != "" {
		search = search + "%"
	}

	// 1. Khởi tạo query cơ bản (CHƯA có Select)
	base := r.db.Model(&ParkingSession{}).
		Joins("JOIN rfid_cards ON rfid_cards.uid = parking_sessions.card_uid").
		Joins("LEFT JOIN users ON users.id = rfid_cards.user_id").
		Where(
			"parking_sessions.entry_time >= ? AND parking_sessions.entry_time < ?",
			start,
			end,
		).
		Where("parking_sessions.lot_id = ?", lotId)

	// 2. Bổ sung điều kiện tìm kiếm nếu có
	if search != "" {
		base = base.Where(
			"parking_sessions.plate_number LIKE ? OR parking_sessions.card_uid LIKE ?",
			search, search,
		)
	}

	// 3. Count tổng số bản ghi TRƯỚC KHI Select
	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 4. Thực hiện Select, Order, Limit, Offset và Find
	sessions := make([]ManageParkingSessionResponse, 0)
	err := base.Select(`
			parking_sessions.id,
			parking_sessions.lot_id,
			parking_sessions.slot_id,
			parking_sessions.card_uid,
			parking_sessions.card_type,
			parking_sessions.plate_number,
			parking_sessions.entry_time,
			parking_sessions.exit_time,
			parking_sessions.fee,
			parking_sessions.is_active,
			COALESCE(CONCAT(users.last_name, ' ', users.first_name), 'Khách vãng lai') AS owner_name
		`).
		Order("exit_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&sessions).Error

	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}

// Gán slot cho phiên gửi xe
func (r *Repository) FindActiveByCardUID(cardUID string) (*ParkingSession, error) {
	var session ParkingSession
	err := r.db.
		Where("card_uid = ? AND is_active = ?", cardUID, true).
		Order("id DESC").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Gán slot cho phiên gửi xe
func (r *Repository) FindActiveByPlateNumber(plateNumber string) (*ParkingSession, error) {
	var session ParkingSession
	err := r.db.
		Where("plate_number = ? AND is_active = ?", plateNumber, true).
		Order("id DESC").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Cập nhật phiên gửi xe theo ID
func (r *Repository) UpdateByID(id uint64, data map[string]any) error {
	return r.db.Model(&ParkingSession{}).Where("id = ?", id).Updates(data).Error
}

// Lấy danh sách phiên gửi xe theo ngày và userID
func (r *Repository) FindByDate(
	date time.Time,
	userID uint64,
	page int,
	pageSize int,
) ([]ParkingSession, int64, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	start := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		date.Location(),
	)

	end := start.Add(24 * time.Hour)

	base := r.db.Model(&ParkingSession{}).
		Joins("JOIN rfid_cards ON rfid_cards.uid = parking_sessions.card_uid").
		Where(
			"parking_sessions.entry_time >= ? AND parking_sessions.entry_time < ? AND rfid_cards.user_id = ?",
			start,
			end,
			userID,
		)

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	sessions := make([]ParkingSession, 0)

	err := base.
		Order("exit_time DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&sessions).Error

	if err != nil {
		return nil, 0, err
	}

	return sessions, total, nil
}
