package entity

import (
	"time"

	"github.com/google/uuid"
)

var demoBullets = []struct {
	Kind    Kind
	Content string
}{
	{KindTask, "Смигрируй эту задачу на сегодня"},
	{KindEvent, "Встреча в 15:00"},
	{KindNote, "Идея для заметок — сюда пишешь мысли без действия"},
}

func NewDemoBullets(userID uuid.UUID) ([]*Bullet, error) {
	bullets := make([]*Bullet, 0, len(demoBullets))
	for _, d := range demoBullets {
		bullet, err := NewBullet(userID, d.Kind, d.Content)
		if err != nil {
			return nil, err
		}

		bullet.CreatedAt = bullet.CreatedAt.Add(-24 * time.Hour)
		bullet.UpdatedAt = bullet.CreatedAt

		bullets = append(bullets, bullet)
	}

	return bullets, nil
}
