package entity

import (
	"time"

	"github.com/google/uuid"
)

var demoBullets = []struct {
	Type    BulletType
	Content string
}{
	{BulletTask, "Смигрируй эту задачу на сегодня"},
	{BulletEvent, "Встреча в 15:00"},
	{BulletNote, "Идея для заметок — сюда пишешь мысли без действия"},
}

func NewDemoBullets(userID uuid.UUID) ([]*Bullet, error) {
	bullets := make([]*Bullet, 0, len(demoBullets))
	for _, d := range demoBullets {
		bullet, err := NewBullet(userID, d.Type, d.Content)
		if err != nil {
			return nil, err
		}

		bullet.CreatedAt = bullet.CreatedAt.Add(-24 * time.Hour)
		bullet.UpdatedAt = bullet.CreatedAt

		bullets = append(bullets, bullet)
	}

	return bullets, nil
}
