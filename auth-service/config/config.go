package config

import (
	"auth/env"
	"common/configloader"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// LoadProperties reads Auth configuration and returns startup failures to the process owner.
func LoadProperties() error {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	if err := configloader.LoadYMLFile(env.ENV_RUNTIME); err != nil {
		return err
	}

	if err := viper.Unmarshal(&Properties); err != nil {
		return fmt.Errorf("parse auth properties: %w", err)
	}

	return nil
}
