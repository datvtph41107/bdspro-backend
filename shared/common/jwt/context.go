package _jwt

import "context"

type principalContextKey struct{}

/**
 * WithPrincipal stores a defensive copy verified at the HTTP boundary.
 */
func WithPrincipal(ctx context.Context, principal *Principal) context.Context {
	if ctx == nil || principal == nil {
		return ctx
	}
	return context.WithValue(ctx, principalContextKey{}, clonePrincipal(principal))
}

/**
 * PrincipalFromRequestContext returns a defensive verified principal copy.
 */
func PrincipalFromRequestContext(ctx context.Context) (*Principal, bool) {
	if ctx == nil {
		return nil, false
	}
	principal, ok := ctx.Value(principalContextKey{}).(*Principal)
	if !ok || principal == nil {
		return nil, false
	}
	return clonePrincipal(principal), true
}

func clonePrincipal(principal *Principal) *Principal {
	if principal == nil {
		return nil
	}
	clone := *principal
	clone.OrganizationId = cloneUint64Pointer(principal.OrganizationId)
	clone.PlanId = cloneUint64Pointer(principal.PlanId)
	clone.RoleIds = append([]uint64(nil), principal.RoleIds...)
	clone.Audience = append(principal.Audience[:0:0], principal.Audience...)
	if principal.ExpiresAt != nil {
		value := *principal.ExpiresAt
		clone.ExpiresAt = &value
	}
	if principal.NotBefore != nil {
		value := *principal.NotBefore
		clone.NotBefore = &value
	}
	if principal.IssuedAt != nil {
		value := *principal.IssuedAt
		clone.IssuedAt = &value
	}
	return &clone
}
