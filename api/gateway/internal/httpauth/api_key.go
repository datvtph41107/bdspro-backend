package httpauth

import (
	"common/jwtverify"
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
	ErrAPIKeyHeaderInvalid    = errors.New("api key header is invalid")
	ErrAPIKeyInvalid          = errors.New("api key is invalid")
	ErrAPIKeyUnavailable      = errors.New("api key verifier is unavailable")
	ErrAuthenticationRequired = errors.New("authentication is required")
	ErrCallerClassification   = errors.New("caller classification is invalid")
	ErrTokenTypeNotAllowed    = errors.New("jwt token type is not allowed for this route")
)

type VerifiedAPIKey struct {
	ID      uint64
	AppName string
}

type APIKeyVerifier interface {
	VerifyAPIKey(context.Context, string) (VerifiedAPIKey, error)
}

type config struct {
	apiKeyVerifier           APIKeyVerifier
	tokenVerifier            jwtverify.Verifier
	providerCredentialRoutes []Route
}

type Option func(*config)

func WithAPIKeyVerifier(v APIKeyVerifier) Option {
	return func(c *config) { c.apiKeyVerifier = v }
}

func WithTokenVerifier(v jwtverify.Verifier) Option {
	return func(c *config) { c.tokenVerifier = v }
}

// WithProviderCredentialRoutes identifies endpoints whose Authorization value
// belongs to an external provider protocol, not to QHPRO JWT authentication.
// The raw header remains available to the transport adapter for digesting.
func WithProviderCredentialRoutes(routes []Route) Option {
	return func(c *config) {
		c.providerCredentialRoutes = append([]Route(nil), routes...)
	}
}

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
		for _, raw := range rawValues {
			value := strings.TrimSpace(raw)
			if value == "" {
				return "", false, ErrAPIKeyHeaderInvalid
			}
			values = append(values, value)
		}
	}
	if !present {
		return "", false, nil
	}
	if len(values) != 1 || !validAPIKey(values[0]) {
		return "", false, ErrAPIKeyHeaderInvalid
	}

	return values[0], true, nil
}

func validAPIKey(value string) bool {
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
