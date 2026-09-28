package config

import (
	"bytes"
	_config "common/config"
	"common/configloader"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// LoadProperties nạp config theo lớp:
//  1. config/runtime.yml (base: CORS, ports, defaults)
//  2. config/<ENV>.yml (môi trường: DSN, secrets) — merge đè lên base
//  3. QHPRO_CONFIG_FILE / CONFIG_FILE (nếu có) — merge cuối cùng
//
// DSN/topology DB lấy từ YAML, không bị .env đè qua AutomaticEnv cho database.dsn.
func LoadProperties() error {
	selection, err := configloader.ResolveRuntimeSelection()
	if err != nil {
		return err
	}

	configureViper()

	if err := readConfigFile("config/runtime.yml", false); err != nil {
		return fmt.Errorf("load runtime base config: %w", err)
	}

	overridePath := firstNonEmpty(
		strings.TrimSpace(os.Getenv("QHPRO_CONFIG_FILE")),
		strings.TrimSpace(os.Getenv("CONFIG_FILE")),
	)
	if overridePath != "" {
		if err := readConfigFile(overridePath, true); err != nil {
			return fmt.Errorf("merge operator config %s: %w", overridePath, err)
		}
	} else if selection.ConfigProfile != "" && selection.ConfigProfile != "runtime" {
		envPath := "config/" + selection.ConfigProfile + ".yml"
		if err := readConfigFile(envPath, true); err != nil {
			return fmt.Errorf("merge environment config %s: %w", envPath, err)
		}
	}

	for _, binding := range []struct {
		key  string
		envs []string
	}{
		{key: "redis.host", envs: []string{"USER_REDIS_ADDRESS", "REDIS_ADDRESS", "REDIS_HOST"}},
		{key: "redis.pass", envs: []string{"USER_REDIS_PASSWORD", "REDIS_PASSWORD"}},
		{key: "redis.db", envs: []string{"USER_REDIS_DB", "REDIS_DB"}},
		{key: "jwt.key-generate", envs: []string{"JWT_KEY_GENERATE"}},
		{key: "outbound_messaging.enabled", envs: []string{"QHPRO_USER_OUTBOUND_MESSAGING_ENABLED"}},
		{key: "rpc.user.address", envs: []string{"RPC_USER_ADDRESS"}},
		{key: "rpc.auth.address", envs: []string{"RPC_AUTH_ADDRESS"}},
		{key: "rpc.organization.address", envs: []string{"RPC_ORGANIZATION_ADDRESS"}},
		{key: "rpc.hub.address", envs: []string{"RPC_HUB_ADDRESS"}},
		{key: "rpc.bdspro.address", envs: []string{"RPC_BDSPRO_ADDRESS"}},
		{key: "rpc.notification.address", envs: []string{"RPC_NOTIFICATION_ADDRESS"}},
		{key: "rpc.payment.address", envs: []string{"RPC_PAYMENT_ADDRESS"}},
		{key: "rpc.chat.address", envs: []string{"RPC_CHAT_ADDRESS"}},
		{key: "rpc.crm.address", envs: []string{"RPC_CRM_ADDRESS"}},
		{key: "rpc.tqd.address", envs: []string{"RPC_TQD_ADDRESS"}},
	} {
		args := append([]string{binding.key}, binding.envs...)
		if err := viper.BindEnv(args...); err != nil {
			return fmt.Errorf("bind %s env: %w", binding.key, err)
		}
	}
	if err := viper.Unmarshal(&Properties); err != nil {
		return fmt.Errorf("parse user properties: %w", err)
	}
	return nil
}

func configureViper() {
	viper.Reset()
	viper.SetConfigType("yaml")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()
}

func readConfigFile(path string, merge bool) error {
	configFile, err := _config.ReadFile(path)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(configFile)
	if merge {
		if err := viper.MergeConfig(reader); err != nil {
			return fmt.Errorf("merge config %s: %w", path, err)
		}
		return nil
	}
	if err := viper.ReadConfig(reader); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}
	return nil
}
