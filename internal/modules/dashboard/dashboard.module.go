package dashboard

import "gorm.io/gorm"

type Module struct {
	Repo    *Repository
	Service *Service
	Handler *Handler
}

func NewModule(db *gorm.DB) *Module {
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	return &Module{
		Repo:    repo,
		Service: service,
		Handler: handler,
	}
}
