package config

import (
	"common/configloader"
	"fmt"
	"os"
	"social/env"

	"github.com/spf13/viper"
)

const C_SESSION_ID = "C_SESSION_ID"

// UserProperties chứa cấu hình của ứng dụng
type UserProperties struct {
	Jwt struct {
		SecretKey           string   `mapstructure:"key-generate"`
		TokenPrefix         string   `mapstructure:"tokenPrefix"`
		TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
		AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
		RefreshExpMinutes   uint64   `mapstructure:"refreshExpAfterMinutes"`
		ListPermit          []string `mapstructure:"listPermit"`
		AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	} `mapstructure:"jwt"`
}

var AppProperties UserProperties

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() error {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	if err := configloader.LoadYMLFile(env.ENV_RUNTIME); err != nil {
		return err
	}

	if err := viper.Unmarshal(&AppProperties); err != nil {
		return fmt.Errorf("unmarshal social config: %w", err)
	}

	return nil
}
