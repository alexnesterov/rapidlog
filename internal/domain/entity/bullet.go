package entity

import "github.com/google/uuid"

type Status string

const (
	Open Status = "OPEN"
	Done Status = "DONE"
)

type Bullet struct {
	ID     uuid.UUID
	Title  string
	Status Status
}
