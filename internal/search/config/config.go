package config

import (
	"common/configloader"
	"fmt"

	"github.com/spf13/viper"
)

const C_SESSION_ID = "C_SESSION_ID"

// AppProperties chứa cấu hình của ứng dụng
type GlobalProperties struct {
	Jwt struct {
		SecretKey           string   `mapstructure:"key-generate"`
		TokenPrefix         string   `mapstructure:"tokenPrefix"`
		TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
		AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
		RefreshExpMinutes   uint64   `mapstructure:"refreshExpAfterMinutes"`
		ListPermit          []string `mapstructure:"listPermit"`
		AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	} `mapstructure:"jwt"`
	Elastic struct {
		Search struct {
			Url  string `mapstructure:"url"`
			User string `mapstructure:"user"`
			Pass string `mapstructure:"password"`
		} `mapstructure:"search"`
	} `mapstructure:"elastic"`
}

var AppProperties GlobalProperties

type ConfigApp struct {
}

func NewConfigApp() *ConfigApp {
	return &ConfigApp{}
}

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() error {
	if _, err := configloader.LoadRuntimeYML(); err != nil {
		return err
	}

	if err := viper.Unmarshal(&AppProperties); err != nil {
		return fmt.Errorf("unmarshal search config: %w", err)
	}

	return nil
}
