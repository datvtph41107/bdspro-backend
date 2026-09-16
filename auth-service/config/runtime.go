package config

import (
	"common/configloader"
	qhprorpc "common/rpc"
	"common/rpcenv"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	defaultPermissionRefreshInterval = 30 * time.Second
	defaultPermissionMaxStaleness    = 2 * time.Minute
	defaultPermissionRoleCacheTTL    = 15 * time.Second
	defaultPermissionRequestTimeout  = 2 * time.Second
)

type PermissionRuntime struct {
	RefreshInterval time.Duration
	MaxStaleness    time.Duration
	RoleCacheTTL    time.Duration
	RequestTimeout  time.Duration
}

type Runtime struct {
	GRPCPort      int
	UserRPCTarget string
	Permission    PermissionRuntime
	Transport     qhprorpc.TransportConfig
}

// NewRuntime materializes Auth process configuration from the Viper registry
// established by LoadProperties plus process environment values. Call it once
// from the process composition root and inject the returned immutable values.
func NewRuntime() (Runtime, error) {
	port := viper.GetInt("server.tcp_port")
	if port <= 0 || port > 65535 {
		return Runtime{}, fmt.Errorf("invalid auth gRPC port %d", port)
	}

	userTarget, err := configloader.RequiredString("rpc.user.address")
	if err != nil {
		return Runtime{}, err
	}

	return Runtime{
		GRPCPort:      port,
		UserRPCTarget: userTarget,
		Permission: PermissionRuntime{
			RefreshInterval: durationFromRuntime(
				"QHPRO_PERMISSION_REFRESH_INTERVAL",
				"permission.refresh_interval",
				defaultPermissionRefreshInterval,
			),
			MaxStaleness: durationFromRuntime(
				"QHPRO_PERMISSION_MAX_STALENESS",
				"permission.max_staleness",
				defaultPermissionMaxStaleness,
			),
			RoleCacheTTL: durationFromRuntime(
				"QHPRO_PERMISSION_ROLE_CACHE_TTL",
				"permission.role_cache_ttl",
				defaultPermissionRoleCacheTTL,
			),
			RequestTimeout: durationFromRuntime(
				"QHPRO_PERMISSION_REQUEST_TIMEOUT",
				"permission.request_timeout",
				defaultPermissionRequestTimeout,
			),
		},
		Transport: rpcenv.LoadTransportConfig(),
	}, nil
}

func durationFromRuntime(
	envKey string,
	configKey string,
	fallback time.Duration,
) time.Duration {
	if value := strings.TrimSpace(os.Getenv(envKey)); value != "" {
		return parseRuntimeDuration(envKey, value, fallback)
	}
	if value := strings.TrimSpace(viper.GetString(configKey)); value != "" {
		return parseRuntimeDuration(configKey, value, fallback)
	}
	return fallback
}

func parseRuntimeDuration(
	key string,
	value string,
	fallback time.Duration,
) time.Duration {
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		slog.Warn(
			"auth runtime duration invalid; using fallback",
			slog.String("key", key),
			slog.String("value", value),
			slog.Duration("fallback", fallback),
		)
		return fallback
	}
	return parsed
}
