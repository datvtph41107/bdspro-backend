package config

import (
	"common/configloader"
	"fmt"
	"os"
	"user/env"

	"github.com/spf13/viper"
)

// LoadProperties đọc cấu hình từ file config.yml
func LoadProperties() error {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	if err := configloader.LoadYMLFile(env.ENV_RUNTIME); err != nil {
		return err
	}
	if err := viper.Unmarshal(&Properties); err != nil {
		return fmt.Errorf("parse user properties: %w", err)
	}
	return nil
}
