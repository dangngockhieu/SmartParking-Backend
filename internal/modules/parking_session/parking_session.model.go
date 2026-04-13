package parking_session

import (
	"time"

	"backend/internal/modules/rfid_card"
)

type ParkingSession struct {
	ID          uint64             `gorm:"primaryKey;autoIncrement" json:"id"`
	LotID       uint64             `gorm:"not null;index" json:"lot_id"`
	SlotID      *uint64            `gorm:"index" json:"slot_id,omitempty"`
	CardUID     string             `gorm:"type:varchar(20);not null;index" json:"card_uid"`
	CardType    rfid_card.CardType `gorm:"type:enum('REGISTERED','GUEST');not null" json:"card_type"`
	PlateNumber string             `gorm:"type:varchar(20);not null;index" json:"plate_number"`
	EntryTime   time.Time          `gorm:"not null;autoCreateTime;index" json:"entry_time"`
	ExitTime    *time.Time         `json:"exit_time,omitempty"`
	Fee         int64              `json:"fee,omitempty"`
	IsActive    bool               `gorm:"not null;default:true;index" json:"is_active"`
}

func (ParkingSession) TableName() string {
	return "parking_sessions"
}
