package config

import (
	"auth/env"
	"bytes"
	"common"
	"log"
	"os"

	_ "embed"

	"github.com/spf13/viper"
)

//go:embed app.yml
var EmbeddedConfig []byte

func LoadConfig() {
	common.LoadConfig(EmbeddedConfig)
	if env.ENV_RUNTIME == "release" {
		config_release, err := os.ReadFile("config/product.yml")
		if err != nil {
			log.Fatal(err)
		}
		err = viper.ReadConfig(bytes.NewReader(config_release))
		if err != nil {
			log.Fatalf("Lỗi khi đọc file config: %v", err)
		}
	}
	if err := viper.Unmarshal(&Properties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
	}
}
