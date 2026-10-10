package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"userv4/config"

	"github.com/spf13/viper"
)

func chdirServiceRoot(t *testing.T) {
	t.Helper()
	serviceRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(serviceRoot); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
}

func clearConfigSelectors(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"QHPRO_CONFIG_FILE", "CONFIG_FILE", "ENV_RUNTIME", "QHPRO_ENVIRONMENT", "QHPRO_EXECUTION_MODE",
	} {
		t.Setenv(key, "")
	}
}

func TestLoadRuntimeUsesLocalYMLWithoutEnv(t *testing.T) {
	chdirServiceRoot(t)
	clearConfigSelectors(t)
	for _, key := range []string{
		"DATABASE_URL", "USER_DATABASE_URL", "DATABASE_DSN",
		"REDIS_ADDRESS", "USER_REDIS_ADDRESS", "REDIS_HOST", "REDIS_PASSWORD", "USER_REDIS_PASSWORD",
		"JWT_KEY_GENERATE",
		"RPC_AUTH_ADDRESS", "RPC_USER_ADDRESS",
	} {
		t.Setenv(key, "")
	}

	runtime, err := config.LoadRuntime()
	if err != nil {
		t.Fatalf("LoadRuntime() error = %v", err)
	}
	if runtime.Database.DSN == "" {
		t.Fatal("database DSN is empty")
	}
	if !strings.Contains(runtime.Redis.Address, "localhost:6379") && !strings.Contains(runtime.Redis.Address, "127.0.0.1:6379") {
		t.Fatalf("redis address = %q, want localhost from local.yml", runtime.Redis.Address)
	}
	if len(runtime.HTTP.AllowedOrigins) == 0 {
		t.Fatal("HTTP allowed_origins empty; want CORS from runtime.yml base")
	}
	found := false
	for _, origin := range runtime.HTTP.AllowedOrigins {
		if origin == "http://localhost:3000" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("allowed_origins = %#v, want http://localhost:3000 from runtime.yml", runtime.HTTP.AllowedOrigins)
	}
	if runtime.HTTPAddress != ":8288" {
		t.Fatalf("HTTPAddress = %q, want :8288 from merged local/runtime server.port", runtime.HTTPAddress)
	}
	if got := viper.GetString("rpc.auth.address"); got != "localhost:8201" {
		t.Fatalf("rpc.auth.address = %q, want local.yml value localhost:8201", got)
	}
	if got := viper.GetString("config-app.qhpro.domain_api"); !strings.Contains(got, "bdspro.com") {
		t.Fatalf("config-app.qhpro.domain_api = %q, want YAML value", got)
	}
}

func TestLoadRuntimeDatabaseDSNIgnoresEnvOverride(t *testing.T) {
	chdirServiceRoot(t)
	clearConfigSelectors(t)
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable")
	t.Setenv("USER_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable")

	runtime, err := config.LoadRuntime()
	if err != nil {
		t.Fatalf("LoadRuntime() error = %v", err)
	}
	if strings.Contains(runtime.Database.DSN, "localhost:5432/user_service") {
		t.Fatalf("database DSN came from .env override, want YAML local.yml value: %q", runtime.Database.DSN)
	}
	if !strings.Contains(runtime.Database.DSN, "db_v4") && !strings.Contains(runtime.Database.DSN, "db_user") {
		t.Fatalf("database DSN = %q, want local.yml db", runtime.Database.DSN)
	}
}

func TestLoadRuntimeRedisIgnoresEnvOverride(t *testing.T) {
	chdirServiceRoot(t)
	clearConfigSelectors(t)
	// Env points at a different host/port than local.yml.
	t.Setenv("USER_REDIS_ADDRESS", "10.255.255.1:6380")
	t.Setenv("REDIS_ADDRESS", "10.255.255.1:6380")
	t.Setenv("REDIS_HOST", "10.255.255.1:6380")

	runtime, err := config.LoadRuntime()
	if err != nil {
		t.Fatalf("LoadRuntime() error = %v", err)
	}
	if strings.Contains(runtime.Redis.Address, "10.255.255.1") {
		t.Fatalf("redis came from .env override, want YAML local.yml value: %q", runtime.Redis.Address)
	}
	if !strings.Contains(runtime.Redis.Address, "localhost:6379") && !strings.Contains(runtime.Redis.Address, "127.0.0.1:6379") {
		t.Fatalf("redis address = %q, want local.yml redis host", runtime.Redis.Address)
	}
}

func TestLoadPropertiesIgnoresEnvForJWTAndRPC(t *testing.T) {
	chdirServiceRoot(t)
	clearConfigSelectors(t)
	t.Setenv("JWT_KEY_GENERATE", "env-jwt-key-that-must-not-win-over-yaml!!!!!!")
	t.Setenv("RPC_AUTH_ADDRESS", "127.0.0.1:9999")
	t.Setenv("RPC_USER_ADDRESS", "127.0.0.1:9998")
	t.Setenv("QHPRO_USER_OUTBOUND_MESSAGING_ENABLED", "true")

	if err := config.LoadProperties(); err != nil {
		t.Fatalf("LoadProperties() error = %v", err)
	}
	if config.Properties.JWT.SecretKey != "key-generate-jwt-token-authentication" {
		t.Fatalf("jwt = %q, want local.yml value", config.Properties.JWT.SecretKey)
	}
	if got := viper.GetString("rpc.auth.address"); got != "localhost:8201" {
		t.Fatalf("rpc.auth.address = %q, want YAML not env", got)
	}
	if got := viper.GetString("rpc.user.address"); got != "localhost:8201" {
		t.Fatalf("rpc.user.address = %q, want YAML not env", got)
	}
	if config.Properties.OutboundMessaging.Enabled {
		t.Fatal("outbound_messaging.enabled came from env; want YAML false")
	}
}

func TestDescribeDatabaseDSNRedactsPassword(t *testing.T) {
	got := config.DescribeDatabaseDSN(
		"host=14.225.210.29 user=hieubap password=secret dbname=db_user port=5432 sslmode=disable",
	)
	if strings.Contains(got, "secret") || strings.Contains(got, "password") {
		t.Fatalf("DescribeDatabaseDSN leaked credentials: %q", got)
	}
	if !strings.Contains(got, "host=14.225.210.29") ||
		!strings.Contains(got, "dbname=db_user") ||
		!strings.Contains(got, "user=hieubap") {
		t.Fatalf("DescribeDatabaseDSN = %q", got)
	}
}
