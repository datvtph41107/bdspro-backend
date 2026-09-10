package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"user/config"
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
}
