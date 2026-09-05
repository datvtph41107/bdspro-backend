package identity

import (
	"context"
	"errors"
	"strings"
)

type CallerKind string

const (
	CallerAnonymous       CallerKind = "anonymous"
	CallerUser            CallerKind = "user"
	CallerAPIKey          CallerKind = "api_key"
	CallerInternalService CallerKind = "internal_service"
)

// Caller classifies the logical source of a request.
type Caller struct {
	Kind      CallerKind
	APIKeyID  uint64
	APIKeyApp string
}

type callerContextKey struct{}

var (
	ErrInvalidCaller         = errors.New("caller classification is invalid")
	ErrCallerContextConflict = errors.New("caller classification conflicts with the existing context")
)

// IsValid reports whether the caller classification is internally consistent.
func (c Caller) IsValid() bool {
	switch c.Kind {
	case CallerAnonymous, CallerInternalService:
		return c.APIKeyID == 0 && strings.TrimSpace(c.APIKeyApp) == ""
	case CallerUser:
		return (c.APIKeyID == 0 && strings.TrimSpace(c.APIKeyApp) == "") ||
			(c.APIKeyID > 0 && validCallerText(c.APIKeyApp))
	case CallerAPIKey:
		return c.APIKeyID > 0 && validCallerText(c.APIKeyApp)
	default:
		return false
	}
}

// BindCaller validates and immutably binds one canonical Caller.
// Rebinding the same value is idempotent; rebinding a different value fails.
func BindCaller(ctx context.Context, caller Caller) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if !caller.IsValid() {
		return ctx, ErrInvalidCaller
	}

	caller.APIKeyApp = strings.TrimSpace(caller.APIKeyApp)
	if existing, ok := CallerFromContext(ctx); ok {
		if existing == caller {
			return ctx, nil
		}
		return ctx, ErrCallerContextConflict
	}

	return context.WithValue(ctx, callerContextKey{}, caller), nil
}

// CallerFromContext returns the canonical logical caller.
func CallerFromContext(ctx context.Context) (Caller, bool) {
	if ctx == nil {
		return Caller{}, false
	}
	caller, ok := ctx.Value(callerContextKey{}).(Caller)
	return caller, ok && caller.IsValid()
}

func validCallerText(value string) bool {
	if value == "" || strings.TrimSpace(value) != value || len(value) > 128 {
		return false
	}
	for _, character := range value {
		if character < 0x20 || character == 0x7f {
			return false
		}
	}
	return true
}
