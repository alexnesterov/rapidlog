package usecase

import (
	"github.com/alexnesterov/rapidlog-api/internal/pkg/transactor"
)

type Usecase struct {
	repo       BulletRepository
	transactor transactor.Transactor
}

func New(r BulletRepository, t transactor.Transactor) *Usecase {
	return &Usecase{
		repo:       r,
		transactor: t,
	}
}
