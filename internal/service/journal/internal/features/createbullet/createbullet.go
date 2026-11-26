package createbullet

import (
	"context"

	"github.com/alexnesterov/rapidlog-api/internal/service/journal/internal/entity"
	"github.com/google/uuid"
)

type Input struct {
	UserID  uuid.UUID
	Kind    entity.Kind
	Content string
}

type Store interface {
	Insert(ctx context.Context, bullet *entity.Bullet) error
}

type Usecase struct {
	store Store
}

func New(store Store) *Usecase {
	return &Usecase{store: store}
}

func (uc *Usecase) Execute(ctx context.Context, in Input) (*entity.Bullet, error) {
	bullet, err := entity.NewBullet(in.UserID, in.Kind, in.Content)
	if err != nil {
		return nil, err
	}

	if err := uc.store.Insert(ctx, bullet); err != nil {
		return nil, err
	}

	return bullet, nil
}
