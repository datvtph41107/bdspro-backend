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
	ErrAPIKeyInvalid     = errors.New("api key is invalid")
	ErrAPIKeyUnavailable = errors.New("api key verifier is unavailable")
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

/**
 * APIKeyFromRequest returns one normalized optional key.
 */
func APIKeyFromRequest(r *http.Request) (string, bool, error) {
	if r == nil {
		return "", false, nil
	}

	values := make([]string, 0, 3)
	for _, header := range []string{APIKeyHeader, LegacyAPIKeyHeader, "api-key"} {
		for _, value := range r.Header.Values(header) {
			value = strings.TrimSpace(value)
			if value != "" {
				values = append(values, value)
			}
		}
	}
	if len(values) == 0 {
		return "", false, nil
	}

	first := values[0]
	for _, value := range values[1:] {
		if value != first {
			return "", false, ErrAPIKeyInvalid
		}
	}
	if !isValidAPIKey(first) {
		return "", false, ErrAPIKeyInvalid
	}
	return first, true, nil
}

func isValidAPIKey(value string) bool {
	if len(value) < 16 || len(value) > MaxAPIKeyLength {
		return false
	}
	for _, r := range value {
		if r <= 0x20 || r == 0x7f {
			return false
		}
	}
	return true
}
