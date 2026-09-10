package usecase

import "github.com/alexnesterov/rapidlog-api/service/identity/internal/domain/port"

type identityService struct {
	userRepo port.UserRepository
}

func NewIdentityService(userRepo port.UserRepository) *identityService {
	return &identityService{
		userRepo: userRepo,
	}
}

var _ port.IdentityService = &identityService{}
