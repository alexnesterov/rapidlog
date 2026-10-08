package httpapi

import (
	"slices"
	"sort"
	"time"

	"github.com/google/uuid"
)

type Bullet struct {
	ID        uuid.UUID `json:"id"`
	Type      string    `json:"type"`
	Signifier string    `json:"signifier"`
	Content   string    `json:"content"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type bulletDayGroup struct {
	Day     string    `json:"day"`
	Bullets []*Bullet `json:"bullets"`
}

func groupBulletsByDay(bullets []*Bullet) []bulletDayGroup {
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
