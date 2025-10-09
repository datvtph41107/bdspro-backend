package config

import (
	"bytes"
	"common"
	_ "embed"
	"log"
	"os"
	"user/config/property"
	"user/env"

	"github.com/spf13/viper"
)

//go:embed app.yml
var EmbeddedConfig []byte

// LoadConfig đọc cấu hình từ file `.env` hoặc `config.yaml`
func LoadProperties() {
	config_release, err := os.ReadFile("config/product.yml")
	if err != nil {
		log.Fatal(err)
	}
	if env.ENV_RUNTIME == "release" {
		err := viper.ReadConfig(bytes.NewReader(config_release))
		if err != nil {
			log.Fatalf("Lỗi khi đọc file config: %v", err)
		}
	} else {
		common.LoadConfig(EmbeddedConfig)
	}
	if err := viper.Unmarshal(&property.AppProperties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
	}
}
