package usecase

import (
	"context"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/service/gate/internal/domain/entity"
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

func (s *bulletService) ListBullets(ctx context.Context, userID uuid.UUID) ([]*entity.Bullet, error) {
	bullets, err := s.repo.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(bullets) > 0 {
		return bullets, nil
	}

	seeded := make([]*entity.Bullet, 0, len(demoBullets))
	for _, d := range demoBullets {
		bullet, err := entity.NewBullet(userID, d.Type, d.Content)
		if err != nil {
			return nil, err
		}

		bullet.CreatedAt = bullet.CreatedAt.Add(-24 * time.Hour)
		bullet.UpdatedAt = bullet.CreatedAt

		seeded = append(seeded, bullet)
	}

	err = s.txMgr.WithTransaction(ctx, func(ctx context.Context) error {
		for _, b := range seeded {
			err := s.repo.Create(ctx, b)
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return seeded, nil
}
