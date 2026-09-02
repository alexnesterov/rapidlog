package entity

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBullet(t *testing.T) {
	userID := uuid.New()

	cases := []struct {
		name       string
		userID     uuid.UUID
		bulletType BulletType
		content    string
		wantErr    error
	}{
		{
			name:       "valid task",
			userID:     userID,
			bulletType: BulletTask,
			content:    "Заголовок",
			wantErr:    nil,
		},
		{
			name:       "valid event",
			userID:     userID,
			bulletType: BulletEvent,
			content:    "Заголовок",
			wantErr:    nil,
		},
		{
			name:       "valid note",
			userID:     userID,
			bulletType: BulletNote,
			content:    "Заголовок",
			wantErr:    nil,
		},
		{
			name:       "invalid type",
			userID:     userID,
			bulletType: BulletType("invalid"),
			content:    "Заголовок",
			wantErr:    ErrTypeInvalid,
		},
		{
			name:       "required content",
			userID:     userID,
			bulletType: BulletTask,
			content:    "",
			wantErr:    ErrContentRequired,
		},
		{
			name:       "content too long",
			userID:     userID,
			bulletType: BulletTask,
			content:    strings.Repeat("a", 201),
			wantErr:    ErrContentTooLong,
		},
		{
			name:       "user id required",
			userID:     uuid.Nil,
			bulletType: BulletTask,
			content:    "Заголовок",
			wantErr:    ErrUserIDRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewBullet(tc.userID, tc.bulletType, tc.content)

			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *ValidationError
				assert.ErrorAs(t, err, &validationErr)
				assert.ErrorIs(t, validationErr.Err, tc.wantErr)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)

			assert.Equal(t, tc.userID, got.UserID)
			assert.Equal(t, tc.bulletType, got.Type)
			assert.Equal(t, tc.content, got.Content)
			assert.Equal(t, SignifierOpen, got.Signifier)
			assert.NotEqual(t, uuid.Nil, got.ID)
			assert.WithinDuration(t, time.Now(), got.CreatedAt, 1*time.Second)
			assert.WithinDuration(t, time.Now(), got.UpdatedAt, 1*time.Second)
		})
	}
}

