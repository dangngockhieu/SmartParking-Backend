package rfid_card

type CreateRfidCardRequest struct {
	UID      string   `json:"uid" binding:"required"`
	CardType CardType `json:"card_type" binding:"required"`
	UserID   *uint64  `json:"user_id"`
}

type UpdateRfidCardRequest struct {
	CardType *CardType `json:"card_type"`
	UserID   *uint64   `json:"user_id"`
	IsActive *bool     `json:"is_active"`
}

type RfidCardResponse struct {
	ID        uint64   `json:"id"`
	UID       string   `json:"uid"`
	CardType  CardType `json:"card_type"`
	UserID    *uint64  `json:"user_id,omitempty"`
	IsActive  bool     `json:"is_active"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type MyRfidCardResponse struct {
	ID           uint64   `json:"id"`
	CardUID      string   `json:"cardUid"`
	UserID       *uint64  `json:"userId"`
	OwnerName    *string  `json:"ownerName"`
	Status       CardType `json:"status"`
	IsActive     bool     `json:"isActive"`
	RegisteredAt *string  `json:"registeredAt"`
}

type RfidCardStatisticsResponse struct {
	TotalCards        int64 `json:"totalCards"`
	RegisteredCards   int64 `json:"registeredCards"`
	UnregisteredCards int64 `json:"unregisteredCards"`
	ActiveCards       int64 `json:"activeCards"`
}

type RfidCardListItem struct {
	ID           uint64   `json:"id"`
	CardUID      string   `json:"cardUid"`
	UserID       *uint64  `json:"userId"`
	OwnerName    *string  `json:"ownerName"`
	Status       CardType `json:"status"`
	IsActive     bool     `json:"isActive"`
	RegisteredAt *string  `json:"registeredAt"`
}

type RfidCardListMeta struct {
	TotalElements int64 `json:"totalElements"`
	TotalPages    int   `json:"totalPages"`
	CurrentPage   int   `json:"currentPage"`
	PageSize      int   `json:"pageSize"`
}

type RfidCardListResponse struct {
	Data []RfidCardListItem `json:"data"`
	Meta RfidCardListMeta   `json:"meta"`
}
