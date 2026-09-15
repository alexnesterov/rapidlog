// Package config contains application configuration
package config

import "time"

type AppConfig struct {
	Name string `mapstructure:"name"`
}

type HTTPConfig struct {
	Port         string        `mapstructure:"port"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
}

type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}

type SessionConfig struct {
	CookieName   string        `mapstructure:"cookie_name"`
	CookieTTL    time.Duration `mapstructure:"cookie_ttl"`
	CookieSecure bool          `mapstructure:"cookie_secure"`
}

type Config struct {
	App     AppConfig
	HTTP    HTTPConfig
	DB      DBConfig
	Session SessionConfig
}
