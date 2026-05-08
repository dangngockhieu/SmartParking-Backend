package dashboard

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

// Cấu trúc để nhận kết quả đếm theo giờ
type HourlyCountRow struct {
	Hour  int   `gorm:"column:hour"`
	Count int64 `gorm:"column:count"`
}

// Đếm số xe vào trong ngày
func (r *Repository) CountTodayIn(start time.Time, end time.Time, lotID *uint64) (int64, error) {
	var total int64

	db := r.db.Table("parking_sessions").
		Where("entry_time >= ? AND entry_time < ?", start, end)

	if lotID != nil {
		db = db.Where("lot_id = ?", *lotID)
	}

	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Đếm số xe ra trong ngày
func (r *Repository) CountTodayOut(start time.Time, end time.Time, lotID *uint64) (int64, error) {
	var total int64

	db := r.db.Table("parking_sessions").
		Where("exit_time IS NOT NULL").
		Where("exit_time >= ? AND exit_time < ?", start, end)

	if lotID != nil {
		db = db.Where("lot_id = ?", *lotID)
	}

	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Đếm số xe hiện tại đang đỗ
func (r *Repository) CountCurrentVehicles(lotID *uint64) (int64, error) {
	var total int64

	db := r.db.Table("parking_slots").
		Where("status = ?", "OCCUPIED")

	if lotID != nil {
		db = db.Where("lot_id = ?", *lotID)
	}

	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Đếm tổng số chỗ trong bãi xe
func (r *Repository) CountCapacity(lotID *uint64) (int64, error) {
	var total int64

	db := r.db.Table("parking_slots")

	if lotID != nil {
		db = db.Where("lot_id = ?", *lotID)
	}

	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Đếm số chỗ trống trong bãi xe
func (r *Repository) CountAvailableSlots(lotID *uint64) (int64, error) {
	var total int64

	db := r.db.Table("parking_slots").
		Where("status = ?", "AVAILABLE")

	if lotID != nil {
		db = db.Where("lot_id = ?", *lotID)
	}

	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Lấy số xe vào theo giờ trong ngày
func (r *Repository) GetHourlyIn(start time.Time, end time.Time, lotID *uint64) ([]HourlyCountRow, error) {
	var rows []HourlyCountRow

	db := r.db.Table("parking_sessions").
		Select("HOUR(entry_time) AS hour, COUNT(*) AS count").
		Where("entry_time >= ? AND entry_time < ?", start, end)

	if lotID != nil {
		db = db.Where("lot_id = ?", *lotID)
	}

	if err := db.
		Group("HOUR(entry_time)").
		Order("hour ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

// Lấy số xe ra theo giờ trong ngày
func (r *Repository) GetHourlyOut(start time.Time, end time.Time, lotID *uint64) ([]HourlyCountRow, error) {
	var rows []HourlyCountRow

	db := r.db.Table("parking_sessions").
		Select("HOUR(exit_time) AS hour, COUNT(*) AS count").
		Where("exit_time IS NOT NULL").
		Where("exit_time >= ? AND exit_time < ?", start, end)

	if lotID != nil {
		db = db.Where("lot_id = ?", *lotID)
	}

	if err := db.
		Group("HOUR(exit_time)").
		Order("hour ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return rows, nil
}

// Lấy tên bãi xe
func (r *Repository) GetLotName(lotID uint64) (string, error) {
	var name string

	err := r.db.Table("parking_lots").
		Select("name").
		Where("id = ?", lotID).
		Scan(&name).Error

	if err != nil {
		return "", err
	}

	return name, nil
}

// Tổng doanh thu
func (r *Repository) GetTotalRevenue(lotID *uint64, start time.Time, end time.Time) (int64, error) {
	var total int64
	db := r.db.Table("parking_sessions").
		Select("COALESCE(SUM(fee), 0)").
		Where("exit_time IS NOT NULL").
		Where("exit_time >= ? AND exit_time < ?", start, end)

	if lotID != nil {
		db = db.Where("lot_id = ?", *lotID)
	}
	if err := db.Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
