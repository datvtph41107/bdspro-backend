package config

import "time"

type SharedConfig struct {
	Auth     AuthConfig       `mapstructure:"auth"`
	Services ServiceDiscovery `mapstructure:"services"`
}

type AuthConfig struct {
	JWTSecret       string        `mapstructure:"jwt_secret"        validate:"required" sensitive:"true"`
	JWTIssuer       string        `mapstructure:"jwt_issuer"        validate:"required"`
	AccessTokenTTL  time.Duration `mapstructure:"access_token_ttl"  default:"15m"`
	RefreshTokenTTL time.Duration `mapstructure:"refresh_token_ttl" default:"168h"`
}

type ServiceDiscovery struct {
	Default   ServiceEndpoint            `mapstructure:"default"`
	Endpoints map[string]ServiceEndpoint `mapstructure:"endpoints"`
}

type ServiceEndpoint struct {
	Addresses []string      `mapstructure:"addresses" validate:"required,min=1"`
	Timeout   time.Duration `mapstructure:"timeout"   default:"3s"`
	Retries   int           `mapstructure:"retries"   default:"1"`
	Backoff   time.Duration `mapstructure:"backoff"   default:"100ms"`
	// Circuit breaker settings
	MaxFailures int           `mapstructure:"max_failures" default:"5"`
	TimeoutOpen time.Duration `mapstructure:"timeout_open" default:"30s"`
}
