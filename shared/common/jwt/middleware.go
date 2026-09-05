package _jwt

import (
	"common/identity"

	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Middleware để xác thực JWT
func JWTAuthMiddleware(pubRoutes []string, tempRoutes []string, options ...AuthOption) gin.HandlerFunc {
	cfg := authConfig{}
	for _, option := range options {
		if option != nil {
			option(&cfg)
		}
	}

	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
		ignoreToken := pubRoutes != nil && isPublicRoute(pubRoutes, c.Request.URL.Path)

		var claims *Principal
		if tokenString != "" {
			tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
			var tokenErr error
			if cfg.tokenVerifier != nil {
				claims, tokenErr = cfg.tokenVerifier.Parse(tokenString)
			} else {
				// Compatibility path for services not yet migrated to explicit JWT verification.
				claims, tokenErr = ParseJWT(tokenString)
			}
			if tokenErr != nil || claims == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Invalid or expired token",
				})
				return
			}
		}

		if claims != nil {
			requestCtx := WithPrincipal(c.Request.Context(), claims)
			boundActorCtx, bindActorErr := identity.BindActor(requestCtx, ActorFromPrincipal(claims))
			if bindActorErr != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": "actor context conflict",
				})
				return
			}
			c.Request = c.Request.WithContext(boundActorCtx)
			c.Set("claims", claims)
			c.Set("organizationId", claims.OrganizationId)
			c.Set("profileId", claims.ProfileId)
			c.Set("planId", claims.PlanId)

			if claims.Type == string(TempToken) {
				isTempRoute := false
				for _, route := range tempRoutes {
					if matchRoute(route, c.Request.URL.Path) {
						isTempRoute = true
						break
					}
				}
				if !isTempRoute {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
						"code":    http.StatusForbidden,
						"message": "Token không đúng",
					})
					return
				}
			}
		}

		resolution, callerErr := ResolveHTTPCaller(
			c.Request.Context(),
			c.Request,
			claims,
			ignoreToken,
			cfg.apiKeyVerifier,
		)
		if callerErr != nil {
			statusCode := http.StatusUnauthorized
			message := "Invalid caller classification"
			switch {
			case errors.Is(callerErr, ErrAPIKeyHeaderInvalid):
				statusCode = http.StatusBadRequest
				message = "Invalid API key header"
			case errors.Is(callerErr, ErrAPIKeyInvalid):
				message = "API key verification failed"
			case errors.Is(callerErr, ErrAPIKeyUnavailable), errors.Is(callerErr, context.DeadlineExceeded):
				statusCode = http.StatusServiceUnavailable
				message = "API key verification unavailable"
			case errors.Is(callerErr, context.Canceled):
				c.Abort()
				return
			case errors.Is(callerErr, ErrAuthenticationRequired):
				message = "Authorization header missing"
			}
			c.AbortWithStatusJSON(statusCode, gin.H{
				"code":    statusCode,
				"message": message,
			})
			return
		}
		boundCallerContext, bindErr := identity.BindCaller(c.Request.Context(), resolution.Caller)
		if bindErr != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "caller context conflict",
			})
			return
		}
		c.Request = c.Request.WithContext(boundCallerContext)
		c.Next()
	}
}
