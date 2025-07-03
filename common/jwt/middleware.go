package _jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Kiểm tra xem đường dẫn có trong danh sách publicRoutes không
func isPublicRoute(pubRoutes []string, path string) bool {
	for _, route := range pubRoutes {
		if strings.HasPrefix(path, route) {
			return true
		}
	}
	return false
}

// Middleware để xác thực JWT
func JWTAuthMiddleware(pubRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Lấy token từ header Authorization
		tokenString := c.GetHeader("Authorization")
		ignoreToken := false
		if pubRoutes != nil && isPublicRoute(pubRoutes, c.Request.URL.Path) {
			ignoreToken = true
		}
		if tokenString == "" && !ignoreToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header missing",
			})
			c.Abort()
			return
		}

		// Xóa "Bearer " nếu có
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse và xác thực token
		claims, err := ParseJWT(tokenString)
		log.Printf("claims=%v, ignoreToken=%v, err=%v", claims, ignoreToken, err)
		if err != nil && !ignoreToken {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// tx := c.MustGet("db").(*gorm.DB) // Lấy DB từ context
		// tx = tx.WithContext(context.WithValue(tx.Statement.Context, "profileId", claims.ProfileId))
		// c.Set("db", tx)

		// Gắn claims vào context để các handler sau có thể sử dụng
		if claims != nil {
			c.Set("claims", claims)
			c.Set("organizationId", claims.OrganizationId)
			c.Set("profileId", claims.ProfileId)
			c.Set("planId", claims.PlanId)
		}

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
		ctx = context.WithValue(ctx, authErrorKey, AuthError{
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
		ctx = context.WithValue(ctx, authErrorKey, AuthError{Message: "invalid or expired token"})
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
			"code":    "500",
			"message": ae.Message,
		})
		return
	}

	// Nếu là lỗi gRPC bình thường
	grpcErr, ok := status.FromError(err)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Internal Server Error",
		})
		return
	}

	httpStatus := httpStatusFromCode(grpcErr.Code())
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]string{
		"code":    "500",
		"message": grpcErr.Message(),
	})
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
