package config

import (
	"bdspro/env"
	"bytes"
	_ "embed"
	"log"
	"os"

	"github.com/spf13/viper"
)

const C_SESSION_ID = "C_SESSION_ID"

// UserProperties chứa cấu hình của ứng dụng
type UserProperties struct {
	Jwt struct {
		SecretKey           string   `mapstructure:"key-generate"`
		TokenPrefix         string   `mapstructure:"tokenPrefix"`
		TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
		AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
		RefreshExpMinutes   uint64   `mapstructure:"refreshExpAfterMinutes"`
		ListPermit          []string `mapstructure:"listPermit"`
		AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	} `mapstructure:"jwt"`
}

var AppProperties UserProperties

//go:embed app.yml
var EmbeddedConfig []byte

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	// viper.SetConfigName("config")   // Tên file không có phần mở rộng
	viper.SetConfigType("yaml") // Loại file config
	// viper.AddConfigPath("./config") // Đọc từ thư mục gốc

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

	if err := viper.Unmarshal(&AppProperties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
		// return nil, err
	}
}
