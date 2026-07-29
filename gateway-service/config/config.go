package config

import (
	"common"
	"gateway/env"
	"os"
)

// LoadConfig đọc cấu hình từ file app.yml
func LoadConfig() {
	// loadBackendLocalEnvFile()

	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")

	common.LoadYMLConfigFile(env.ENV_RUNTIME)
}
