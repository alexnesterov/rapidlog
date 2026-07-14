package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	cases := []struct {
		name string
		cfg  *Config
		want error
	}{
		{
			name: "empty app name",
			cfg: &Config{
				App: AppConfig{
					Name: "",
				},
			},
			want: ErrEmptyAppName,
		},
		{
			name: "valid config",
			cfg: &Config{
				App: AppConfig{
					Name: "DOTLINE",
				},
			},
			want: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate(tc.cfg)
			assert.ErrorIs(t, err, tc.want)
		})
	}
}
