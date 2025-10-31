package env

type Runtime struct {
	Environment string
	ConfigPath  string
}

var ENV_RUNTIME *Runtime

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
