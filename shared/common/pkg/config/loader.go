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

	if err := v.MergeInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read %s config: %w", env, err)
		}
	}

	// 4. local.yaml (gitignore)
	v.SetConfigName("local")
	if err := v.MergeInConfig(); err != nil {
		// optional
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
