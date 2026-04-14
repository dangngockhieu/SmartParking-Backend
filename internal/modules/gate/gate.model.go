package gate

type GateType string

const (
	GateTypeEntry GateType = "ENTRY"
	GateTypeExit  GateType = "EXIT"
)

type Gate struct {
	ID         uint64   `gorm:"primaryKey;autoIncrement"`
	Name       string   `gorm:"type:varchar(50);not null"`
	Type       GateType `gorm:"type:enum('ENTRY','EXIT');not null"`
	MacAddress string   `gorm:"type:varchar(50);not null;uniqueIndex"`
	LotID      uint64   `gorm:"not null;index"`
	IsActive   bool     `gorm:"not null;default:true"`
}

func (Gate) TableName() string {
	return "gates"
}
