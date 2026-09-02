package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/alexnesterov/rapidlog-api/internal/domain/port"
	"github.com/google/uuid"
)

var demoBullets = []struct {
	Type    entity.BulletType
	Content string
}{
	{
		Type:    entity.BulletTask,
		Content: "Смигрируй эту задачу на сегодня",
	},
	{
		Type:    entity.BulletEvent,
		Content: "Встреча в 15:00",
	},
	{
		Type:    entity.BulletNote,
		Content: "Идея для заметок — сюда пишешь мысли без действия",
	},
}

func (s *userService) ResolveUser(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err == nil {
		return user.ID, nil
	}
	if !errors.Is(err, port.ErrNotFound) {
		return uuid.Nil, err
	}

	newUser := entity.NewUser()

	err = s.txMgr.WithTransaction(ctx, func(ctx context.Context) error {
		if err := s.userRepo.Create(ctx, newUser); err != nil {
			return err
		}

		for _, b := range demoBullets {
			bullet, err := entity.NewBullet(newUser.ID, b.Type, b.Content)
			if err != nil {
				return err
			}

			bullet.CreatedAt = bullet.CreatedAt.Add(-24 * time.Hour)
			bullet.UpdatedAt = bullet.CreatedAt

			if err := s.bulletRepo.Create(ctx, bullet); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	return newUser.ID, nil
}