func TestBullet_Validate(t *testing.T) {
	cases := []struct {
		name    string
		bullet  Bullet
		wantErr error
	}{
		{
			name: "type invalid",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: "Заголовок",
				Type:    "invalid",
			},
			wantErr: ErrTypeInvalid,
		},
		{
			name: "type empty",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: "Заголовок",
				Type:    "",
			},
			wantErr: ErrTypeInvalid,
		},
		{
			name: "valid content",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: "Заголовок",
				Type:    BulletTask,
			},
			wantErr: nil,
		},
		{
			name: "empty content",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: "",
				Type:    BulletTask,
			},
			wantErr: ErrContentRequired,
		},
		{
			name: "content too long",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: strings.Repeat("а", 201),
				Type:    BulletTask,
			},
			wantErr: ErrContentTooLong,
		},
		{
			name: "200 cyrillic chars is valid",
			bullet: Bullet{
				UserID:  uuid.New(),
				Content: strings.Repeat("ф", 200),
				Type:    BulletTask,
			},
			wantErr: nil,
		},
		{
			name: "user id required",
			bullet: Bullet{
				UserID:  uuid.Nil,
				Content: "Заголовок",
				Type:    BulletTask,
			},
			wantErr: ErrUserIDRequired,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bullet := tc.bullet
			err := bullet.Validate()

			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)

				var validationErr *ValidationError
				assert.ErrorAs(t, err, &validationErr)

				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestBullet_Migrate(t *testing.T) {
	cases := []struct {
		name       string
		bulletType BulletType
		signifier  Signifier
		createdAt  time.Time
		wantErr    error
	}{
		{
			name:       "task open",
			bulletType: BulletTask,
			signifier:  SignifierOpen,
			createdAt:  time.Now().Add(-24 * time.Hour),
			wantErr:    nil,
		},
		{
			name:       "task completed",
			bulletType: BulletTask,
			signifier:  SignifierCompleted,
			createdAt:  time.Now().Add(-24 * time.Hour),
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "task migrated",
			bulletType: BulletTask,
			signifier:  SignifierMigrated,
			createdAt:  time.Now().Add(-24 * time.Hour),
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "task scheduled",
			bulletType: BulletTask,
			signifier:  SignifierScheduled,
			createdAt:  time.Now().Add(-24 * time.Hour),
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "task cancelled",
			bulletType: BulletTask,
			signifier:  SignifierCancelled,
			createdAt:  time.Now().Add(-24 * time.Hour),
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "event open",
			bulletType: BulletEvent,
			signifier:  SignifierOpen,
			createdAt:  time.Now().Add(-24 * time.Hour),
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "event canceled",
			bulletType: BulletEvent,
			signifier:  SignifierCancelled,
			createdAt:  time.Now().Add(-24 * time.Hour),
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "note open",
			bulletType: BulletNote,
			signifier:  SignifierOpen,
			createdAt:  time.Now().Add(-24 * time.Hour),
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "note canceled",
			bulletType: BulletNote,
			signifier:  SignifierCancelled,
			createdAt:  time.Now().Add(-24 * time.Hour),
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "created at today",
			bulletType: BulletTask,
			signifier:  SignifierOpen,
			createdAt:  time.Now(),
			wantErr:    ErrTodayTask,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bullet := Bullet{
				ID:        uuid.New(),
				Type:      tc.bulletType,
				Signifier: tc.signifier,
				Content:   "Test task",
				UserID:    uuid.New(),
				CreatedAt: tc.createdAt,
				UpdatedAt: time.Now(),
			}

			initialUpdatedAt := bullet.UpdatedAt

			got, err := bullet.Migrate()
			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *ValidationError
				assert.ErrorAs(t, err, &validationErr)
				assert.ErrorIs(t, validationErr.Err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, SignifierMigrated, bullet.Signifier)
			assert.True(t, bullet.UpdatedAt.After(initialUpdatedAt))

			assert.Equal(t, bullet.UserID, got.UserID)
			assert.NotEqual(t, bullet.ID, got.ID)
			assert.Equal(t, bullet.Type, got.Type)
			assert.Equal(t, SignifierOpen, got.Signifier)
			assert.Equal(t, bullet.Content, got.Content)
			assert.Equal(t, bullet.UpdatedAt, got.CreatedAt)
			assert.Equal(t, bullet.UpdatedAt, got.UpdatedAt)
		})
	}
}

func TestBullet_Cancel(t *testing.T) {
	cases := []struct {
		name    string
		bullet  Bullet
		wantErr error
	}{
		{
			name:    "task open",
			bullet:  Bullet{Type: BulletTask, Signifier: SignifierOpen},
			wantErr: nil,
		},
		{
			name:    "event open",
			bullet:  Bullet{Type: BulletEvent, Signifier: SignifierOpen},
			wantErr: nil,
		},
		{
			name:    "note open",
			bullet:  Bullet{Type: BulletNote, Signifier: SignifierOpen},
			wantErr: nil,
		},
		{
			name:    "task completed",
			bullet:  Bullet{Type: BulletTask, Signifier: SignifierCompleted},
			wantErr: ErrNotOpenBullet,
		},
		{
			name:    "task migrated",
			bullet:  Bullet{Type: BulletTask, Signifier: SignifierMigrated},
			wantErr: ErrNotOpenBullet,
		},
		{
			name:    "task canceled",
			bullet:  Bullet{Type: BulletTask, Signifier: SignifierCancelled},
			wantErr: ErrNotOpenBullet,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			initialUpdatedAt := tc.bullet.UpdatedAt

			err := tc.bullet.Cancel()
			if tc.wantErr != nil {
				require.Error(t, err)
				var validationErr *ValidationError
				assert.ErrorAs(t, err, &validationErr)
				assert.ErrorIs(t, err, tc.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, SignifierCancelled, tc.bullet.Signifier)
			assert.True(t, tc.bullet.UpdatedAt.After(initialUpdatedAt))
		})
	}
}
