package env

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Runtime struct {
	Environment string
	ConfigPath  string
}

var ENV_RUNTIME *Runtime

func init() {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}

	configPath := "config"
	if os.Getenv("CONFIG_PATH") != "" {
		configPath = os.Getenv("CONFIG_PATH")
	}

	ENV_RUNTIME = &Runtime{
		Environment: environment,
		ConfigPath:  configPath,
	}

	loadConfig()
}

func loadConfig() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(ENV_RUNTIME.ConfigPath)
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Auto read environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Config file not found: %v", err)
	}

	log.Printf("Assistant Service loaded with environment: %s", ENV_RUNTIME.Environment)
}

// GetEnvironment returns current environment
func GetEnvironment() string {
	return ENV_RUNTIME.Environment
}

// IsProduction checks if running in production
func IsProduction() bool {
	return ENV_RUNTIME.Environment == "production" || ENV_RUNTIME.Environment == "release"
}

// IsDevelopment checks if running in development
func IsDevelopment() bool {
	return ENV_RUNTIME.Environment == "development"
}
