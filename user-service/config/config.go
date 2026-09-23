package config

import (
	"common/configloader"
	"fmt"

	"github.com/spf13/viper"
)

// LoadProperties đọc cấu hình từ config/<ENV_RUNTIME>.yml (hoặc CONFIG_FILE
// khi operator mount). DSN/topology DB lấy từ file, không bị .env đè.
func LoadProperties() error {
	_, err := configloader.LoadRuntimeYML()
	if err != nil {
		return err
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
