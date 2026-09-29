// Package rpcenv loads the process environment contract used to configure
// shared RPC service assertions.
package rpcenv

import (
	"common/identity"
	"common/rpc"
	"encoding/json"
	"os"
	"strings"
	"time"
)

// LoadTransportConfig resolves the QHPRO RPC service-assertion environment
// contract into the canonical transport configuration.
func LoadTransportConfig() rpc.TransportConfig {
	verificationSecrets, configured := verificationSecretsFromEnv(
		"QHPRO_INTERNAL_METADATA_VERIFICATION_KEYS",
	)

	return rpc.TransportConfig{
		ServiceAssertion: rpc.ServiceAssertionConfig{
			ServiceID: firstNonEmpty(
				os.Getenv("QHPRO_SERVICE_ID"),
				os.Getenv("SERVICE_NAME"),
			),
			Secret: strings.TrimSpace(
				os.Getenv("QHPRO_INTERNAL_METADATA_SECRET"),
			),
			VerificationSecrets:        verificationSecrets,
			VerificationKeysConfigured: configured,
			MaxAge: durationFromEnv(
				"QHPRO_INTERNAL_METADATA_MAX_AGE",
				30*time.Second,
			),
			ClockSkew: durationFromEnv(
				"QHPRO_INTERNAL_METADATA_CLOCK_SKEW",
				5*time.Second,
			),
		},
		RequireServiceAssertion: strings.EqualFold(
			strings.TrimSpace(
				os.Getenv("QHPRO_TRUSTED_METADATA_MODE"),
			),
			"enforce",
		),
	}
}

func durationFromEnv(
	key string,
	fallback time.Duration,
) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			return value
		}
	}
	return ""
}

func verificationSecretsFromEnv(
	key string,
) (map[string]string, bool) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil, false
	}

	var values map[string]string
	if err := json.Unmarshal([]byte(raw), &values); err != nil {
		// Environment explicitly configured a keyring but it is invalid.
		// Keep configured=true so verification cannot silently fall back to
		// the shared secret.
		return nil, true
	}

	result := make(map[string]string, len(values))
	for serviceID, secret := range values {
		serviceID = strings.TrimSpace(serviceID)
		secret = strings.TrimSpace(secret)

		if !(identity.ServiceCaller{
			ServiceID: serviceID,
		}).IsValid() || secret == "" {
			continue
		}

		result[serviceID] = secret
	}

	if len(result) == 0 {
		return nil, true
	}
	return result, true
}
