package entity

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusOpen Status = "OPEN"
	StatusDone Status = "DONE"
)

type Bullet struct {
	ID           uuid.UUID
	CollectionID uuid.UUID
	Title        string
	Date         time.Time
	Status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
