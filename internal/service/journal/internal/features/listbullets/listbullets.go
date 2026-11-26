package listbullets

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

type Input struct {
	UserID uuid.UUID
}

type Store interface {
	List(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error)
}

type Usecase struct {
	store Store
}

func New(store Store) *Usecase {
	return &Usecase{store: store}
}

func (uc *Usecase) Execute(ctx context.Context, in Input) ([]*entity.Bullet, error) {
	bullets, err := uc.store.List(ctx, in.UserID)
	if err != nil {
		return nil, err
	}

	return bullets, nil
}
