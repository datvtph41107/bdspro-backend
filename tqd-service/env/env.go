package env

import (
	"os"
	"strconv"
)

// ENV_RUNTIME represents the runtime environment
var ENV_RUNTIME string

func init() {
	ENV_RUNTIME = os.Getenv("ENV")
	if ENV_RUNTIME == "" {
		ENV_RUNTIME = "development"
	}
}

// GetEnv gets environment variable with default value
func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetEnvInt gets environment variable as integer with default value
func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// GetEnvBool gets environment variable as boolean with default value
func GetEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
