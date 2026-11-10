package usecase

import (
	"github.com/alexnesterov/rapidlog-api/internal/pkg/transactor"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port"
)

type Usecase struct {
	repo       port.BulletRepository
	transactor transactor.Transactor
}

func New(r port.BulletRepository, t transactor.Transactor) *Usecase {
	return &Usecase{
		repo:       r,
		transactor: t,
	}
}
