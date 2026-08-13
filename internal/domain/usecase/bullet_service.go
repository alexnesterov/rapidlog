package usecase

import (
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
)

type bulletService struct {
	repo  port.BulletRepository
	txMgr port.TransactionManager
}

func NewBulletService(r port.BulletRepository, txMgr port.TransactionManager) *bulletService {
	return &bulletService{
		repo:  r,
		txMgr: txMgr,
	}
}

var _ port.BulletService = &bulletService{}
