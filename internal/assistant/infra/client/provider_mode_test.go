package client

import (
	"strings"
	"testing"
)

func TestValidateProviderCallFailsClosed(t *testing.T) {
	t.Run("stub never calls external provider", func(t *testing.T) {
		err := validateProviderCall("openai", "stub", "secret")
		if err == nil || !strings.Contains(err.Error(), "disabled") {
			t.Fatalf("validateProviderCall() error = %v, want disabled", err)
		}
	})

	t.Run("live requires key", func(t *testing.T) {
		err := validateProviderCall("gemini", "live", "")
		if err == nil || !strings.Contains(err.Error(), "key is required") {
			t.Fatalf("validateProviderCall() error = %v, want missing key", err)
		}
	})

	t.Run("live accepts injected key", func(t *testing.T) {
		if err := validateProviderCall("deepseek", "live", "secret"); err != nil {
			t.Fatalf("validateProviderCall() error = %v", err)
		}
	})

	t.Run("unknown mode is rejected", func(t *testing.T) {
		if err := validateProviderCall("openai", "maybe", "secret"); err == nil {
			t.Fatal("validateProviderCall() accepted unknown mode")
		}
	})
}
