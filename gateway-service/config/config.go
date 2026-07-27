package config

import (
	"common"
	"gateway/env"
	"os"
	"strings"
)

// LoadConfig đọc cấu hình từ file app.yml
func LoadConfig() {
	loadBackendLocalEnvFile()

	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	if strings.TrimSpace(env.ENV_RUNTIME) == "" {
		env.ENV_RUNTIME = "local"
		_ = os.Setenv("ENV_RUNTIME", env.ENV_RUNTIME)
	}

	common.LoadYMLConfigFile(env.ENV_RUNTIME)
}
