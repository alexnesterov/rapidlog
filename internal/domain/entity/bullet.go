package entity

import (
	"errors"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

var ErrBulletContentRequired = errors.New("bullet content is required")
var ErrBulletContentTooLong = errors.New("bullet content is too long")
var ErrBulletTypeInvalid = errors.New("bullet type must be task, event or note")
var ErrBulletTodayTask = errors.New("bullet is already scheduled for today")
var ErrBulletUserIDRequired = errors.New("bullet user id is required")

var ErrBulletAlreadyCompleted = errors.New("bullet already completed")
var ErrBulletMustBeTaskToBeCompleted = errors.New("bullet must be task to be completed")
var ErrBulletMustBeOpenToBeCompleted = errors.New("bullet must be open to be completed")

var ErrBulletAlreadyCancelled = errors.New("bullet already cancelled")
var ErrBulletMustBeOpenToBeCancelled = errors.New("bullet must be open to be cancelled")

var ErrBulletMustBeOpenTaskToBeMigrated = errors.New("bullet must be open task to be migrated")

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
	UserID    uuid.UUID  `json:"user_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewBullet(userID uuid.UUID, bulletType BulletType, content string) (*Bullet, error) {
	now := time.Now()

	bullet := &Bullet{
		ID:        uuid.New(),
		Type:      bulletType,
		Signifier: SignifierOpen,
		Content:   content,
		UserID:    userID,
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
		return &ValidationError{Err: ErrBulletTypeInvalid}
	}

	if b.Content == "" {
		return &ValidationError{Err: ErrBulletContentRequired}
	}

	if utf8.RuneCountInString(b.Content) > 200 {
		return &ValidationError{Err: ErrBulletContentTooLong}
	}

	if b.UserID == uuid.Nil {
		return &ValidationError{Err: ErrBulletUserIDRequired}
	}

	return nil
}

func (b *Bullet) Complete() error {
	if b.Signifier == SignifierCompleted {
		return &ValidationError{Err: ErrBulletAlreadyCompleted}
	}
	if b.Signifier != SignifierOpen {
		return &ValidationError{Err: ErrBulletMustBeOpenToBeCompleted}
	}

	if b.Type != BulletTask {
		return &ValidationError{Err: ErrBulletMustBeTaskToBeCompleted}
	}

	b.Signifier = SignifierCompleted
	b.UpdatedAt = time.Now()

	return nil
}

func (b *Bullet) Cancel() error {
	if b.Signifier == SignifierCancelled {
		return &ValidationError{Err: ErrBulletAlreadyCancelled}
	}
	if b.Signifier != SignifierOpen {
		return &ValidationError{Err: ErrBulletMustBeOpenToBeCancelled}
	}

	b.Signifier = SignifierCancelled
	b.UpdatedAt = time.Now()

	return nil
}

func (b *Bullet) Migrate() (*Bullet, error) {
	if b.Type != BulletTask || b.Signifier != SignifierOpen {
		return nil, &ValidationError{Err: ErrBulletMustBeOpenTaskToBeMigrated}
	}

	if b.CreatedAt.Format("2006-01-02") == time.Now().Format("2006-01-02") {
		return nil, &ValidationError{Err: ErrBulletTodayTask}
	}

	now := time.Now()
	b.Signifier = SignifierMigrated
	b.UpdatedAt = now

	migrated := &Bullet{
		ID:        uuid.New(),
		Type:      b.Type,
		Signifier: SignifierOpen,
		Content:   b.Content,
		UserID:    b.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return migrated, nil
}
