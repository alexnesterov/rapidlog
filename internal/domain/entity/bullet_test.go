package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestBullet_ValidateMigrate(t *testing.T) {

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
			name:       "event cancelled",
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
			name:       "note cancelled",
			bulletType: BulletNote,
			signifier:  SignifierCancelled,
			wantErr:    ErrNotOpenTask,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bullet := Bullet{
				Type:      tc.bulletType,
				Signifier: tc.signifier,
			}

			err := bullet.ValidateMigrate()
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
