package entity

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBullet_Validate(t *testing.T) {
	cases := []struct {
		name    string
		bullet  Bullet
		wantErr error
	}{
		{
			name: "type invalid",
			bullet: Bullet{
				Content: "Заголовок",
				Type:    "invalid",
			},
			wantErr: ErrTypeInvalid,
		},
		{
			name: "type empty",
			bullet: Bullet{
				Content: "Заголовок",
				Type:    "",
			},
			wantErr: ErrTypeInvalid,
		},
		{
			name: "valid content",
			bullet: Bullet{
				Content: "Заголовок",
				Type:    BulletTask,
			},
			wantErr: nil,
		},
		{
			name: "empty content",
			bullet: Bullet{
				Content: "",
				Type:    BulletTask,
			},
			wantErr: ErrContentRequired,
		},
		{
			name: "content too long",
			bullet: Bullet{
				Content: strings.Repeat("а", 201),
				Type:    BulletTask,
			},
			wantErr: ErrContentTooLong,
		},
		{
			name: "200 cyrillic chars is valid",
			bullet: Bullet{
				Content: strings.Repeat("ф", 200),
				Type:    BulletTask,
			},
			wantErr: nil,
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
		wantErr    error
	}{
		{
			name:       "task open",
			bulletType: BulletTask,
			signifier:  SignifierOpen,
			wantErr:    nil,
		},
		{
			name:       "task completed",
			bulletType: BulletTask,
			signifier:  SignifierCompleted,
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "task migrated",
			bulletType: BulletTask,
			signifier:  SignifierMigrated,
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "task scheduled",
			bulletType: BulletTask,
			signifier:  SignifierScheduled,
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "task cancelled",
			bulletType: BulletTask,
			signifier:  SignifierCancelled,
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "event open",
			bulletType: BulletEvent,
			signifier:  SignifierOpen,
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "event canceled",
			bulletType: BulletEvent,
			signifier:  SignifierCancelled,
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "note open",
			bulletType: BulletNote,
			signifier:  SignifierOpen,
			wantErr:    ErrNotOpenTask,
		},
		{
			name:       "note canceled",
			bulletType: BulletNote,
			signifier:  SignifierCancelled,
			wantErr:    ErrNotOpenTask,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bullet := Bullet{
				ID:        uuid.New(),
				Type:      tc.bulletType,
				Signifier: tc.signifier,
				Content:   "Test task",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			initialUpdatedAt := bullet.UpdatedAt

			got, err := bullet.Migrate()
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)

				var validationErr *ValidationError
				assert.ErrorAs(t, err, &validationErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, SignifierMigrated, bullet.Signifier)
			assert.True(t, bullet.UpdatedAt.After(initialUpdatedAt))

			assert.NotEqual(t, bullet.ID, got.ID)
			assert.Equal(t, bullet.Type, got.Type)
			assert.Equal(t, SignifierOpen, got.Signifier)
			assert.Equal(t, bullet.Content, got.Content)
			assert.Equal(t, bullet.UpdatedAt, got.CreatedAt)
			assert.Equal(t, bullet.UpdatedAt, got.UpdatedAt)
		})
	}
}
