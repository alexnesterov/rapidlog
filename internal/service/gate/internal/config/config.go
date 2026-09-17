// Package config contains application configuration
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type HTTPConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}

type IdentityConfig struct {
	Addr string `mapstructure:"addr"`
}

type SessionConfig struct {
	CookieName   string        `mapstructure:"cookie_name"`
	CookieTTL    time.Duration `mapstructure:"cookie_ttl"`
	CookieSecure bool          `mapstructure:"cookie_secure"`
}

type Config struct {
	HTTP     HTTPConfig
	DB       DBConfig
	Identity IdentityConfig
	Session  SessionConfig
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("rapidlog")

	v.SetDefault("db.dsn", "postgres://rapidlog:rapidlog@localhost:5432/rapidlog")
	v.SetDefault("http.port", 8080)
	v.SetDefault("identity.addr", "localhost:50051")

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

	if cfg.DB.DSN == "" {
		return nil, fmt.Errorf("config: empty database dsn")
	}
	if cfg.HTTP.Port == "" {
		return nil, fmt.Errorf("config: empty http port")
	}
	if cfg.Identity.Addr == "" {
		return nil, fmt.Errorf("config: empty identity addr")
	}

	return &cfg, nil
}
