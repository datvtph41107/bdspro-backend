package config

import (
	"common/configloader"
	"fmt"

	"github.com/spf13/viper"
)

// LoadConfig loads the environment-specific File-service configuration.
// Environment/YAML lookup remains confined to the composition boundary.
func LoadConfig() error {
	if _, err := configloader.LoadRuntimeYML(); err != nil {
		return err
	}
	// Root .env dùng tên có namespace để không đụng DATABASE_DSN của service
	// khác. Alias cũ vẫn được giữ sau nó cho compatibility khi deploy.
	for key, names := range map[string][]string{
		"database.dsn": {"FILE_DATABASE_DSN", "DATABASE_DSN"},
	} {
		args := append([]string{key}, names...)
		if err := viper.BindEnv(args...); err != nil {
			return fmt.Errorf("bind %s env: %w", key, err)
		}
	}

	return nil
}
