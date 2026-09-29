package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type testConfig struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
	Name string `yaml:"name"`
}

func TestReadFileRequiresExplicitPath(t *testing.T) {
	_, err := ReadFile("")
	if err == nil || !strings.Contains(err.Error(), "config path is empty") {
		t.Fatalf("ReadFile error = %v", err)
	}
}

func TestReadFileReturnsPathAwareError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.yml")
	_, err := ReadFile(path)
	if err == nil || !strings.Contains(err.Error(), "read config "+path) {
		t.Fatalf("ReadFile error = %v", err)
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("ReadFile error = %v, want os.ErrNotExist", err)
	}
}

func TestDecodeYAMLRequiresOutput(t *testing.T) {
	if err := DecodeYAML([]byte("name: qhpro\n"), nil); err == nil {
		t.Fatal("DecodeYAML accepted nil output")
	}
}

func TestLoadYAMLDecodesTypedConfig(t *testing.T) {
	path := filepath.Join(t.TempDir(), "service.yml")
	if err := os.WriteFile(path, []byte("name: gateway\nserver:\n  port: 8080\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var cfg testConfig
	if err := LoadYAML(path, &cfg); err != nil {
		t.Fatalf("LoadYAML error = %v", err)
	}
	if cfg.Name != "gateway" || cfg.Server.Port != 8080 {
		t.Fatalf("config = %+v", cfg)
	}
}

func TestLoadYAMLReturnsPathAwareParseError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "broken.yml")
	if err := os.WriteFile(path, []byte("server: [\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var cfg testConfig
	err := LoadYAML(path, &cfg)
	if err == nil || !strings.Contains(err.Error(), "parse config "+path) {
		t.Fatalf("LoadYAML error = %v", err)
	}
}
