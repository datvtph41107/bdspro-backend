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
	w.WriteHeader(httpStatus)
	//
	resp := map[string]interface{}{
		"code":    grpcErr.Code(),
		"message": grpcErr.Message(),
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
