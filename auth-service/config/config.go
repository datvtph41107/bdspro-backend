package config

import (
	"auth/env"
	"common"
	"log"
	"os"

	"github.com/spf13/viper"
)

// LoadProperties đọc cấu hình từ file config.yml
func LoadProperties() {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	common.LoadYMLConfigFile(env.ENV_RUNTIME)

	if err := viper.Unmarshal(&Properties); err != nil {
		log.Fatalf("Lỗi parse config: %v", err)
	}
}
