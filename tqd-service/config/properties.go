package config

import (
	"bytes"
	_ "embed"
	"log"
	"os"
	"tqd/env"

	"github.com/spf13/viper"
)

//go:embed config.yml
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
func LoadConfig() *Properties {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	viper.SetConfigType("yaml") // Loại file config

	// Đọc biến môi trường (nếu có)
	viper.AutomaticEnv()

	var err error
	if env.ENV_RUNTIME == "release" {
		config_release, err := os.ReadFile("config/product.yml")
		if err != nil {
			log.Fatal(err)
		}
		err = viper.ReadConfig(bytes.NewReader(config_release))
		if err != nil {
			log.Fatalf("Lỗi khi đọc file config: %v", err)
		}
	} else {
		// Đọc file config
		err = viper.ReadConfig(bytes.NewReader(EmbeddedConfig))
	}

	if err != nil {
		log.Fatalf("Lỗi khi đọc file config: %v", err)
	}

	AppConfig = &Properties{}
	err = viper.Unmarshal(AppConfig)
	if err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	return AppConfig
}
