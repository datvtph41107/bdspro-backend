package config

import (
	_ "embed"
	"log"

	"github.com/spf13/viper"
)

const C_SESSION_ID = "C_SESSION_ID"

//go:embed app.yml
var ConfigFile []byte

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

	Firebase struct {
		File    string `mapstructure:"file"`
		FileURL string `mapstructure:"fileUrl"`
	} `mapstructure:"firebase"`
}

var AppProperties UserProperties

// LoadConfig đọc cấu hình từ file `.env` hoặc `config.yaml`
func LoadProperties() {
	log.Printf("Đọc config")

	// var config UserProperties
	if err := viper.Unmarshal(&AppProperties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
		// return nil, err
	}

	// return &config, nil
}
