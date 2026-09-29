package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestNewRuntimeMaterializesProviderConfig(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	t.Setenv(providerModeEnv, "stub")
	t.Setenv("OPENAI_API_KEY", "env-openai-key")

	viper.Set("app.port.grpc", "8218")
	viper.Set("openai.api_key", "yaml-openai-key")
	viper.Set("openai.base_url", "https://example.invalid/v1")
	viper.Set("openai.model", "test-model")
	viper.Set("openai.timeout_seconds", 45)
	viper.Set("openai.max_tokens", 1234)

	runtime, err := NewRuntime()
	if err != nil {
		t.Fatalf("NewRuntime() error = %v", err)
	}
	if runtime.GRPCPort != "8218" {
		t.Fatalf("GRPCPort = %q, want 8218", runtime.GRPCPort)
	}
	if runtime.ProviderMode != "stub" {
		t.Fatalf("ProviderMode = %q, want stub", runtime.ProviderMode)
	}
	if runtime.OpenAI.APIKey != "env-openai-key" {
		t.Fatalf("OpenAI.APIKey did not prefer environment override")
	}
	if runtime.OpenAI.BaseURL != "https://example.invalid/v1" {
		t.Fatalf("OpenAI.BaseURL = %q", runtime.OpenAI.BaseURL)
	}
	if runtime.OpenAI.Model != "test-model" {
		t.Fatalf("OpenAI.Model = %q", runtime.OpenAI.Model)
	}
	if runtime.OpenAI.Timeout != 45*time.Second {
		t.Fatalf("OpenAI.Timeout = %s", runtime.OpenAI.Timeout)
	}
	if runtime.OpenAI.MaxTokens != 1234 {
		t.Fatalf("OpenAI.MaxTokens = %d", runtime.OpenAI.MaxTokens)
	}
}

func TestNewRuntimeRejectsUnknownProviderMode(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv(providerModeEnv, "maybe")

	if _, err := NewRuntime(); err == nil {
		t.Fatal("NewRuntime() accepted unknown provider mode")
	}
}
