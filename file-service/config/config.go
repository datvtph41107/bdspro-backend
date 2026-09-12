package config

import (
	"common/configloader"
)

// LoadConfig loads the environment-specific File-service configuration.
// Environment/YAML lookup remains confined to the composition boundary.
//
// Database DSN is owned by the environment YAML:
//   - workstation: config/runtime.yml merged with config/local.yml
//   - deploy: CONFIG_FILE / QHPRO_CONFIG_FILE, or DATABASE_DSN via AutomaticEnv
//
// Do not BindEnv database.dsn to service .env — that fights local.yml the same
// way payment-service stopped owning DSN through PAYMENT_DATABASE_URL.
func LoadConfig() error {
	if _, err := configloader.LoadRuntimeYML(); err != nil {
		return err
	}
	return nil
}
