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
			name: "valid title",
			bullet: Bullet{
				Title: "Заголовок",
			},
			wantErr: nil,
		},
		{
			name: "empty title",
			bullet: Bullet{
				Title: "",
			},
			wantErr: ErrTitleRequired,
		},
		{
			name: "title too long",
			bullet: Bullet{
				Title: strings.Repeat("а", 201),
			},
			wantErr: ErrTitleTooLong,
		},
		{
			name: "200 cyrillic chars is valid",
			bullet: Bullet{
				Title: strings.Repeat("ф", 200),
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
				return
			}

			assert.NoError(t, err)
		})
	}
}
