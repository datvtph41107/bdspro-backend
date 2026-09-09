package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func prepareRuntimeTest(t *testing.T) {
	t.Helper()
	viper.Reset()
	t.Cleanup(viper.Reset)

	for _, key := range []string{
		"QHPRO_PERMISSION_REFRESH_INTERVAL",
		"QHPRO_PERMISSION_MAX_STALENESS",
		"QHPRO_PERMISSION_ROLE_CACHE_TTL",
		"QHPRO_PERMISSION_REQUEST_TIMEOUT",
		"QHPRO_SERVICE_ID",
		"SERVICE_NAME",
		"QHPRO_INTERNAL_METADATA_SECRET",
		"QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS",
		"QHPRO_INTERNAL_METADATA_MAX_AGE",
		"QHPRO_INTERNAL_METADATA_CLOCK_SKEW",
		"QHPRO_TRUSTED_METADATA_MODE",
	} {
		t.Setenv(key, "")
	}
}

func TestNewRuntimeMaterializesTypedOwnership(t *testing.T) {
	prepareRuntimeTest(t)

	viper.Set("server.tcp_port", 8216)
	viper.Set("rpc.user.address", "user:8201")
	viper.Set("permission.refresh_interval", "30s")
	viper.Set("permission.max_staleness", "2m")
	viper.Set("permission.role_cache_ttl", "15s")
	viper.Set("permission.request_timeout", "2s")

	t.Setenv("QHPRO_PERMISSION_REFRESH_INTERVAL", "45s")
	t.Setenv("QHPRO_SERVICE_ID", "auth-service")
	t.Setenv("QHPRO_INTERNAL_METADATA_SECRET", "test-secret")
	t.Setenv("QHPRO_TRUSTED_METADATA_MODE", "enforce")

	runtime, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	if runtime.GRPCPort != 8216 {
		t.Fatalf("GRPCPort = %d", runtime.GRPCPort)
	}
	if runtime.UserRPCTarget != "user:8201" {
		t.Fatalf("UserRPCTarget = %q", runtime.UserRPCTarget)
	}
	if runtime.Permission.RefreshInterval != 45*time.Second {
		t.Fatalf("RefreshInterval = %s", runtime.Permission.RefreshInterval)
	}
	if runtime.Permission.MaxStaleness != 2*time.Minute {
		t.Fatalf("MaxStaleness = %s", runtime.Permission.MaxStaleness)
	}
	if runtime.Permission.RoleCacheTTL != 15*time.Second {
		t.Fatalf("RoleCacheTTL = %s", runtime.Permission.RoleCacheTTL)
	}
	if runtime.Permission.RequestTimeout != 2*time.Second {
		t.Fatalf("RequestTimeout = %s", runtime.Permission.RequestTimeout)
	}
	if runtime.Transport.ServiceAssertion.ServiceID != "auth-service" {
		t.Fatalf(
			"transport service ID = %q",
			runtime.Transport.ServiceAssertion.ServiceID,
		)
	}
	if !runtime.Transport.RequireServiceAssertion {
		t.Fatal("trusted metadata enforce mode must require service assertion")
	}
}

func TestNewRuntimePreservesPermissionFallbacks(t *testing.T) {
	prepareRuntimeTest(t)

	viper.Set("server.tcp_port", 8216)
	viper.Set("rpc.user.address", "localhost:8201")
	t.Setenv("QHPRO_PERMISSION_REQUEST_TIMEOUT", "invalid")

	runtime, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}

	if runtime.Permission.RefreshInterval != defaultPermissionRefreshInterval {
		t.Fatalf("RefreshInterval = %s", runtime.Permission.RefreshInterval)
	}
	if runtime.Permission.MaxStaleness != defaultPermissionMaxStaleness {
		t.Fatalf("MaxStaleness = %s", runtime.Permission.MaxStaleness)
	}
	if runtime.Permission.RoleCacheTTL != defaultPermissionRoleCacheTTL {
		t.Fatalf("RoleCacheTTL = %s", runtime.Permission.RoleCacheTTL)
	}
	if runtime.Permission.RequestTimeout != defaultPermissionRequestTimeout {
		t.Fatalf("RequestTimeout = %s", runtime.Permission.RequestTimeout)
	}
}

func TestNewRuntimeRejectsInvalidProcessTopology(t *testing.T) {
	prepareRuntimeTest(t)

	viper.Set("server.tcp_port", 0)
	viper.Set("rpc.user.address", "user:8201")
	if _, err := NewRuntime(); err == nil {
		t.Fatal("zero gRPC port must fail")
	}

	viper.Set("server.tcp_port", 8216)
	viper.Set("rpc.user.address", "")
	if _, err := NewRuntime(); err == nil {
		t.Fatal("missing User RPC target must fail")
	}
}
