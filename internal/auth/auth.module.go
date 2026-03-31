package auth

import (
	"backend/internal/auth/mail"
	"backend/internal/auth/token"
	"backend/pkg/database"

	"gorm.io/gorm"
)

type Module struct {
	Repository *Repository
	Service    *Service
	Handler    *Handler
}

func NewModule(
	db *gorm.DB,
	redis *database.RedisClient,
	tokenService *token.Service,
	mailService *mail.Service,
) *Module {
	repo := NewRepository(db)
	service := NewService(repo, tokenService, mailService, redis)
	handler := NewHandler(service, mailService)

	return &Module{
		Repository: repo,
		Service:    service,
		Handler:    handler,
	}
}
