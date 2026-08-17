package _jwt

import (
	_callercontext "common/callercontext"

	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	sharepb "pb/types/shared"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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

		rawAPIKey, hasAPIKey, apiKeyErr := APIKeyFromRequest(c.Request)
		if apiKeyErr != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Invalid API key header",
			})
			return
		}

		var verifiedAPIKey VerifiedAPIKey
		if hasAPIKey {
			if cfg.apiKeyVerifier == nil {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"code":    http.StatusServiceUnavailable,
					"message": "API key verification unavailable",
				})
				return
			}

			verifiedAPIKey, apiKeyErr = cfg.apiKeyVerifier.VerifyAPIKey(c.Request.Context(), rawAPIKey)
			if apiKeyErr != nil {
				statusCode := http.StatusUnauthorized
				if errors.Is(apiKeyErr, ErrAPIKeyUnavailable) {
					statusCode = http.StatusServiceUnavailable
				}
				c.AbortWithStatusJSON(statusCode, gin.H{
					"code":    statusCode,
					"message": "API key verification failed",
				})
				return
			}
			if verifiedAPIKey.ID == 0 || strings.TrimSpace(verifiedAPIKey.AppName) == "" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Invalid API key",
				})
				return
			}

			c.Request.Header.Set(APIKeyHeader, rawAPIKey)
			c.Set("api-key", rawAPIKey)
			c.Set("api-key-id", verifiedAPIKey.ID)
			c.Set("api-key-app", verifiedAPIKey.AppName)
		}

		tokenString := strings.TrimSpace(c.GetHeader("Authorization"))
		ignoreToken := pubRoutes != nil && isPublicRoute(pubRoutes, c.Request.URL.Path)
		if tokenString == "" && !ignoreToken && !hasAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Authorization header missing",
			})
			return
		}

		var claims *Principal
		var tokenErr error
		if tokenString != "" {
			tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))
			claims, tokenErr = ParseJWT(tokenString)
		}
		if tokenErr != nil && !ignoreToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid or expired token",
			})
			return
		}

		caller := _callercontext.Caller{Kind: _callercontext.KindAnonymous}
		if claims != nil && tokenErr == nil {
			c.Request = c.Request.WithContext(WithPrincipal(c.Request.Context(), claims))
			c.Set("claims", claims)
			c.Set("organizationId", claims.OrganizationId)
			c.Set("profileId", claims.ProfileId)
			c.Set("planId", claims.PlanId)
			caller.Kind = _callercontext.KindUser

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
		} else if hasAPIKey {
			caller.Kind = _callercontext.KindAPIKey
		}

		if hasAPIKey {
			caller.APIKeyID = verifiedAPIKey.ID
			caller.APIKeyApp = verifiedAPIKey.AppName
		}
		if !caller.IsValid() {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid caller classification",
			})
			return
		}
		c.Request = c.Request.WithContext(_callercontext.WithCaller(c.Request.Context(), caller))
		c.Next()
	}
}

type AuthError struct {
	Message string
}

type authErrorKeyType string

const authErrorKey authErrorKeyType = "auth_error"

func ExtractMetadataFromRequest(ctx context.Context, r *http.Request) metadata.MD {
	md := metadata.MD{}

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		// Gắn lỗi vào context để xử lý sau
		_ = context.WithValue(ctx, authErrorKey, AuthError{
			Message: "Missing Authorization header",
		})
		return md
	}

	// Parse token từ "Bearer xxx"
	splitToken := strings.Split(tokenString, "Bearer ")
	var tokenStr string
	if len(splitToken) == 2 {
		tokenStr = strings.TrimSpace(splitToken[1])
	} else {
		tokenStr = strings.TrimSpace(splitToken[0])
	}

	claims, err := ParseJWT(tokenStr)
	if err != nil || claims == nil {
		_ = context.WithValue(ctx, authErrorKey, AuthError{Message: "invalid or expired token"})
		return md
	}

	md.Append("profileId", fmt.Sprintf("%d", claims.ProfileId))
	if claims.OrganizationId != nil {
		md.Append("organizationId", fmt.Sprintf("%d", *claims.OrganizationId))
	}
	if claims.Session != 0 {
		md.Append("session", fmt.Sprintf("%d", claims.Session))
	}
	if claims.Role != "" {
		md.Append("role", claims.Role)
	}
	if claims.Type != "" {
		md.Append("type", claims.Type)
	}
	if claims.AuthID != 0 {
		md.Append("authId", fmt.Sprintf("%d", claims.AuthID))
	}
	if claims.PlanId != nil {
		md.Append("planId", fmt.Sprintf("%d", *claims.PlanId))
	}
	if claims.PlanFrom != nil && !claims.PlanFrom.Time.IsZero() {
		md.Append("planFrom", claims.PlanFrom.Time.Format(time.RFC3339))
	}
	md.Append("method", r.Method)

	return md
}

func ErrorToGatewayResponse(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("ErrorToGatewayResponse", err)

	// Nếu context có lỗi auth -> xử lý ở đây
	if ae, ok := ctx.Value(authErrorKey).(AuthError); ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"code":    "401",
			"message": ae.Message,
		})
		return
	}

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
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
