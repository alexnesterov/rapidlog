package cancelbullet

import (
	"context"
	"errors"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

type Input struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type Store interface {
	Get(ctx context.Context, id, userID uuid.UUID) (*entity.Bullet, error)
	UpdateSignifier(ctx context.Context, bullet *entity.Bullet) error
}

type Usecase struct {
	store Store
}

func New(store Store) *Usecase {
	return &Usecase{store: store}
}

func (uc *Usecase) Execute(ctx context.Context, in Input) (*entity.Bullet, error) {
	bullet, err := uc.store.Get(ctx, in.ID, in.UserID)
	if err != nil {
		return nil, err
	}

	err = bullet.Cancel()
	if errors.Is(err, entity.ErrBulletAlreadyCancelled) {
		return bullet, nil
	}
	if err != nil {
		return nil, err
	}

	if err := uc.store.UpdateSignifier(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}
