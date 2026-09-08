package configloader

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// RequiredString resolves one process configuration value after the service
// config loader has established file/env precedence. Composition providers use
// this boundary instead of reading Viper directly.
func RequiredString(key string) (string, error) {
	value := strings.TrimSpace(viper.GetString(key))
	if value == "" {
		return "", fmt.Errorf("config %s is required", key)
	}
	return value, nil
}
