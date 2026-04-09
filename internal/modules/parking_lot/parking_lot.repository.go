package parking_lot

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
}

type parkingLotStatsRow struct {
	Status string
	Count  int64
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// Tạo mới bãi đỗ xe
func (r *Repository) Create(lot *ParkingLot) error {
	return r.db.Create(lot).Error
}

// Lấy danh sách tất cả bãi đỗ xe
func (r *Repository) FindAll() ([]ParkingLot, error) {
	var lots []ParkingLot
	err := r.db.
		Select("id", "name", "location").
		Order("id ASC").
		Find(&lots).Error
	if err != nil {
		return nil, err
	}
	return lots, nil
}

// Lấy thông tin chi tiết bãi đỗ theo ID
func (r *Repository) FindByID(id uint64) (*ParkingLot, error) {
	var lot ParkingLot
	err := r.db.
		Select("id", "name", "location").
		First(&lot, id).Error
	if err != nil {
		return nil, err
	}
	return &lot, nil
}

// Cập nhật thông tin bãi đỗ theo ID
func (r *Repository) UpdateByID(id uint64, data map[string]interface{}) error {
	return r.db.Model(&ParkingLot{}).Where("id = ?", id).Updates(data).Error
}

// Xóa bãi đỗ theo ID
func (r *Repository) FindSlotsByLotID(lotID uint64) ([]ParkingLotSlotResponse, error) {
	var slots []ParkingLotSlotResponse

	err := r.db.
		Table("parking_slots").
		Select("id, name, status, device_mac, port_number").
		Where("lot_id = ?", lotID).
		Order("name ASC").
		Scan(&slots).Error

	if err != nil {
		return nil, err
	}

	return slots, nil
}

// Lấy danh sách cổng của bãi đỗ theo ID
func (r *Repository) FindGatesByLotID(lotID uint64) ([]ParkingLotGateResponse, error) {
	var gates []ParkingLotGateResponse

	err := r.db.
		Table("gates").
		Select("id, name, type, mac_address, is_active").
		Where("lot_id = ?", lotID).
		Order("id ASC").
		Scan(&gates).Error
	if err != nil {
		return nil, err
	}

	return gates, nil
}

// Lấy thông tin chi tiết bãi đỗ theo ID, bao gồm danh sách slot và thống kê
func (r *Repository) CountStatsByLotID(lotID uint64) ([]parkingLotStatsRow, error) {
	var rows []parkingLotStatsRow

	err := r.db.
		Table("parking_slots").
		Select("status, COUNT(*) as count").
		Where("lot_id = ?", lotID).
		Group("status").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// Lấy user ID theo email
func (r *Repository) GetUserIDByEmail(email string) (*uint64, error) {
	var userID uint64
	err := r.db.Where("email = ?", email).Select("id").Scan(&userID).Error
	if err != nil {
		return nil, err
	}
	return &userID, nil
}
