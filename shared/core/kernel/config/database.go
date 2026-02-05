package config

import (
	"time"
)

type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver"`
	DSN             string        `mapstructure:"dsn"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	LogLevel        string        `mapstructure:"log_level"`
}

func DefaultConfig() DatabaseConfig {
	return DatabaseConfig{
		Driver:          "postgres",
		MaxOpenConns:    50,
		MaxIdleConns:    25,
		ConnMaxLifetime: 5 * time.Minute,
		LogLevel:        "warn",
	}
}
