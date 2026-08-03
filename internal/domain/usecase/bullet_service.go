package usecase

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
)

type bulletService struct {
	repo port.BulletRepository
}

func NewBulletService(r port.BulletRepository) *bulletService {
	return &bulletService{
		repo: r,
	}
}

var _ port.BulletService = &bulletService{}
