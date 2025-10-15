package config

import (
	"common"
	"file/env"
	"os"
)

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	common.LoadYMLConfigFile(env.ENV_RUNTIME)
}
