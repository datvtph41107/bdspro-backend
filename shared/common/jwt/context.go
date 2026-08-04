package _jwt

import "context"

type principalContextKey struct{}

/**
 * WithPrincipal stores the principal verified at the HTTP boundary.
 */
func WithPrincipal(ctx context.Context, principal *Principal) context.Context {
	if ctx == nil || principal == nil {
		return ctx
	}
	return context.WithValue(ctx, principalContextKey{}, principal)
}

/**
 * PrincipalFromRequestContext returns the verified HTTP principal.
 */
func PrincipalFromRequestContext(ctx context.Context) (*Principal, bool) {
	if ctx == nil {
		return nil, false
	}
	principal, ok := ctx.Value(principalContextKey{}).(*Principal)
	return principal, ok && principal != nil
}
