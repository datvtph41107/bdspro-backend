package config

import (
	"common/configloader"
	"fmt"
	"hub/env"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const C_SESSION_ID = "C_SESSION_ID"

type UserRPCProperties struct {
	Address string `mapstructure:"address"`
}

type RPCProperties struct {
	User UserRPCProperties `mapstructure:"user"`
}

type HubProperties struct {
	Applink struct {
		BaseURL string `mapstructure:"base_url"`
	} `mapstructure:"applink"`
	Jwt struct {
		SecretKey           string   `mapstructure:"key-generate"`
		TokenPrefix         string   `mapstructure:"tokenPrefix"`
		TokenExpirationDays int      `mapstructure:"tokenExpirationAfterDays"`
		AccessExpMinutes    int      `mapstructure:"accessExpAfterMinutes"`
		RefreshExpMinutes   uint64   `mapstructure:"refreshExpAfterMinutes"`
		ListPermit          []string `mapstructure:"listPermit"`
		AuthorizationHeader string   `mapstructure:"authorizationHeader"`
	} `mapstructure:"jwt"`
	Security struct {
		APIKey           string   `mapstructure:"api_key"`
		ProtectedMethods []string `mapstructure:"protected_methods"`
	} `mapstructure:"security"`
	RPC RPCProperties `mapstructure:"rpc"`
}

var AppProperties HubProperties

// LoadConfig đọc cấu hình từ file config.yml
func LoadConfig() error {
	env.ENV_RUNTIME = os.Getenv("ENV_RUNTIME")
	if err := configloader.LoadYMLFile(env.ENV_RUNTIME); err != nil {
		return err
	}

	if err := viper.Unmarshal(&AppProperties); err != nil {
		return fmt.Errorf("unmarshal hub config: %w", err)
	}

	if err := AppProperties.Validate(); err != nil {
		return fmt.Errorf("validate hub config: %w", err)
	}

	return nil
}

func (c HubProperties) Validate() error {
	if strings.TrimSpace(c.RPC.User.Address) == "" {
		return fmt.Errorf("rpc.user.address is required")
	}
	return nil
}
