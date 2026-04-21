package rfid_card

import "time"

type CardType string

const (
	CardTypeRegistered CardType = "REGISTERED"
	CardTypeGuest      CardType = "GUEST"
)

type RfidCard struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	UID       string    `gorm:"type:varchar(20);not null;uniqueIndex" json:"uid"`
	CardType  CardType  `gorm:"type:enum('REGISTERED','GUEST');not null;default:'REGISTERED'" json:"card_type"`
	UserID    *uint64   `gorm:"column:user_id" json:"user_id,omitempty"`
	IsActive  bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (RfidCard) TableName() string {
	return "rfid_cards"
}
