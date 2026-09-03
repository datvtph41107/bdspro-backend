package config

import (
	"common/configloader"
	"fmt"
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
	_, err := configloader.LoadRuntimeYML()
	if err != nil {
		return err
	}
	// YAML giữ topology/default; các giá trị runtime nhạy cảm hoặc khác nhau
	// giữa host và container phải override được tại process boundary. Bind rõ
	// alias để Unmarshal và các compatibility provider cùng thấy một nguồn thật.
	for key, names := range map[string][]string{
		"database.dsn":             {"HUB_DATABASE_URL", "DATABASE_DSN"},
		"redis.host":               {"REDIS_HOST", "HUB_REDIS_ADDRESS"},
		"redis.db":                 {"REDIS_DB", "HUB_REDIS_DB"},
		"jwt.key-generate":         {"JWT_KEY_GENERATE"},
		"security.api_key":         {"HUB_API_KEY"},
		"rpc.user.address":         {"HUB_USER_GRPC_ADDRESS"},
		"rpc.bdspro.address":       {"HUB_BDSPRO_GRPC_ADDRESS"},
		"rpc.notification.address": {"HUB_NOTIFICATION_GRPC_ADDRESS"},
	} {
		args := append([]string{key}, names...)
		if err := viper.BindEnv(args...); err != nil {
			return fmt.Errorf("bind %s env: %w", key, err)
		}
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
