package parking_lot

type CreateParkingLotRequest struct {
	Name     string `json:"name" binding:"required"`
	Location string `json:"location" binding:"required"`
}

type UpdateParkingLotRequest struct {
	Name     *string `json:"name"`
	Location *string `json:"location"`
}

type ParkingLotResponse struct {
	ID       uint64  `json:"id"`
	Name     string  `json:"name"`
	Location *string `json:"location,omitempty"`
}

type ParkingLotSlotResponse struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	DeviceMac  string `json:"device_mac"`
	PortNumber int    `json:"port_number"`
}

type ParkingLotGateResponse struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	MacAddress string `json:"mac_address"`
	IsActive   bool   `json:"is_active"`
}

type ParkingLotStatsResponse struct {
	Total     int64 `json:"total"`
	Available int64 `json:"available"`
	Occupied  int64 `json:"occupied"`
	Maintain  int64 `json:"maintain"`
}

type ParkingLotDetailResponse struct {
	ID       uint64                   `json:"id"`
	Name     string                   `json:"name"`
	Location *string                  `json:"location,omitempty"`
	Slots    []ParkingLotSlotResponse `json:"slots"`
	Stats    ParkingLotStatsResponse  `json:"stats"`
}
