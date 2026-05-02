package wallet

import (
	"backend/configs"

	"gorm.io/gorm"
)

type Module struct {
	Repository *Repository
	Service    *Service
	Handler    *Handler
}

func NewModule(db *gorm.DB, cfg *configs.Config) *Module {
	repo := NewRepository(db)
	payosClient := NewPayOSClient(cfg)
	service := NewService(repo, payosClient, cfg.PayOSReturnURL, cfg.PaymentUpdateStatusCancelURL)
	handler := NewHandler(service, cfg.PayOSCancelURL)

	return &Module{
		Repository: repo,
		Service:    service,
		Handler:    handler,
	}
}
