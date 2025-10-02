package config

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/spf13/viper"
)

//go:embed config.yml
var configFile []byte

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

// LoadConfig loads configuration from config.yml
func LoadConfig() *Properties {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	err := viper.ReadConfig(bytes.NewReader(configFile))
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	AppConfig = &Properties{}
	err = viper.Unmarshal(AppConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return AppConfig
}
