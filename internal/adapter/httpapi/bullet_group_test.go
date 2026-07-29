package httpapi

import (
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/stretchr/testify/assert"
)

func TestGroupBulletsByDate(t *testing.T) {
	dayA := time.Date(2026, 07, 26, 9, 0, 0, 0, time.UTC)
	dayB := time.Date(2026, 07, 27, 9, 0, 0, 0, time.UTC)
	dayC := time.Date(2026, 07, 28, 9, 0, 0, 0, time.UTC)

	bullets := []*entity.Bullet{
		{Title: "Заголовок 4", CreatedAt: dayC.Add(1 * time.Hour)},
		{Title: "Заголовок 1", CreatedAt: dayA.Add(1 * time.Hour)},
		{Title: "Заголовок 3", CreatedAt: dayB.Add(6 * time.Hour)},
		{Title: "Заголовок 2", CreatedAt: dayB.Add(1 * time.Hour)},
	}

	grouped := groupBulletsByDay(bullets)

	assert.Equal(t, []bulletGroup{
		{Day: dayC.Format("2006-01-02"), Bullets: []*entity.Bullet{{Title: "Заголовок 4", CreatedAt: dayC.Add(1 * time.Hour)}}},
		{Day: dayB.Format("2006-01-02"), Bullets: []*entity.Bullet{{Title: "Заголовок 2", CreatedAt: dayB.Add(1 * time.Hour)}, {Title: "Заголовок 3", CreatedAt: dayB.Add(6 * time.Hour)}}},
		{Day: dayA.Format("2006-01-02"), Bullets: []*entity.Bullet{{Title: "Заголовок 1", CreatedAt: dayA.Add(1 * time.Hour)}}},
	}, grouped)
}
