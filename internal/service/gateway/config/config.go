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

type KratosConfig struct {
	// Публичный API Kratos: проверка сессии (whoami) и прокси /auth/* для браузера
	PublicURL string `mapstructure:"public_url"`
}

type Config struct {
	Port    string `mapstructure:"port"`
	HTTP    HTTPConfig
	Journal string `mapstructure:"journal"`
	Kratos  KratosConfig
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("gateway")

	v.SetDefault("port", 8080)
	v.SetDefault("journal", "localhost:50052")

	v.SetDefault("http.read_timeout", 5*time.Second)
	v.SetDefault("http.write_timeout", 10*time.Second)
	v.SetDefault("http.idle_timeout", 15*time.Second)

	v.SetDefault("kratos.public_url", "http://localhost:4433")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if cfg.Port == "" {
		return nil, fmt.Errorf("config: empty port")
	}
	if cfg.Journal == "" {
		return nil, fmt.Errorf("config: empty journal")
	}
	if cfg.Kratos.PublicURL == "" {
		return nil, fmt.Errorf("config: empty kratos public url")
	}

	return &cfg, nil
}
