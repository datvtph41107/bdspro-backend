package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var validate = validator.New()

// MustLoad loads config for a service, panics on error.
func MustLoad(serviceName string, cfg interface{}) {
	if err := load(serviceName, cfg); err != nil {
		panic(fmt.Sprintf("[%s] CONFIG FATAL: %v", serviceName, err))
	}
}

func load(serviceName string, cfg interface{}) error {
	// 1. Set defaults from struct tags (defaults.Set)
	if err := defaults.Set(cfg); err != nil {
		return fmt.Errorf("set defaults: %w", err)
	}

	v := viper.New()

	// 2. Env prefix
	prefix := strings.ToUpper(serviceName)
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// 3. Read environment-specific config file (dev.yaml, staging.yaml, prod.yaml)
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}
	v.SetConfigName(env)
	v.SetConfigType("yaml")
	v.AddConfigPath(fmt.Sprintf("./%s/config", serviceName))
	v.AddConfigPath(fmt.Sprintf("./%s", serviceName))
	// QHPro: also resolve configuration relative to the running service directory
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read %s config: %w", env, err)
		}
	}

	// 4. local config override
	// QHPro: explicitly merge a service-local local.yml/local.yaml overlay
	localConfigLoaded := false
	for _, localConfigFile := range []string{
		"./config/local.yml",
		"./config/local.yaml",
		"./local.yml",
		"./local.yaml",
	} {
		if _, err := os.Stat(localConfigFile); err != nil {
			continue
		}

		v.SetConfigFile(localConfigFile)
		if err := v.MergeInConfig(); err != nil {
			return fmt.Errorf("merge local config %s: %w", localConfigFile, err)
		}

		localConfigLoaded = true
		break
	}

	if !localConfigLoaded {
		v.SetConfigName("local")
		_ = v.MergeInConfig()
	}

	// 5. Unmarshal
	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	// 6. Schema validation
	if err := validate.Struct(cfg); err != nil {
		return fmt.Errorf("schema validation: %w", err)
	}

	// 7. Business validation (if implemented)
	if v, ok := cfg.(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return fmt.Errorf("business validation: %w", err)
		}
	}

	return nil
}
