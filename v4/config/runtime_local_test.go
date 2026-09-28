package config_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"userv4/config"
)

func TestLoadRuntimeUsesLocalYMLWithoutEnv(t *testing.T) {
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

	for _, key := range []string{
		"DATABASE_URL", "USER_DATABASE_URL", "DATABASE_DSN",
		"REDIS_ADDRESS", "USER_REDIS_ADDRESS", "JWT_KEY_GENERATE",
		"QHPRO_CONFIG_FILE", "CONFIG_FILE", "ENV_RUNTIME", "QHPRO_ENVIRONMENT",
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
	if runtime.Redis.Address == "" {
		t.Fatal("redis address is empty")
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
}

func TestLoadRuntimeDatabaseDSNIgnoresEnvOverride(t *testing.T) {
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

	for _, key := range []string{"QHPRO_CONFIG_FILE", "CONFIG_FILE"} {
		t.Setenv(key, "")
	}
	t.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable")
	t.Setenv("USER_DATABASE_URL", "postgres://postgres:postgres@localhost:5432/user_service?sslmode=disable")

	runtime, err := config.LoadRuntime()
	if err != nil {
		t.Fatalf("LoadRuntime() error = %v", err)
	}
	if strings.Contains(runtime.Database.DSN, "localhost:5432/user_service") {
		t.Fatalf("database DSN came from .env override, want YAML local.yml value: %q", runtime.Database.DSN)
	}
	if !strings.Contains(runtime.Database.DSN, "db_user") {
		t.Fatalf("database DSN = %q, want local.yml db_user", runtime.Database.DSN)
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
