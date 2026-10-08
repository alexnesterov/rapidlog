// Package config contains application configuration
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type HTTPConfig struct {
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type SessionConfig struct {
	CookieName   string        `mapstructure:"cookie_name"`
	CookieTTL    time.Duration `mapstructure:"cookie_ttl"`
	CookieSecure bool          `mapstructure:"cookie_secure"`
}

type Config struct {
	DSN      string `mapstructure:"dsn"`
	Port     string `mapstructure:"port"`
	HTTP     HTTPConfig
	Identity string `mapstructure:"identity"`
	Session  SessionConfig
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("gate")

	v.SetDefault("dsn", "postgres://rapidlog:rapidlog@localhost:5432/rapidlog")
	v.SetDefault("port", 8080)
	v.SetDefault("identity", "localhost:50051")

	v.SetDefault("http.read_timeout", 5*time.Second)
	v.SetDefault("http.write_timeout", 10*time.Second)
	v.SetDefault("http.idle_timeout", 15*time.Second)

	v.SetDefault("session.cookie_name", "session_id")
	v.SetDefault("session.cookie_ttl", 365*24*time.Hour)
	v.SetDefault("session.cookie_secure", true)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if cfg.DSN == "" {
		return nil, fmt.Errorf("config: empty database dsn")
	}
	if cfg.Port == "" {
		return nil, fmt.Errorf("config: empty port")
	}
	if cfg.Identity == "" {
		return nil, fmt.Errorf("config: empty identity")
	}

	return &cfg, nil
}
