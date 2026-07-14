package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

var ErrEmptyAppName = errors.New("empty app name")

func setDefaults() {
	viper.SetDefault("app.name", "dotline")

	viper.SetDefault("http.port", 8080)
	viper.SetDefault("http.read_timeout", 5*time.Second)
	viper.SetDefault("http.write_timeout", 10*time.Second)
	viper.SetDefault("http.idle_timeout", 15*time.Second)
}

func validate(cfg *Config) error {
	if cfg.App.Name == "" {
		return ErrEmptyAppName
	}
	return nil
}

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("configs")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("DOTLINE")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		var notFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &notFoundError) {
			return nil, fmt.Errorf("config: read file: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("config: validate: %w", err)
	}

	return &cfg, nil
}
