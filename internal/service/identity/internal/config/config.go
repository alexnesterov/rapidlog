package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type DBConfig struct {
	DSN string `mapstructure:"dsn"`
}

type GRPCConfig struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	DB   DBConfig
	GRPC GRPCConfig
}

func Load() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.SetEnvPrefix("identity")

	v.SetDefault("db.dsn", "postgres://rapidlog:rapidlog@localhost:5432/identity")
	v.SetDefault("grpc.port", "50051")

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("config: unmarshal: %w", err)
	}

	if cfg.DB.DSN == "" {
		return nil, fmt.Errorf("config: empty database dsn")
	}
	if cfg.GRPC.Port == "" {
		return nil, fmt.Errorf("config: empty grpc port")
	}

	return &cfg, nil
}
