package config

import "time"

// shared/pkg/config/health.go
type HealthConfig struct {
	LivenessPath  string        `mapstructure:"liveness_path"  validate:"required"`
	ReadinessPath string        `mapstructure:"readiness_path" validate:"required"`
	CheckInterval time.Duration `mapstructure:"check_interval"`
	// Upstream health checks
	CheckUpstreams  bool          `mapstructure:"check_upstreams"`
	UpstreamTimeout time.Duration `mapstructure:"upstream_timeout"`
}
