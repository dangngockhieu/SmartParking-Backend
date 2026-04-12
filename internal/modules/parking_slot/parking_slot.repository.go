package parking_slot

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByID(id uint64) (*ParkingSlot, error) {
	var slot ParkingSlot
	err := r.db.First(&slot, id).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *Repository) FindByDeviceMacAndPort(deviceMac string, port int) (*ParkingSlot, error) {
	var slot ParkingSlot
	err := r.db.
		Where("device_mac = ? AND port_number = ?", deviceMac, port).
		First(&slot).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *Repository) Create(slot *ParkingSlot) error {
	return r.db.Create(slot).Error
}

func (r *Repository) UpdateStatus(id uint64, status SlotStatus) error {
	return r.db.Model(&ParkingSlot{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Repository) FindConflictSlot(deviceMac string, portNumber int, excludeID uint64) (*ParkingSlot, error) {
	var slot ParkingSlot
	err := r.db.
		Where("device_mac = ? AND port_number = ? AND id <> ?", deviceMac, portNumber, excludeID).
		First(&slot).Error
	if err != nil {
		return nil, err
	}
	return &slot, nil
}

func (r *Repository) ChangeDevice(
	slotID uint64,
	deviceMac string,
	portNumber int,
) (*ParkingSlot, error) {
	var updated ParkingSlot

	err := r.db.Transaction(func(tx *gorm.DB) error {
		var slot ParkingSlot
		if err := tx.First(&slot, slotID).Error; err != nil {
			return err
		}

		if err := tx.Model(&ParkingSlot{}).
			Where("id = ?", slotID).
			Updates(map[string]any{
				"device_mac":  deviceMac,
				"port_number": portNumber,
			}).Error; err != nil {
			return err
		}

		if err := tx.First(&updated, slotID).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *Repository) IsAvailable(lotID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&ParkingSlot{}).
		Where("lot_id = ? AND status = ?", lotID, "AVAILABLE").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
