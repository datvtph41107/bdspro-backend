package _utils

import (
	_constants "common/constants"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type Principal struct {
	UserID    uint64  `json:"sub"`
	ProfileId uint64  `json:"profile"`
	PlanId    *uint64 `json:"plan"`
	Type      string  `json:"type"`
	Role      string  `json:"role"`
	Session   uint    `json:"session"`
	jwt.RegisteredClaims
}

func GetProfileIdWithContext(ctx context.Context) uint64 {
	// clms := ctx.Value("claims")
	// if claims, ok := clms.(*Principal); ok {
	// 	return claims.ProfileId
	// }
	// return 0
	profileId := ctx.Value(_constants.ProfileIDKey)
	if profileId == nil {
		return 0
	}
	return profileId.(uint64)
}

func GetOrganizationIdFromContext(c context.Context) uint64 {
	// todo: đổi sang profile sau khi xong
	organizationId := c.Value(_constants.OrganizationIDKey)
	if organizationId == nil {
		return 0
	}
	return organizationId.(uint64)
}

// func GetProfileIdFromContext(ctx context.Context) uint64 {
// 	return ctx.Value("profileId").(uint64)
// 	// if claims, ok := clms.(*Principal); ok {
// 	// 	return claims.ProfileId
// 	// }
// 	// return 0
// }

func CloneContext(ctx context.Context) context.Context {
	cloneCtx := context.Background()
	profileId := GetProfileIdWithContext(ctx)
	organizationId := GetOrganizationIdFromContext(ctx)
	cloneCtx = context.WithValue(cloneCtx, _constants.ProfileIDKey, profileId)
	cloneCtx = context.WithValue(cloneCtx, _constants.OrganizationIDKey, organizationId)
	return cloneCtx
}
