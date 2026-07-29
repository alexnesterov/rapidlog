package httpapi

import (
	"slices"
	"sort"

	"github.com/alexnesterov/rapidlog-api/internal/domain/entity"
)

type bulletDayGroup struct {
	Day     string           `json:"day"`
	Bullets []*entity.Bullet `json:"bullets"`
}

func groupBulletsByDay(bullets []*entity.Bullet) []bulletDayGroup {
	sort.Slice(bullets, func(i, j int) bool {
		return bullets[i].CreatedAt.Before(bullets[j].CreatedAt)
	})

	groups := []bulletDayGroup{}

	for _, b := range bullets {
		day := b.CreatedAt.Format("2006-01-02")

		if len(groups) > 0 && groups[len(groups)-1].Day == day {
			groups[len(groups)-1].Bullets = append(groups[len(groups)-1].Bullets, b)
			continue
		}

		g := bulletDayGroup{
			Day:     day,
			Bullets: []*entity.Bullet{b},
		}

		groups = append(groups, g)
	}

	slices.Reverse(groups)

	return groups
}
