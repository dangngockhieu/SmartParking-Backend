package token

import "time"

type RefreshToken struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement"`
	TokenHash string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Device    *string   `gorm:"type:varchar(100)"`
	IP        *string   `gorm:"type:varchar(45)"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UserID    uint64    `gorm:"not null;index"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
