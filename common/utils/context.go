package _utils

import (
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
	profileId := ctx.Value("profileId")
	if profileId == nil {
		return 0
	}
	return profileId.(uint64)
}

func GetOrganizationIdFromContext(c context.Context) uint64 {
	organizationId := c.Value("profileId")
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
