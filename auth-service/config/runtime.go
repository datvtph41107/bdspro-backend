package config

import (
	"common"
	"log"

	_ "embed"

	"github.com/spf13/viper"
)

//go:embed app.yml
var EmbeddedConfig []byte

func LoadConfig() {
	common.LoadConfig(EmbeddedConfig)
	if err := viper.Unmarshal(&Properties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
	}
}
