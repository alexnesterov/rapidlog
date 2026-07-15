package entity

import (
	"time"

	"github.com/google/uuid"
)

type Collection struct {
	ID        uuid.UUID
	Topic     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
