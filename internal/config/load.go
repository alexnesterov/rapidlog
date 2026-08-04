package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var ErrEmptyAppName = errors.New("empty app name")
var ErrEmptyDSN = errors.New("empty database dsn")

func setDefaults(v *viper.Viper) {
	v.SetDefault("app.name", "rapidlog")

	v.SetDefault("http.port", 8080)
	v.SetDefault("http.read_timeout", 5*time.Second)
	v.SetDefault("http.write_timeout", 10*time.Second)
	v.SetDefault("http.idle_timeout", 15*time.Second)

	v.SetDefault("db.dsn", "postgres://rapidlog:rapidlog@localhost:5432/rapidlog")
}

func validate(cfg *Config) error {
	if cfg.App.Name == "" {
		return ErrEmptyAppName
	}
	if cfg.DB.DSN == "" {
		return ErrEmptyDSN
	}
	return nil
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath("configs")
	v.AddConfigPath(".")

	v.AutomaticEnv()

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("rapidlog")

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		var notFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &notFoundError) {
			return nil, fmt.Errorf("config: read file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("config: validate: %w", err)
	}

	return &cfg, nil
}
