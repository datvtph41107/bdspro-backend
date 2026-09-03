package config

import "testing"

func TestParseConfigRejectsMalformedYAML(t *testing.T) {
	if _, err := parseConfig([]byte("weights: [")); err == nil {
		t.Fatal("expected malformed embedded YAML to return an error")
	}
}

func TestEmbeddedConfigParses(t *testing.T) {
	cfg, err := parseConfig(configYAML)
	if err != nil {
		t.Fatalf("embedded resolver config must remain valid: %v", err)
	}
	if cfg == nil {
		t.Fatal("expected resolver config")
	}
}
