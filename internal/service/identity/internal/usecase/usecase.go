package usecase

import "github.com/alexnesterov/rapidlog-api/internal/service/identity/internal/port"

type Usecase struct {
	repo port.UserRepository
}

func New(r port.UserRepository) *Usecase {
	return &Usecase{
		repo: r,
	}
}
