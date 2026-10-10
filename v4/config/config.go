package config

import (
	"bytes"
	_config "common/config"
	"common/configloader"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	tqdconfig "userv4/modules/tqd/config"
)

// LoadProperties nạp config theo lớp (YAML là nguồn duy nhất cho topology/secrets app):
//  1. config/runtime.yml (base: CORS, ports, defaults)
//  2. config/<ENV>.yml (môi trường: DSN, Redis, JWT, RPC, …) — merge đè lên base
//  3. QHPRO_CONFIG_FILE / CONFIG_FILE (nếu có) — merge cuối cùng
//
// Không BindEnv / AutomaticEnv cho redis, database, jwt, rpc, oauth, …
// Process env chỉ chọn profile (QHPRO_ENVIRONMENT / ENV_RUNTIME) hoặc đường dẫn
// operator overlay — không đè giá trị trong YAML.
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

	if err := viper.Unmarshal(&Properties); err != nil {
		return fmt.Errorf("parse user properties: %w", err)
	}
	tqdconfig.InitLayerResolver(&Properties.TQD.Layer)
	return nil
}

func configureViper() {
	viper.Reset()
	viper.SetConfigType("yaml")
	// No AutomaticEnv / BindEnv: workstation often sources root/.env and
	// legacy user-service/.env which would otherwise clobber local.yml.
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
