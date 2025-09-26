package config

import (
	"common"

	_ "embed"
	"log"

	"github.com/spf13/viper"
)

//go:embed app.yml
var EmbeddedConfig []byte

func LoadConfig() {
	common.LoadConfig(EmbeddedConfig)
	if err := viper.Unmarshal(&AppProperties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
	}
}
