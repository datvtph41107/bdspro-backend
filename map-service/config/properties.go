package config

import (
	"common/configloader"
	_ "embed"
	"fmt"

	"github.com/spf13/viper"
)

//go:embed runtime.yml
var EmbeddedConfig []byte

// Properties holds all configuration properties
type Properties struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Eureka   EurekaConfig   `mapstructure:"eureka"`
	Key      KeyConfig      `mapstructure:"key"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port int `mapstructure:"port"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	KeyGenerate         string   `mapstructure:"key-generate"`
	TokenPrefix         string   `mapstructure:"tokenPrefix"`
	TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
	AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
	RefreshExpMinutes   int      `mapstructure:"refreshExpAfterMinutes"`
	AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	ListPermit          []string `mapstructure:"listPermit"`
}

// EurekaConfig holds Eureka configuration
type EurekaConfig struct {
	URL  string `mapstructure:"url"`
	Port int    `mapstructure:"port"`
}

// KeyConfig holds key configuration
type KeyConfig struct {
	JWT string `mapstructure:"jwt"`
}

var AppConfig *Properties

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() (*Properties, error) {
	if _, err := configloader.LoadRuntimeYML(); err != nil {
		return nil, fmt.Errorf("load map runtime config: %w", err)
	}

	AppConfig = &Properties{}
	if err := viper.Unmarshal(AppConfig); err != nil {
		return nil, fmt.Errorf("parse map runtime config: %w", err)
	}

	return AppConfig, nil
}
