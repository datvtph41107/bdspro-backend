package _utils

import (
	_constants "common/constants"
	"context"

	"github.com/gin-gonic/gin"
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

func GetSessionIdFromContext(c context.Context) uint64 {
	sessionId := c.Value(_constants.SessionKey)
	if sessionId == nil {
		return 0
	}
	return sessionId.(uint64)
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

func GetPlanIdFromContext(ctx context.Context) uint64 {
	planId := ctx.Value(_constants.PlanIDKey)
	if planId == nil {
		return 0
	}
	return planId.(uint64)
}

// New functions for investment module
const (
	UserIDKey = "user_id"
)

func GetUserIDFromContext(ctx context.Context) uint64 {
	if userID, ok := ctx.Value(UserIDKey).(uint64); ok {
		return userID
	}
	return 0
}

func SetUserIDToContext(ctx context.Context, userID uint64) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// Helper function to create proper context from gin context for HTTP routes
// This bridges the gap between HTTP middleware and gRPC middleware patterns
func CreateContextFromGin(c *gin.Context) context.Context {
	ctx := c.Request.Context()

	// Get profile ID from gin context and set to proper context
	if profileId, exists := c.Get("profileId"); exists {
		if pid, ok := profileId.(uint64); ok {
			ctx = context.WithValue(ctx, _constants.ProfileIDKey, pid)
		}
	}

	// Get organization ID from gin context and set to proper context
	if organizationId, exists := c.Get("organizationId"); exists {
		if orgId, ok := organizationId.(*uint64); ok && orgId != nil {
			ctx = context.WithValue(ctx, _constants.OrganizationIDKey, *orgId)
		}
	}

	// Get plan ID if available
	if planId, exists := c.Get("planId"); exists {
		if pid, ok := planId.(*uint64); ok && pid != nil {
			ctx = context.WithValue(ctx, _constants.PlanIDKey, *pid)
		}
	}

	// Get auth ID if available
	if claims, exists := c.Get("claims"); exists {
		if principal, ok := claims.(*Principal); ok {
			ctx = context.WithValue(ctx, _constants.AuthIDKey, principal.UserID)
			ctx = context.WithValue(ctx, _constants.RoleKey, principal.Role)
			ctx = context.WithValue(ctx, _constants.TypeKey, principal.Type)
			ctx = context.WithValue(ctx, _constants.SessionKey, uint64(principal.Session))
		}
	}

	return ctx
}
