package config

import (
	shared "common/pkg/config"
	"fmt"
)

type GatewayConfig struct {
	shared.BaseConfig   `mapstructure:",squash"`
	shared.SharedConfig `mapstructure:",squash"`
	Gateway             GatewaySpecific `mapstructure:"gateway"`
}

type GatewaySpecific struct {
	GlobalRateLimit int      `mapstructure:"global_rate_limit"   default:"1000" reload:"hot"`
	PerIPRateLimit  int      `mapstructure:"per_ip_rate_limit"   default:"50"   reload:"hot"`
	UpstreamTimeout int      `mapstructure:"upstream_timeout_sec" default:"10"  reload:"cold"`
	AllowedOrigins  []string `mapstructure:"allowed_origins"`
	PublicRoutes    []string `mapstructure:"public_routes"        reload:"hot"`
	TempRoutes      []string `mapstructure:"temp_routes"`
}

func (c *GatewayConfig) Validate() error {
	if c.Gateway.PerIPRateLimit > c.Gateway.GlobalRateLimit {
		return fmt.Errorf("per_ip_rate_limit (%d) cannot exceed global_rate_limit (%d)",
			c.Gateway.PerIPRateLimit, c.Gateway.GlobalRateLimit)
	}
	if c.Feature.MaintenanceMode && c.Server.Port == 80 {
		return fmt.Errorf("maintenance mode cannot run on privileged port 80")
	}
	return nil
}

func MustLoad() *GatewayConfig {
	var cfg GatewayConfig
	shared.MustLoad("gateway", &cfg)
	return &cfg
}
