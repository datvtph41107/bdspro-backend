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
				abortProblem(c, http.StatusServiceUnavailable, "auth.token_verifier_unavailable", "token verifier unavailable")
				return
			}
			var err error
			principal, err = cfg.tokenVerifier.Parse(rawToken)
			if err != nil || principal == nil {
				abortProblem(c, http.StatusUnauthorized, "auth.invalid_token", "invalid or expired token")
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
					abortProblem(c, http.StatusForbidden, "auth.temporary_token_forbidden", "temporary token is not allowed for this route")
					return
				}
			default:
				abortProblem(c, http.StatusUnauthorized, "auth.token_type_not_allowed", "token type is not allowed")
				return
			}
			requestCtx := WithPrincipal(c.Request.Context(), principal)
			bound, err := identity.BindActor(requestCtx, ActorFromPrincipal(principal))
			if err != nil {
				abortInternalProblem(c, "auth.actor_context_conflict")
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
			statusCode := http.StatusUnauthorized
			code := "auth.caller_invalid"
			message := "invalid caller classification"
			switch {
			case errors.Is(err, ErrAPIKeyHeaderInvalid):
				statusCode, code, message = http.StatusBadRequest, "auth.api_key_header_invalid", "invalid API key header"
			case errors.Is(err, ErrAPIKeyInvalid):
				code, message = "auth.api_key_invalid", "API key verification failed"
			case errors.Is(err, ErrAPIKeyUnavailable), errors.Is(err, context.DeadlineExceeded):
				statusCode, code, message = http.StatusServiceUnavailable, "auth.api_key_unavailable", "API key verification unavailable"
			case errors.Is(err, context.Canceled):
				c.Abort()
				return
			case errors.Is(err, ErrAuthenticationRequired):
				code, message = "auth.authentication_required", "authorization header missing"
			case errors.Is(err, ErrTokenTypeNotAllowed):
				code, message = "auth.token_type_not_allowed", "token type is not allowed"
			}
			abortProblem(c, statusCode, code, message)
			return
		}
		bound, err := identity.BindCaller(c.Request.Context(), resolution.Caller)
		if err != nil {
			abortInternalProblem(c, "auth.caller_context_conflict")
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
