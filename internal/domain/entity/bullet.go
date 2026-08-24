package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrContentRequired = errors.New("content is required")
var ErrContentTooLong = errors.New("content is too long")
var ErrTypeInvalid = errors.New("type must be task, event or note")

var ErrNotOpenTask = errors.New("bullet is not an open task")
var ErrTodayTask = errors.New("bullet is already scheduled for today")

type BulletType string

const (
	BulletTask  BulletType = "task"
	BulletEvent BulletType = "event"
	BulletNote  BulletType = "note"
)

type Signifier string

const (
	SignifierOpen      Signifier = "open"
	SignifierCompleted Signifier = "completed"
	SignifierMigrated  Signifier = "migrated"
	SignifierScheduled Signifier = "scheduled"
	SignifierCancelled Signifier = "cancelled"
)

type Bullet struct {
	ID        uuid.UUID  `json:"id"`
	Type      BulletType `json:"type"`
	Signifier Signifier  `json:"signifier"`
	Content   string     `json:"content"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewBullet(bulletType BulletType, content string) (*Bullet, error) {
	now := time.Now()

	bullet := &Bullet{
		ID:        uuid.New(),
		Type:      bulletType,
		Signifier: SignifierOpen,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := bullet.Validate(); err != nil {
		return nil, err
	}

	return bullet, nil
}

func (b *Bullet) Validate() error {
	switch b.Type {
	case BulletTask, BulletEvent, BulletNote:
	default:
		return &ValidationError{Err: ErrTypeInvalid}
	}

	if b.Content == "" {
		return &ValidationError{Err: ErrContentRequired}
	}

	if utf8.RuneCountInString(b.Content) > 200 {
		return &ValidationError{Err: ErrContentTooLong}
	}

	return nil
}

func (b *Bullet) Migrate() (*Bullet, error) {
	if b.Type != BulletTask || b.Signifier != SignifierOpen {
		return nil, &ValidationError{Err: ErrNotOpenTask}
	}

	if b.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return nil, &ValidationError{Err: ErrTodayTask}
	}

	now := time.Now()
	b.Signifier = SignifierMigrated
	b.UpdatedAt = now

	migrated := &Bullet{
		ID:        uuid.New(),
		Type:      b.Type,
		Signifier: SignifierOpen,
		Content:   b.Content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return migrated, nil
}
