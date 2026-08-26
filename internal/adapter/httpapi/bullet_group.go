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
		if bullets[i].CreatedAt.Equal(bullets[j].CreatedAt) {
			return bullets[i].ID.String() < bullets[j].ID.String()
		}
		return bullets[i].CreatedAt.Before(bullets[j].CreatedAt)
	})

	var groups []bulletDayGroup
	var current *bulletDayGroup

	for _, b := range bullets {
		day := b.CreatedAt.Format("2006-01-02")

		if current == nil || current.Day != day {
			groups = append(groups, bulletDayGroup{Day: day})
			current = &groups[len(groups)-1]
		}

		current.Bullets = append(current.Bullets, b)
	}

	slices.Reverse(groups)

	if groups == nil {
		groups = []bulletDayGroup{}
	}

	return groups
}
