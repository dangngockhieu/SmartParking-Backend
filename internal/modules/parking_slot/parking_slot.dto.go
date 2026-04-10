package parking_slot

type AdminUpdateParkingSlotRequest struct {
	Status SlotStatus `json:"status" binding:"required"`
}

type SensorUpdateParkingSlotRequest struct {
	Mac        string `json:"mac" binding:"required"`
	Port       int    `json:"port" binding:"required"`
	IsOccupied *bool  `json:"is_occupied" binding:"required"`
}

type CreateParkingSlotRequest struct {
	Name       string `json:"name" binding:"required"`
	LotID      uint64 `json:"lot_id" binding:"required"`
	DeviceMac  string `json:"device_mac" binding:"required"`
	PortNumber int    `json:"port_number" binding:"required"`
}

type ChangeSlotDeviceRequest struct {
	DeviceMac  string `json:"device_mac" binding:"required"`
	PortNumber int    `json:"port_number" binding:"required"`
}

type UpdateParkingSlotResponse struct {
	Changed   bool       `json:"changed"`
	ID        uint64     `json:"id"`
	LotID     uint64     `json:"lot_id"`
	Name      string     `json:"name"`
	Message   string     `json:"message"`
	OldStatus SlotStatus `json:"old_status"`
	NewStatus SlotStatus `json:"new_status"`
}
