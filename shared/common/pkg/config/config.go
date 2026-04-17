package config

import "time"

type BaseConfig struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	Feature  FeatureConfig  `mapstructure:"feature"`
}

type AppConfig struct {
	Name            string        `mapstructure:"name"                     validate:"required"`
	Environment     string        `mapstructure:"environment"             validate:"oneof=development staging production test"`
	Version         string        `mapstructure:"version"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"         default:"15s"`
}

type ServerConfig struct {
	Port            int           `mapstructure:"port"                    validate:"required,min=1024,max=65535" default:"8080" reload:"cold"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"            default:"30s" reload:"cold"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"           default:"30s" reload:"cold"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"        default:"15s" reload:"cold"`
}

type DatabaseConfig struct {
	DSN             string        `mapstructure:"dsn"                     sensitive:"true"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"          default:"25"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"          default:"5"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"       default:"300s"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"  validate:"oneof=debug info warn error" default:"info" reload:"hot"`
	Format string `mapstructure:"format" validate:"oneof=json text" default:"json" reload:"cold"`
}

type FeatureConfig struct {
	EnableNewFlow   bool `mapstructure:"enable_new_flow"   reload:"hot"`
	MaintenanceMode bool `mapstructure:"maintenance_mode"  reload:"hot"`
}
