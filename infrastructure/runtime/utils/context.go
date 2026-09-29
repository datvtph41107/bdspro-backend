package _utils

import (
	_enums "common/domain/enum"
	_jwt "common/jwt"
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
	profileId := ctx.Value(_enums.ProfileIDKey)
	if profileId == nil {
		return 0
	}
	return profileId.(uint64)
}

func GetOriginIdFromContext(ctx context.Context) uint64 {
	profileId := ctx.Value(_enums.OriginIDKey)
	if profileId == nil {
		return 0
	}
	return profileId.(uint64)
}

func GetrRoleWithContext(ctx context.Context) uint64 {
	// clms := ctx.Value("claims")
	// if claims, ok := clms.(*Principal); ok {
	// 	return claims.ProfileId
	// }
	// return 0
	roleKey := ctx.Value(_enums.RoleKey)
	if roleKey == nil {
		return 0
	}
	return roleKey.(uint64)
}

func GetOrganizationIdFromContext(c context.Context) uint64 {
	// todo: đổi sang profile sau khi xong
	organizationId := c.Value(_enums.OrganizationIDKey)
	if organizationId == nil {
		return 0
	}
	return organizationId.(uint64)
}

func GetSessionIdFromContext(c context.Context) uint64 {
	sessionId := c.Value(_enums.SessionKey)
	if sessionId == nil {
		return 0
	}
	return sessionId.(uint64)
}

func GetAuthIdFromContext(c context.Context) uint64 {
	authId := c.Value(_enums.AuthIDKey)
	if authId == nil {
		return 0
	}
	return authId.(uint64)
}

func GetDeviceIdFromContext(c context.Context) string {
	deviceID := c.Value(_enums.DeviceIDKey)
	if deviceID == nil {
		return ""
	}
	if id, ok := deviceID.(string); ok {
		return id
	}
	return ""
}

func GetRoleIdsFromContext(ctx context.Context) []uint64 {
	roleIds := ctx.Value(_enums.RoleIDsKey)
	if roleIds == nil {
		return nil
	}
	if ids, ok := roleIds.([]uint64); ok {
		return ids
	}
	return nil
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
	originId := GetOriginIdFromContext(ctx)
	organizationId := GetOrganizationIdFromContext(ctx)
	cloneCtx = context.WithValue(cloneCtx, _enums.ProfileIDKey, profileId)
	cloneCtx = context.WithValue(cloneCtx, _enums.OriginIDKey, originId)
	cloneCtx = context.WithValue(cloneCtx, _enums.OrganizationIDKey, organizationId)
	roleIds := GetRoleIdsFromContext(ctx)
	if len(roleIds) > 0 {
		cloneCtx = context.WithValue(cloneCtx, _enums.RoleIDsKey, roleIds)
	}
	return cloneCtx
}

func GetPlanIdFromContext(ctx context.Context) uint64 {
	planId := ctx.Value(_enums.PlanIDKey)
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
			ctx = context.WithValue(ctx, _enums.ProfileIDKey, pid)
		}
	}

	if originId, exists := c.Get("originId"); exists {
		if oid, ok := originId.(uint64); ok {
			ctx = context.WithValue(ctx, _enums.OriginIDKey, oid)
		}
	}

	// Get organization ID from gin context and set to proper context
	if organizationId, exists := c.Get("organizationId"); exists {
		if orgId, ok := organizationId.(*uint64); ok && orgId != nil {
			ctx = context.WithValue(ctx, _enums.OrganizationIDKey, *orgId)
		}
	}

	// Get plan ID if available
	if planId, exists := c.Get("planId"); exists {
		if pid, ok := planId.(*uint64); ok && pid != nil {
			ctx = context.WithValue(ctx, _enums.PlanIDKey, *pid)
		}
	}

	if deviceId, exists := c.Get("device-id"); exists {
		if did, ok := deviceId.(string); ok {
			ctx = context.WithValue(ctx, _enums.DeviceIDKey, did)
		}
	}

	// Get auth ID if available
	if claims, exists := c.Get("claims"); exists {
		if principal, ok := claims.(*_jwt.Principal); ok {
			ctx = context.WithValue(ctx, _enums.AuthIDKey, principal.AuthID)
			ctx = context.WithValue(ctx, _enums.RoleKey, principal.Role)
			ctx = context.WithValue(ctx, _enums.TypeKey, principal.Type)
			ctx = context.WithValue(ctx, _enums.SessionKey, uint64(principal.Session))
			if len(principal.RoleIds) > 0 {
				ctx = context.WithValue(ctx, _enums.RoleIDsKey, principal.RoleIds)
			}
		}
	}

	return ctx
}
