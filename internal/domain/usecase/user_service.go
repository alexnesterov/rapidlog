package usecase

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
)

type userService struct {
	userRepo   port.UserRepository
	bulletRepo port.BulletRepository
	txMgr      port.TransactionManager
}

func NewUserService(
	userRepo port.UserRepository,
	bulletRepo port.BulletRepository,
	txMgr port.TransactionManager,
) *userService {
	return &userService{
		userRepo:   userRepo,
		bulletRepo: bulletRepo,
		txMgr:      txMgr,
	}
}

var _ port.UserService = &userService{}
