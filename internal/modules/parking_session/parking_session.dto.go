package parking_session

import "time"

type ParkingSessionResponse struct {
	ID          uint64     `json:"id"`
	LotID       uint64     `json:"lot_id"`
	SlotID      *uint64    `json:"slot_id,omitempty"`
	CardUID     string     `json:"card_uid"`
	CardType    string     `json:"card_type"`
	PlateNumber string     `json:"plate_number"`
	EntryTime   time.Time  `json:"entry_time"`
	ExitTime    *time.Time `json:"exit_time,omitempty"`
	Fee         int64      `json:"fee,omitempty"`
	IsActive    bool       `json:"is_active"`
}

type ManageParkingSessionResponse struct {
	ID          uint64     `json:"id"`
	LotID       uint64     `json:"lot_id"`
	SlotID      *uint64    `json:"slot_id,omitempty"`
	CardUID     string     `json:"card_uid"`
	CardType    string     `json:"card_type"`
	PlateNumber string     `json:"plate_number"`
	EntryTime   time.Time  `json:"entry_time"`
	ExitTime    *time.Time `json:"exit_time,omitempty"`
	Fee         int64      `json:"fee,omitempty"`
	IsActive    bool       `json:"is_active"`
	OwnerName   *string    `json:"owner_name"`
}

type ParkingSessionListMeta struct {
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
	CurrentPage   int   `json:"currentPage"`
	PageSize      int   `json:"pageSize"`
}

type ParkingSessionListResponse struct {
	Data []ParkingSessionResponse `json:"data"`
	Meta ParkingSessionListMeta   `json:"meta"`
}

type ManageParkingSessionListResponse struct {
	Data []ManageParkingSessionResponse `json:"data"`
	Meta ParkingSessionListMeta         `json:"meta"`
}

type CreateParkingSessionInput struct {
	LotID       uint64 `json:"lot_id"`
	CardUID     string `json:"card_uid"`
	CardType    string `json:"card_type"`
	PlateNumber string `json:"plate_number"`
}

type AssignSlotInput struct {
	SessionID uint64 `json:"session_id"`
	SlotID    uint64 `json:"slot_id"`
}

type FinishParkingSessionInput struct {
	SessionID uint64 `json:"session_id"`
	Fee       int64  `json:"fee"`
}
