package rfid_card

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Create thêm một thẻ RFID mới vào cơ sở dữ liệu
func (r *Repository) Create(card *RfidCard) error {
	return r.db.Create(card).Error
}

// FindByID tìm kiếm thẻ RFID theo ID
func (r *Repository) FindByID(id uint64) (*RfidCard, error) {
	var card RfidCard
	err := r.db.First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

type RfidCardDetailRow struct {
	ID        uint64
	UID       string
	CardType  CardType
	UserID    *uint64
	OwnerName *string
	IsActive  bool
	CreatedAt time.Time
}

// Lấy thẻ theo userId
func (r *Repository) GetByUserID(userID uint64) (*RfidCardDetailRow, error) {
	var row RfidCardDetailRow

	err := r.db.Table("rfid_cards AS rc").
		Joins("LEFT JOIN users AS u ON u.id = rc.user_id").
		Where("rc.user_id = ?", userID).
		Select(`
			rc.id,
			rc.uid,
			rc.card_type,
			rc.user_id,
			NULLIF(
				TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))),
				''
			) AS owner_name,
			rc.is_active,
			rc.created_at
		`).
		First(&row).Error

	if err != nil {
		return nil, err
	}

	return &row, nil
}

// Update thẻ RFID theo ID với dữ liệu mới
func (r *Repository) UpdateByID(id uint64, data map[string]any) error {
	return r.db.Model(&RfidCard{}).Where("id = ?", id).Session(&gorm.Session{}).UpdateColumns(data).Error
}

func (r *Repository) FindByUID(uid string) (*RfidCard, error) {
	var card RfidCard
	err := r.db.Where("uid = ?", uid).First(&card).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

type RfidCardListRow struct {
	ID          uint64
	UID         string
	CardType    CardType
	UserID      *uint64
	OwnerName   *string
	IsActive    bool
	CreatedAt   time.Time
	PlateNumber *string
}

func (r *Repository) CountStatistics() (int64, int64, int64, int64, error) {
	var total int64
	if err := r.db.Model(&RfidCard{}).
		Count(&total).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	var registered int64
	if err := r.db.Model(&RfidCard{}).
		Where("card_type = ?", CardTypeRegistered).
		Count(&registered).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	var unregistered int64
	if err := r.db.Model(&RfidCard{}).
		Where("card_type = ?", CardTypeGuest).
		Count(&unregistered).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	var active int64
	if err := r.db.Model(&RfidCard{}).
		Where("is_active = ?", true).
		Count(&active).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	return total, registered, unregistered, active, nil
}

func (r *Repository) FindWithFilters(
	status string,
	keyword string,
	page int,
	pageSize int,
) ([]RfidCardListRow, int64, error) {
	if page < 1 {
		page = 1
	}

	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	latestBase := r.db.Table("parking_sessions AS ps1").
		Select("ps1.card_uid, ps1.plate_number, ps1.entry_time")

	maxEntry := r.db.Table("parking_sessions").
		Select("card_uid, MAX(entry_time) AS max_entry")

	maxEntry = maxEntry.Group("card_uid")

	latest := latestBase.
		Joins(
			"JOIN (?) AS ps2 ON ps1.card_uid = ps2.card_uid AND ps1.entry_time = ps2.max_entry",
			maxEntry,
		)

	base := r.db.Table("rfid_cards AS rc").
		Joins("LEFT JOIN users AS u ON u.id = rc.user_id").
		Joins("LEFT JOIN (?) AS ps ON ps.card_uid = rc.uid", latest)

	if status != "" {
		base = base.Where("rc.card_type = ?", status)
	}

	keyword = strings.TrimSpace(keyword)
	if keyword != "" {
		like := keyword + "%"
		base = base.Where(
			`rc.uid LIKE ?`,
			like,
		)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).
		Distinct("rc.id").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	rows := make([]RfidCardListRow, 0)

	err := base.
		Select(`
			rc.id,
			rc.uid,
			rc.card_type,
			rc.user_id,
			NULLIF(
				TRIM(CONCAT(COALESCE(u.first_name, ''), ' ', COALESCE(u.last_name, ''))),
				''
			) AS owner_name,
			rc.is_active,
			rc.created_at,
			ps.plate_number
		`).
		Order("rc.id DESC").
		Offset(offset).
		Limit(pageSize).
		Scan(&rows).Error

	if err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}

func (r *Repository) GetUserIDByEmail(email string) (*uint64, error) {
	var userID uint64
	if err := r.db.Table("users").
		Select("id").
		Where("email = ?", email).
		Scan(&userID).Error; err != nil {
		return nil, err
	}

	return &userID, nil
}
