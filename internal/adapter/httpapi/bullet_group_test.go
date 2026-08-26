package httpapi

import (
	"testing"
	"time"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGroupBulletsByDay(t *testing.T) {
	dayA := time.Date(2026, 07, 26, 9, 0, 0, 0, time.UTC)
	dayB := time.Date(2026, 07, 27, 9, 0, 0, 0, time.UTC)
	dayC := time.Date(2026, 07, 28, 9, 0, 0, 0, time.UTC)

	bullets := []*entity.Bullet{
		{Content: "Заголовок 4", CreatedAt: dayC.Add(1 * time.Hour)},
		{Content: "Заголовок 1", CreatedAt: dayA.Add(1 * time.Hour)},
		{Content: "Заголовок 3", CreatedAt: dayB.Add(6 * time.Hour)},
		{Content: "Заголовок 2", CreatedAt: dayB.Add(1 * time.Hour)},
	}

	grouped := groupBulletsByDay(bullets)

	assert.Equal(t, []bulletDayGroup{
		{Day: dayC.Format("2006-01-02"), Bullets: []*entity.Bullet{{Content: "Заголовок 4", CreatedAt: dayC.Add(1 * time.Hour)}}},
		{Day: dayB.Format("2006-01-02"), Bullets: []*entity.Bullet{{Content: "Заголовок 2", CreatedAt: dayB.Add(1 * time.Hour)}, {Content: "Заголовок 3", CreatedAt: dayB.Add(6 * time.Hour)}}},
		{Day: dayA.Format("2006-01-02"), Bullets: []*entity.Bullet{{Content: "Заголовок 1", CreatedAt: dayA.Add(1 * time.Hour)}}},
	}, grouped)
}

func TestGroupBulletsByDay_EqualCreatedAtTieBreaksByID(t *testing.T) {
	day := time.Date(2026, 8, 25, 9, 0, 0, 0, time.UTC)
	idLow := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	idHigh := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	cases := []struct {
		name    string
		bullets []*entity.Bullet
	}{
		{
			name: "high id first in input",
			bullets: []*entity.Bullet{
				{ID: idHigh, Content: "Второй по ID", CreatedAt: day},
				{ID: idLow, Content: "Первый по ID", CreatedAt: day},
			},
		},
		{
			name: "low id first in input",
			bullets: []*entity.Bullet{
				{ID: idLow, Content: "Первый по ID", CreatedAt: day},
				{ID: idHigh, Content: "Второй по ID", CreatedAt: day},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			grouped := groupBulletsByDay(tc.bullets)

			assert.Equal(t, []bulletDayGroup{
				{
					Day: day.Format("2006-01-02"),
					Bullets: []*entity.Bullet{
						{ID: idLow, Content: "Первый по ID", CreatedAt: day},
						{ID: idHigh, Content: "Второй по ID", CreatedAt: day},
					},
				},
			}, grouped)
		})
	}
}

func TestGroupBulletsByDay_EmptyAndNilInput(t *testing.T) {
	cases := []struct {
		name    string
		bullets []*entity.Bullet
	}{
		{name: "empty input", bullets: []*entity.Bullet{}},
		{name: "nil input", bullets: nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			grouped := groupBulletsByDay(tc.bullets)

			assert.NotNil(t, grouped)
			assert.Equal(t, []bulletDayGroup{}, grouped)
		})
	}
}
