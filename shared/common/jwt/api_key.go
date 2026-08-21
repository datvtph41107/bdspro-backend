package _jwt

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

const (
	APIKeyHeader       = "X-API-Key"
	LegacyAPIKeyHeader = "API-KEY"
	MaxAPIKeyLength    = 256
)

var (
	ErrAPIKeyHeaderInvalid = errors.New("api key header is invalid")
	ErrAPIKeyInvalid       = errors.New("api key is invalid")
	ErrAPIKeyUnavailable   = errors.New("api key verifier is unavailable")
)

/**
 * VerifiedAPIKey is the non-secret identity returned by the API key authority.
 */
type VerifiedAPIKey struct {
	ID      uint64
	AppName string
}

/**
 * APIKeyVerifier validates a raw key against its authority.
 */
type APIKeyVerifier interface {
	VerifyAPIKey(ctx context.Context, rawKey string) (VerifiedAPIKey, error)
}

type authConfig struct {
	apiKeyVerifier APIKeyVerifier
	tokenVerifier  TokenVerifier
}

/**
 * AuthOption configures HTTP authentication without changing route contracts.
 */
type AuthOption func(*authConfig)

/**
 * WithAPIKeyVerifier enables verified API key callers.
 */
func WithAPIKeyVerifier(verifier APIKeyVerifier) AuthOption {
	return func(cfg *authConfig) {
		cfg.apiKeyVerifier = verifier
	}
}

// WithTokenVerifier installs an explicit JWT verification dependency.
// Gateway uses this to avoid process-global Viper configuration.
func WithTokenVerifier(verifier TokenVerifier) AuthOption {
	return func(cfg *authConfig) {
		cfg.tokenVerifier = verifier
	}
}

/**
 * APIKeyFromRequest returns one normalized optional key.
 */
func APIKeyFromRequest(r *http.Request) (string, bool, error) {
	if r == nil {
		return "", false, nil
	}

	values := make([]string, 0, 2)
	present := false
	for _, header := range []string{APIKeyHeader, LegacyAPIKeyHeader} {
		rawValues, exists := r.Header[http.CanonicalHeaderKey(header)]
		if !exists {
			continue
		}
		present = true
		for _, rawValue := range rawValues {
			value := strings.TrimSpace(rawValue)
			if value == "" {
				return "", false, ErrAPIKeyHeaderInvalid
			}
			values = append(values, value)
		}
	}
	if !present {
		return "", false, nil
	}
	if len(values) != 1 || !isValidAPIKey(values[0]) {
		return "", false, ErrAPIKeyHeaderInvalid
	}
	return values[0], true, nil
}

func isValidAPIKey(value string) bool {
	if len(value) < 16 || len(value) > MaxAPIKeyLength {
		return false
	}
	for _, r := range value {
		if r < 0x21 || r > 0x7e || r == ',' {
			return false
		}
	}
	return true
}
