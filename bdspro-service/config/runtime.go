package config

import (
	"common/configloader"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

const C_SESSION_ID = "C_SESSION_ID"

// UserProperties chứa cấu hình của ứng dụng
type UserProperties struct {
	Redis struct {
		Host     string `mapstructure:"host"`
		Password string `mapstructure:"pass"`
	} `mapstructure:"redis"`
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
func LoadConfig() {
	if _, err := configloader.LoadRuntimeYML(); err != nil {
		slog.Error(
			"load BDSPro runtime config failed",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&AppProperties); err != nil {
		slog.Error(
			"parse BDSPro runtime config failed",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
}
