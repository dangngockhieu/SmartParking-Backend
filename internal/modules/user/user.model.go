package user

import "time"

type Role string

const (
	RoleUser    Role = "USER"
	RoleManager Role = "MANAGER"
	RoleAdmin   Role = "ADMIN"
)

type User struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement"`
	FirstName  string    `gorm:"type:varchar(100);not null"`
	LastName   string    `gorm:"type:varchar(100);not null"`
	Email      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	Password   string    `gorm:"type:varchar(255);not null"`
	Role       Role      `gorm:"type:enum('USER','MANAGER','ADMIN');default:'USER';not null"`
	IsVerified bool      `gorm:"default:false;not null"`
	Money      int64     `gorm:"column:money;not null;default:0" json:"money"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
