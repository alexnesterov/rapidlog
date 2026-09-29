package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DSN  string `mapstructure:"dsn"`
	Port string `mapstructure:"port"`
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("identity")

	v.SetDefault("dsn", "postgres://rapidlog:rapidlog@localhost:5432/identity")
	v.SetDefault("port", "50051")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if cfg.DSN == "" {
		return nil, fmt.Errorf("config: empty database dsn")
	}
	if cfg.Port == "" {
		return nil, fmt.Errorf("config: empty grpc port")
	}

	return &cfg, nil
}
