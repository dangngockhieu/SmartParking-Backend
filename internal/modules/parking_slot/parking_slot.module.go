package parking_slot

import (
	"backend/internal/realtime/parking"

	"gorm.io/gorm"
)

type Module struct {
	Repository *Repository
	Service    *Service
	Handler    *Handler
}

func NewModule(db *gorm.DB, hub *parking.Hub) *Module {
	repo := NewRepository(db)
	service := NewService(repo, hub)
	handler := NewHandler(service)

	return &Module{
		Repository: repo,
		Service:    service,
		Handler:    handler,
	}
}
