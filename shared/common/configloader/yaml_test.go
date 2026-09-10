package configloader

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func TestConfigFilePathDefaultsToLocal(t *testing.T) {
	if got := configFilePath(""); got != "config/local.yml" {
		t.Fatalf("configFilePath(\"\") = %q", got)
	}
	if got := configFilePath(" develop "); got != "config/develop.yml" {
		t.Fatalf("configFilePath(develop) = %q", got)
	}
}

func TestLoadYMLFileReturnsReadError(t *testing.T) {
	withWorkingDir(t, t.TempDir(), func() {
		err := LoadYMLFile("missing")
		if err == nil || !strings.Contains(err.Error(), "read config config/missing.yml") {
			t.Fatalf("LoadYMLFile() error = %v, want read error", err)
		}
	})
}

func TestLoadYMLFileHonorsOperatorConfigOverride(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "mounted.yml")
	if err := os.WriteFile(path, []byte("runtime:\n  mode: mounted\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	withWorkingDir(t, t.TempDir(), func() {
		t.Setenv("QHPRO_CONFIG_FILE", path)
		if err := LoadYMLFile("production"); err != nil {
			t.Fatalf("LoadYMLFile() with override: %v", err)
		}
		if got := viper.GetString("runtime.mode"); got != "mounted" {
			t.Fatalf("runtime.mode = %q, want mounted", got)
		}
	})
}

func TestLoadYMLFileReturnsParseError(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "config", "local.yml"), []byte("server: [\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	withWorkingDir(t, root, func() {
		err := LoadYMLFile("")
		if err == nil || !strings.Contains(err.Error(), "parse config config/local.yml") {
			t.Fatalf("LoadYMLFile() error = %v, want parse error", err)
		}
	})
}

func TestLoadYMLFileLoadsPresentConfig(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "config", "local.yml"), []byte("server:\n  port: 8000\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	withWorkingDir(t, root, func() {
		if err := LoadYMLFile(""); err != nil {
			t.Fatalf("LoadYMLFile() error = %v", err)
		}
	})
}

func TestLoadRuntimeYMLMergesLocalOverlay(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "config", "runtime.yml"), []byte("database:\n  dsn: runtime-dsn\nredis:\n  host: runtime-redis\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "config", "local.yml"), []byte("database:\n  dsn: local-dsn\njwt:\n  key-generate: local-jwt-key-at-least-32-characters\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	withWorkingDir(t, root, func() {
		t.Setenv(EnvironmentKey, "")
		t.Setenv(ExecutionModeKey, "")
		t.Setenv(legacyProfileKey, "")
		t.Setenv("QHPRO_CONFIG_FILE", "")
		t.Setenv("CONFIG_FILE", "")
		t.Setenv("JWT_KEY_GENERATE", "")
		t.Setenv("DATABASE_DSN", "")
		t.Setenv("DATABASE_URL", "")
		viper.Reset()
		if _, err := LoadRuntimeYML(); err != nil {
			t.Fatalf("LoadRuntimeYML() error = %v", err)
		}
		if got := viper.GetString("database.dsn"); got != "local-dsn" {
			t.Fatalf("database.dsn = %q, want local-dsn", got)
		}
		if got := viper.GetString("redis.host"); got != "runtime-redis" {
			t.Fatalf("redis.host = %q, want runtime-redis", got)
		}
		if got := viper.GetString("jwt.key-generate"); got != "local-jwt-key-at-least-32-characters" {
			t.Fatalf("jwt.key-generate = %q", got)
		}
	})
}

func TestLoadRuntimeYMLSkipsMissingLocalOverlay(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "config"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "config", "runtime.yml"), []byte("database:\n  dsn: runtime-dsn\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	withWorkingDir(t, root, func() {
		t.Setenv(EnvironmentKey, "")
		t.Setenv(ExecutionModeKey, "")
		t.Setenv(legacyProfileKey, "")
		t.Setenv("QHPRO_CONFIG_FILE", "")
		t.Setenv("CONFIG_FILE", "")
		viper.Reset()
		if _, err := LoadRuntimeYML(); err != nil {
			t.Fatalf("LoadRuntimeYML() error = %v", err)
		}
		if got := viper.GetString("database.dsn"); got != "runtime-dsn" {
			t.Fatalf("database.dsn = %q, want runtime-dsn", got)
		}
	})
}

func withWorkingDir(t *testing.T, dir string, run func()) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Chdir(old); err != nil {
			t.Fatalf("restore working directory: %v", err)
		}
	}()
	run()
}
