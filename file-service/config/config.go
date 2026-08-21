package config

import (
	"common/configloader"
	"file/env"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// LoadConfig loads the environment-specific YAML and binds legacy File gRPC
// environment names at the composition/config boundary. Business/runtime
// helpers consume typed RuntimeConfig and never read these env vars directly.
func LoadConfig() error {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	if err := configloader.LoadYMLFile(env.ENV_RUNTIME); err != nil {
		return err
	}

	bindings := map[string]string{
		"internal_grpc.enabled": "QHPRO_FILE_GRPC_ENABLED",
		"internal_grpc.port":    "QHPRO_FILE_GRPC_PORT",
	}
	for key, environment := range bindings {
		if err := viper.BindEnv(key, environment); err != nil {
			return fmt.Errorf("bind %s to %s: %w", key, environment, err)
		}
	}
	return nil
}
