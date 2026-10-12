package usecase

import (
	"github.com/alexnesterov/rapidlog-api/internal/pkg/transactor"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port"
)

type bulletService struct {
	repo       port.BulletRepository
	transactor transactor.Transactor
}

func NewBulletService(r port.BulletRepository, t transactor.Transactor) *bulletService {
	return &bulletService{
		repo:       r,
		transactor: t,
	}
}

var _ port.BulletService = &bulletService{}
