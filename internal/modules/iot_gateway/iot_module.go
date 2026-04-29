package iot_gateway

import (
	"backend/internal/modules/gate"
	"backend/internal/modules/parking_session"
	"backend/internal/modules/parking_slot"
	"backend/internal/modules/rfid_card"
	"backend/internal/modules/user"
)

type Module struct {
	PlateCache *PlateCache
	Service    *Service
	Handler    *Handler
}

func NewModule(
	gateService *gate.Service,
	rfidService *rfid_card.Service,
	sessionService *parking_session.Service,
	parkingSlotService *parking_slot.Service,
	userService *user.Service,
) *Module {
	plateCache := NewPlateCache()
	service := NewService(plateCache, gateService, rfidService, sessionService, parkingSlotService, userService)
	handler := NewHandler(service)

	return &Module{
		PlateCache: plateCache,
		Service:    service,
		Handler:    handler,
	}
}
