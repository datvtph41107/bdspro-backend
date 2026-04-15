// shared/pkg/config/loader.go
package config

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// LoadBase: load config chung cho tất cả service.
// serviceName được dùng làm env prefix: "chat" → "CHAT_SERVER_PORT"
//
// Tại sao cần serviceName?
// Vì 5 service chạy trên cùng một machine (không có K8s) có thể conflict env.
// CHAT_SERVER_PORT=8080 và CRM_SERVER_PORT=8081 — tường biên rõ ràng.
func LoadConfig(serviceName string, cfg interface{}) error {
	v := viper.New()

	// Convention thống nhất: PREFIX_SECTION_KEY
	// chat → CHAT_SERVER_PORT, CHAT_DATABASE_DSN
	// crm  → CRM_SERVER_PORT, CRM_DATABASE_DSN
	prefix := strings.ToUpper(serviceName)
	v.SetEnvPrefix(prefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Config file: {serviceName}/config.yaml
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(fmt.Sprintf("./%s", serviceName))
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("read config: %w", err)
		}
		// Không có file config là OK — env vars là đủ cho production
	}

	if err := v.Unmarshal(cfg); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return validateStruct(cfg)
}

func MustLoad(serviceName string, cfg interface{}) {
	if err := LoadConfig(serviceName, cfg); err != nil {
		panic(fmt.Sprintf("[%s] CONFIG FATAL: %v", serviceName, err))
	}
}

func validateStruct(cfg interface{}) error {
	v := validator.New()
	if err := v.Struct(cfg); err != nil {
		return formatErrors(err)
	}
	return nil
}

// formatErrors: lỗi validation phải đọc được, không phải stack trace
func formatErrors(err error) error {
	var msgs []string
	for _, e := range err.(validator.ValidationErrors) {
		msgs = append(msgs, fmt.Sprintf(
			"  field '%s': failed '%s' (got: '%v')",
			e.Field(), e.Tag(), e.Value(),
		))
	}
	return fmt.Errorf("config validation failed:\n%s", strings.Join(msgs, "\n"))
}
