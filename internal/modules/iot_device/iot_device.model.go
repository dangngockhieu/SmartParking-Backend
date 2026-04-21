package iot_device

import "time"

type DeviceStatus string

const (
	DeviceStatusActive   DeviceStatus = "ACTIVE"
	DeviceStatusInactive DeviceStatus = "INACTIVE"
	DeviceStatusError    DeviceStatus = "ERROR"
)

type IoTDevice struct {
	MacAddress string       `gorm:"primaryKey;type:varchar(50)"`
	DeviceName *string      `gorm:"type:varchar(50)"`
	Status     DeviceStatus `gorm:"type:enum('ACTIVE','INACTIVE','ERROR');default:'ACTIVE';not null"`
	LotID      *uint64      `gorm:"index"`
	LastSeen   *time.Time
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (IoTDevice) TableName() string {
	return "iot_devices"
}
