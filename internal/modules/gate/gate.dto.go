package gate

// ─── Request ─────────────────────────────────────────────────────────────────

type CreateGateRequest struct {
	Name       string   `json:"name"        binding:"required,max=50"`
	Type       GateType `json:"type"        binding:"required,oneof=ENTRY EXIT"`
	MacAddress string   `json:"mac_address" binding:"required,max=50"`
	LotID      uint64   `json:"lot_id"      binding:"required"`
}

type UpdateGateRequest struct {
	MacAddress string `json:"mac_address" binding:"omitempty,max=50"`
}

// ─── URI params ──────────────────────────────────────────────────────────────

type GateURIParams struct {
	GateID uint64 `uri:"gateId" binding:"required"`
}

type LotURIParams struct {
	LotID uint64 `uri:"lotId" binding:"required"`
}

// ─── Response ─────────────────────────────────────────────────────────────────

type GateResponse struct {
	ID         uint64   `json:"id"`
	Name       string   `json:"name"`
	Type       GateType `json:"type"`
	MacAddress string   `json:"mac_address"`
	LotID      uint64   `json:"lot_id"`
	IsActive   bool     `json:"is_active"`
}

func ToGateResponse(g *Gate) GateResponse {
	return GateResponse{
		ID:         g.ID,
		Name:       g.Name,
		Type:       g.Type,
		MacAddress: g.MacAddress,
		LotID:      g.LotID,
		IsActive:   g.IsActive,
	}
}
