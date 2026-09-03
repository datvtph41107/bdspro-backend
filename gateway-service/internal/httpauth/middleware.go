package httpauth

import (
	"common/identity"
	"common/jwtverify"
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware(publicRoutes, tempRoutes []Route, options ...Option) gin.HandlerFunc {
	cfg := config{}
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
		providerCredential := matchesRequest(
			cfg.providerCredentialRoutes,
			c.Request.Method,
			c.Request.URL.Path,
		)
		ignoreToken := providerCredential
		for _, route := range publicRoutes {
			if MatchRequest(route, c.Request.Method, c.Request.URL.Path) {
				ignoreToken = true
				break
			}
		}
		var principal *jwtverify.Principal
		rawToken := strings.TrimSpace(c.GetHeader("Authorization"))
		if rawToken != "" && !providerCredential {
			rawToken = strings.TrimSpace(strings.TrimPrefix(rawToken, "Bearer "))
			if cfg.tokenVerifier == nil {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"code": http.StatusServiceUnavailable, "message": "token verifier unavailable"})
				return
			}
			var err error
			principal, err = cfg.tokenVerifier.Parse(rawToken)
			if err != nil || principal == nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "Invalid or expired token"})
				return
			}
			switch jwtverify.TokenType(principal.Type) {
			case jwtverify.AccessToken:
			case jwtverify.TempToken:
				allowed := false
				for _, route := range tempRoutes {
					if MatchRequest(route, c.Request.Method, c.Request.URL.Path) {
						allowed = true
						break
					}
				}
				if !allowed {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": http.StatusForbidden, "message": "Token không đúng"})
					return
				}
			default:
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"code": http.StatusUnauthorized, "message": "Token type is not allowed"})
				return
			}
			requestCtx := WithPrincipal(c.Request.Context(), principal)
			bound, err := identity.BindActor(requestCtx, ActorFromPrincipal(principal))
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "actor context conflict"})
				return
			}
			c.Request = c.Request.WithContext(bound)
			c.Set("claims", principal)
			c.Set("organizationId", principal.OrganizationID)
			c.Set("profileId", principal.ProfileID)
			c.Set("planId", principal.PlanID)
		}
		resolution, err := ResolveHTTPCaller(c.Request.Context(), c.Request, principal, ignoreToken, cfg.apiKeyVerifier)
		if err != nil {
			statusCode, message := http.StatusUnauthorized, "Invalid caller classification"
			switch {
			case errors.Is(err, ErrAPIKeyHeaderInvalid):
				statusCode, message = http.StatusBadRequest, "Invalid API key header"
			case errors.Is(err, ErrAPIKeyInvalid):
				message = "API key verification failed"
			case errors.Is(err, ErrAPIKeyUnavailable), errors.Is(err, context.DeadlineExceeded):
				statusCode, message = http.StatusServiceUnavailable, "API key verification unavailable"
			case errors.Is(err, context.Canceled):
				c.Abort()
				return
			case errors.Is(err, ErrAuthenticationRequired):
				message = "Authorization header missing"
			case errors.Is(err, ErrTokenTypeNotAllowed):
				message = "Token type is not allowed"
			}
			c.AbortWithStatusJSON(statusCode, gin.H{"code": statusCode, "message": message})
			return
		}
		bound, err := identity.BindCaller(c.Request.Context(), resolution.Caller)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "caller context conflict"})
			return
		}
		c.Request = c.Request.WithContext(bound)
		c.Next()
	}
}

func matchesRequest(routes []Route, method, path string) bool {
	for _, route := range routes {
		if MatchRequest(route, method, path) {
			return true
		}
	}
	return false
}
