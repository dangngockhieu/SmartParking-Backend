package parking_slot

import "time"

type SlotStatus string

const (
	SlotStatusAvailable SlotStatus = "AVAILABLE"
	SlotStatusOccupied  SlotStatus = "OCCUPIED"
	SlotStatusMaintain  SlotStatus = "MAINTAIN"
)

type ParkingSlot struct {
	ID         uint64     `gorm:"primaryKey;autoIncrement"`
	Name       string     `gorm:"type:varchar(10);not null;uniqueIndex:uk_lot_name"`
	LotID      uint64     `gorm:"not null;index;uniqueIndex:uk_lot_name;index:idx_lot_status"`
	DeviceMac  string     `gorm:"type:varchar(50);not null;index;uniqueIndex:uk_device_port"`
	PortNumber int        `gorm:"not null;uniqueIndex:uk_device_port"`
	Status     SlotStatus `gorm:"type:enum('AVAILABLE','OCCUPIED','MAINTAIN');default:'AVAILABLE';not null;index:idx_lot_status"`
	CreatedAt  time.Time  `gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `gorm:"autoUpdateTime"`
}

func (ParkingSlot) TableName() string {
	return "parking_slots"
}
