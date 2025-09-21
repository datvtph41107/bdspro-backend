package config

import (
	_ "embed"
	"log"
	"user/config/property"

	"github.com/spf13/viper"
)

//go:embed app.yml
var ConfigFile []byte

// LoadConfig đọc cấu hình từ file `.env` hoặc `config.yaml`
func LoadProperties() {
	log.Printf("Đọc config")

	// var config UserProperties
	if err := viper.Unmarshal(&property.AppProperties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
		// return nil, err
	}

	// return &config, nil
}
