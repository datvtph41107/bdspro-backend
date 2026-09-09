package client

import (
	"fmt"
	"strings"
)

// validateProviderCall is deliberately pure: process configuration is resolved
// once by assistant/config and injected into each provider client.
func validateProviderCall(provider, mode, apiKey string) error {
	mode = strings.ToLower(strings.TrimSpace(mode))
	if mode == "" {
		mode = "live"
	}
	switch mode {
	case "stub", "disabled", "off":
		return fmt.Errorf("%s provider is disabled by provider mode=%s", provider, mode)
	case "live":
		if strings.TrimSpace(apiKey) == "" {
			return fmt.Errorf("%s API key is required in live provider mode", provider)
		}
		return nil
	default:
		return fmt.Errorf("unsupported provider mode=%q (want stub or live)", mode)
	}
}
