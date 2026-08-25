package httpauth

import (
	"common/jwtverify"
	"context"
)

type principalContextKey struct{}

func WithPrincipal(ctx context.Context, p *jwtverify.Principal) context.Context {
	if ctx == nil || p == nil {
		return ctx
	}

	copy := *p
	copy.OrganizationID = cloneUint64(p.OrganizationID)
	copy.PlanID = cloneUint64(p.PlanID)
	copy.RoleIDs = append([]uint64(nil), p.RoleIDs...)

	return context.WithValue(ctx, principalContextKey{}, &copy)
}

func PrincipalFromContext(ctx context.Context) (*jwtverify.Principal, bool) {
	if ctx == nil {
		return nil, false
	}

	p, ok := ctx.Value(principalContextKey{}).(*jwtverify.Principal)
	if !ok || p == nil {
		return nil, false
	}

	copy := *p
	copy.OrganizationID = cloneUint64(p.OrganizationID)
	copy.PlanID = cloneUint64(p.PlanID)
	copy.RoleIDs = append([]uint64(nil), p.RoleIDs...)

	return &copy, true
}
