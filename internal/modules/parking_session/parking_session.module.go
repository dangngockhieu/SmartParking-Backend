package parking_session

import "gorm.io/gorm"

type Module struct {
	Repository *Repository
	Service    *Service
	Handler    *Handler
}

func NewModule(db *gorm.DB) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		Repository: repo,
		Service:    service,
		Handler:    handler,
	}
}
