package _jwt

import (
	"common/identity"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// JWTAuthHTTPMiddleware authenticates HTTP requests using JWT and/or API key,
// binding Actor/Caller onto request context. It does not use Gin.
func JWTAuthHTTPMiddleware(pubRoutes []string, tempRoutes []string, options ...AuthOption) func(http.Handler) http.Handler {
	cfg := authConfig{}
	for _, option := range options {
		if option != nil {
			option(&cfg)
		}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			req, status, message, silent := authenticateRequest(r, pubRoutes, tempRoutes, cfg)
			if silent {
				return
			}
			if status != 0 {
				writeAuthJSON(w, status, message)
				return
			}
			next.ServeHTTP(w, req)
		})
	}
}

func writeAuthJSON(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"code":    status,
		"message": message,
	})
}

// authenticateRequest validates JWT/API-key and returns an updated request with
// identity bound on context. On failure status != 0; silent is true for canceled contexts.
func authenticateRequest(r *http.Request, pubRoutes, tempRoutes []string, cfg authConfig) (*http.Request, int, string, bool) {
	tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
	ignoreToken := pubRoutes != nil && isPublicRoute(pubRoutes, r.URL.Path)

	var claims *Principal
	if tokenString != "" {
		tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
		var tokenErr error
		if cfg.tokenVerifier != nil {
			claims, tokenErr = cfg.tokenVerifier.Parse(tokenString)
		} else {
			claims, tokenErr = ParseJWT(tokenString)
		}
		if tokenErr != nil || claims == nil {
			return r, http.StatusUnauthorized, "Invalid or expired token", false
		}
	}

	req := r
	if claims != nil {
		requestCtx := WithPrincipal(req.Context(), claims)
		boundActorCtx, bindActorErr := identity.BindActor(requestCtx, ActorFromPrincipal(claims))
		if bindActorErr != nil {
			return r, http.StatusInternalServerError, "actor context conflict", false
		}
		req = req.WithContext(boundActorCtx)

		if claims.Type == string(TempToken) {
			isTempRoute := false
			for _, route := range tempRoutes {
				if matchRoute(route, req.URL.Path) {
					isTempRoute = true
					break
				}
			}
			if !isTempRoute {
				return r, http.StatusForbidden, "Token không đúng", false
			}
		}
	}

	resolution, callerErr := ResolveHTTPCaller(
		req.Context(),
		req,
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
			return r, 0, "", true
		case errors.Is(callerErr, ErrAuthenticationRequired):
			message = "Authorization header missing"
		}
		return r, statusCode, message, false
	}
	boundCallerContext, bindErr := identity.BindCaller(req.Context(), resolution.Caller)
	if bindErr != nil {
		return r, http.StatusInternalServerError, "caller context conflict", false
	}
	return req.WithContext(boundCallerContext), 0, "", false
}
