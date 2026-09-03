package identity

import (
	"context"
	"errors"
)

const maxServiceIDLength = 128

// ServiceCaller identifies the authenticated service that directly called
// the current service.
type ServiceCaller struct {
	ServiceID string
}

type serviceCallerContextKey struct{}

var (
	ErrInvalidServiceCaller         = errors.New("service caller identity is invalid")
	ErrServiceCallerContextConflict = errors.New("service caller identity conflicts with the existing context")
)

// IsValid reports whether the service caller identity is usable.
func (c ServiceCaller) IsValid() bool {
	if len(c.ServiceID) < 1 || len(c.ServiceID) > maxServiceIDLength {
		return false
	}
	for _, r := range c.ServiceID {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' || r == ':' {
			continue
		}
		return false
	}
	return true
}

// BindServiceCaller validates and immutably binds one authenticated service caller.
// Rebinding the same value is idempotent; rebinding a different value fails.
func BindServiceCaller(ctx context.Context, caller ServiceCaller) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if !caller.IsValid() {
		return ctx, ErrInvalidServiceCaller
	}
	if existing, ok := ServiceCallerFromContext(ctx); ok {
		if existing == caller {
			return ctx, nil
		}
		return ctx, ErrServiceCallerContextConflict
	}
	return context.WithValue(ctx, serviceCallerContextKey{}, caller), nil
}

// ServiceCallerFromContext returns the authenticated immediate service caller.
func ServiceCallerFromContext(ctx context.Context) (ServiceCaller, bool) {
	if ctx == nil {
		return ServiceCaller{}, false
	}
	caller, ok := ctx.Value(serviceCallerContextKey{}).(ServiceCaller)
	return caller, ok && caller.IsValid()
}
