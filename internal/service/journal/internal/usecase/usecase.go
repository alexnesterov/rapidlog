package usecase

import (
	"github.com/alexnesterov/rapidlog-api/internal/pkg/transactor"
	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/port"
)

type usecase struct {
	repo       port.BulletRepository
	transactor transactor.Transactor
}

func New(r port.BulletRepository, t transactor.Transactor) *usecase {
	return &usecase{
		repo:       r,
		transactor: t,
	}
}

var _ port.Usecase = &usecase{}
