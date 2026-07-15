package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestLoad_Defaults(t *testing.T) {
	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "rapidlog", cfg.App.Name)
	assert.Equal(t, "8080", cfg.HTTP.Port)
	assert.Equal(t, 5*time.Second, cfg.HTTP.ReadTimeout)
	assert.Equal(t, 10*time.Second, cfg.HTTP.WriteTimeout)
	assert.Equal(t, 15*time.Second, cfg.HTTP.IdleTimeout)
}

func TestLoad_FromFile(t *testing.T) {
	dir := t.TempDir()
	err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(`
app:
  name: "custom app"

http:
  port: 9090
  read_timeout: 1s
  write_timeout: 2s
  idle_timeout: 3s
`), 0o644)
	require.NoError(t, err)

	t.Chdir(dir)

	cfg, err := Load()
	require.NoError(t, err)

	assert.Equal(t, "custom app", cfg.App.Name)
	assert.Equal(t, "9090", cfg.HTTP.Port)
	assert.Equal(t, 1*time.Second, cfg.HTTP.ReadTimeout)
	assert.Equal(t, 2*time.Second, cfg.HTTP.WriteTimeout)
	assert.Equal(t, 3*time.Second, cfg.HTTP.IdleTimeout)
}
