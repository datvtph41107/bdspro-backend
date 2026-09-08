package client

import (
	"fmt"
	"os"
	"strings"
)

const providerModeEnv = "QHPRO_AI_PROVIDER_MODE"

// validateProviderCall chặn outbound request khi develop dùng stub và
// buộc runtime live phải inject key, không bao giờ commit secret vào Git.
func validateProviderCall(provider, apiKey string) error {
	mode := strings.ToLower(strings.TrimSpace(os.Getenv(providerModeEnv)))
	if mode == "" {
		mode = "live"
	}
	switch mode {
	case "stub", "disabled", "off":
		return fmt.Errorf("%s provider is disabled by %s=%s", provider, providerModeEnv, mode)
	case "live":
		if strings.TrimSpace(apiKey) == "" {
			return fmt.Errorf("%s API key is required in live provider mode", provider)
		}
		return nil
	default:
		return fmt.Errorf("unsupported %s=%q (want stub or live)", providerModeEnv, mode)
	}
}
