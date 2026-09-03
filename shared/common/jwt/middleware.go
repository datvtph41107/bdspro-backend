package _jwt

import (
	_callercontext "common/callercontext"
	_identitycontext "common/identitycontext"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	sharepb "pb/types/shared"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
			claims, tokenErr = ParseJWT(tokenString)
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
			boundActorCtx, bindActorErr := _identitycontext.BindActor(requestCtx, ActorFromPrincipal(claims))
			if bindActorErr != nil {
				message := "invalid actor"
				if errors.Is(bindActorErr, _identitycontext.ErrActorContextConflict) {
					message = "actor context conflict"
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":    http.StatusInternalServerError,
					"message": message,
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
		boundCallerContext, bindErr := _callercontext.BindCaller(c.Request.Context(), resolution.Caller)
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

func ErrorToGatewayResponse(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("ErrorToGatewayResponse", err)

	// Nếu là lỗi gRPC bình thường
	grpcErr, ok := status.FromError(err)
	if !ok {
		if message, isConn := friendlyUpstreamBusyMessage(err.Error()); isConn {
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":    http.StatusServiceUnavailable,
				"message": message,
			})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}

	if ok {
		for _, d := range grpcErr.Details() {
			switch info := d.(type) {
			case *sharepb.ErrorResponse:
				fmt.Printf("Custom code = %d, message = %s\n", info.Code, info.Message)
				// Business error trong body; HTTP 200 để FE đọc res.code (tránh axios reject).
				// (NoRoute từng để status=404 → login "404" dù auth thành công.)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"code":    info.Code,
					"message": info.Message,
					"second":  info.Second,
				})
				return
			}
		}
	}

	httpStatus := httpStatusFromCode(grpcErr.Code())
	message := grpcErr.Message()
	if friendly, isConn := friendlyUpstreamBusyMessage(message); isConn {
		message = friendly
		httpStatus = http.StatusServiceUnavailable
	} else if friendly, isConn := friendlyUpstreamBusyMessage(err.Error()); isConn {
		message = friendly
		httpStatus = http.StatusServiceUnavailable
	}
	w.WriteHeader(httpStatus)
	resp := map[string]interface{}{
		"code":    grpcErr.Code(),
		"message": message,
	}

	for _, d := range grpcErr.Details() {
		switch br := d.(type) {

		case *errdetails.BadRequest:
			var errors []map[string]string
			for _, v := range br.FieldViolations {
				errors = append(errors, map[string]string{
					"field": v.Field,
					"error": v.Description,
				})
			}
			resp["errors"] = errors
		}
	}

	json.NewEncoder(w).Encode(resp)
}

var (
	dialTCPPortRe           = regexp.MustCompile(`(?i)dial tcp (?:\[[^\]]+\]|[^:\s]+):(\d{2,5})`)
	connectionRefusedPortRe = regexp.MustCompile(`:(\d{2,5}): connect: connection refused`)
)

func isUpstreamConnectionError(msg string) bool {
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "error while dialing") ||
		strings.Contains(lower, "connection refused") ||
		strings.Contains(lower, "dial tcp") ||
		strings.Contains(lower, "all subconns are in transientfailure") ||
		(strings.Contains(lower, "connection error") && strings.Contains(lower, "transport"))
}

func extractUpstreamPort(msg string) string {
	if m := dialTCPPortRe.FindStringSubmatch(msg); len(m) == 2 {
		return m[1]
	}
	if m := connectionRefusedPortRe.FindStringSubmatch(msg); len(m) == 2 {
		return m[1]
	}
	return ""
}

func friendlyUpstreamBusyMessage(msg string) (string, bool) {
	if !isUpstreamConnectionError(msg) {
		return msg, false
	}
	if port := extractUpstreamPort(msg); port != "" {
		return fmt.Sprintf("máy chủ đang bận vui lòng thử lại sau (%s)", port), true
	}
	return "máy chủ đang bận vui lòng thử lại sau", true
}

func httpStatusFromCode(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
