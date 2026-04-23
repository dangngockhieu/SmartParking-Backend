package iot_device

import (
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

type IoTDeviceWithLot struct {
	MacAddress string       `gorm:"column:mac_address"`
	DeviceName *string      `gorm:"column:device_name"`
	Status     DeviceStatus `gorm:"column:status"`
	LotID      *uint64      `gorm:"column:lot_id"`
	LotName    *string      `gorm:"column:lot_name"`
	LastSeen   *time.Time   `gorm:"column:last_seen"`
	CreatedAt  time.Time    `gorm:"column:created_at"`
	UpdatedAt  time.Time    `gorm:"column:updated_at"`
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByMacAddress(macAddress string) (*IoTDevice, error) {
	var device IoTDevice
	err := r.db.Where("mac_address = ?", macAddress).First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

func (r *Repository) CreateDevice(device *IoTDevice) error {
	return r.db.Create(device).Error
}

func (r *Repository) FindAll(filters GetIoTDevicesQuery) ([]IoTDeviceWithLot, error) {
	var devices []IoTDeviceWithLot
	query := r.db.
		Table("iot_devices AS d").
		Select(`
			d.mac_address,
			d.device_name,
			d.status,
			d.lot_id,
			l.name AS lot_name,
			d.last_seen,
			d.created_at,
			d.updated_at
		`).
		Joins("LEFT JOIN parking_lots AS l ON l.id = d.lot_id")

	if filters.LotID != nil {
		query = query.Where("d.lot_id = ?", *filters.LotID)
	}

	if filters.Status != nil {
		query = query.Where("d.status = ?", *filters.Status)
	}

	if filters.Keyword != "" {
		prefixPattern := filters.Keyword + "%"
		query = query.Where("(d.mac_address LIKE ? OR d.device_name LIKE ?)", prefixPattern, prefixPattern)
	}

	if err := query.Order("d.created_at desc").Scan(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *Repository) UpdateByMacAddress(macAddress string, updates map[string]interface{}) error {
	return r.db.Model(&IoTDevice{}).Where("mac_address = ?", macAddress).Updates(updates).Error
}

func (r *Repository) FindByMacAddressWithLot(macAddress string) (*IoTDeviceWithLot, error) {
	var device IoTDeviceWithLot
	err := r.db.
		Table("iot_devices AS d").
		Select(`
			d.mac_address,
			d.device_name,
			d.status,
			d.lot_id,
			l.name AS lot_name,
			d.last_seen,
			d.created_at,
			d.updated_at
		`).
		Joins("LEFT JOIN parking_lots AS l ON l.id = d.lot_id").
		Where("d.mac_address = ?", macAddress).
		First(&device).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}
