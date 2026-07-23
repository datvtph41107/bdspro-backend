package _jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	sharepb "pb/types/shared"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Kiểm tra xem đường dẫn có trong danh sách publicRoutes không
func isPublicRoute(pubRoutes []string, path string) bool {
	for _, route := range pubRoutes {
		if matchRoute(route, path) {
			return true
		}
	}
	return false
}

// matchRoute kiểm tra xem path có match với pattern không
// Hỗ trợ:
// - * : wildcard cho bất kỳ ký tự nào (có thể rỗng)
// - : : parameter cho số hoặc chữ (không được rỗng)
func matchRoute(pattern, path string) bool {
	// Nếu pattern không chứa wildcard hoặc parameter, dùng prefix matching
	if !strings.Contains(pattern, "*") && !strings.Contains(pattern, ":") {
		return strings.HasPrefix(path, pattern)
	}

	// Split pattern và path thành segments
	patternSegments := strings.Split(strings.Trim(pattern, "/"), "/")
	pathSegments := strings.Split(strings.Trim(path, "/"), "/")

	// Loại bỏ segments rỗng cuối cùng từ path
	for len(pathSegments) > 0 && pathSegments[len(pathSegments)-1] == "" {
		pathSegments = pathSegments[:len(pathSegments)-1]
	}

	// Loại bỏ segments rỗng cuối cùng từ pattern, nhưng giữ lại wildcard
	for len(patternSegments) > 0 && patternSegments[len(patternSegments)-1] == "" {
		patternSegments = patternSegments[:len(patternSegments)-1]
	}

	// Nếu pattern kết thúc bằng wildcard (*), cho phép path có thêm 1 segment
	if len(patternSegments) > 0 && patternSegments[len(patternSegments)-1] == "*" {
		// Loại bỏ wildcard cuối cùng
		patternSegments = patternSegments[:len(patternSegments)-1]
		// Path phải có đúng số segments bằng pattern + 1 (cho wildcard)
		if len(pathSegments) != len(patternSegments)+1 {
			return false
		}
		// Chỉ kiểm tra các segments tương ứng
		for i := 0; i < len(patternSegments); i++ {
			patternSeg := patternSegments[i]
			pathSeg := pathSegments[i]

			// Nếu pattern segment là parameter (:), match nếu không rỗng
			if patternSeg == ":" {
				if pathSeg == "" {
					return false
				}
				continue
			}

			// Nếu pattern segment không phải parameter, phải khớp chính xác
			if patternSeg != pathSeg {
				return false
			}
		}
		return true
	}

	// Nếu số lượng segments không khớp, không match
	if len(patternSegments) != len(pathSegments) {
		return false
	}

	// Kiểm tra từng segment
	for i, patternSeg := range patternSegments {
		pathSeg := pathSegments[i]

		// Nếu pattern segment là wildcard (*), match bất kỳ gì (kể cả rỗng)
		if patternSeg == "*" {
			continue
		}

		// Nếu pattern segment là parameter (:), match nếu không rỗng
		if patternSeg == ":" {
			if pathSeg == "" {
				return false
			}
			continue
		}

		// Nếu pattern segment không phải wildcard hay parameter, phải khớp chính xác
		if patternSeg != pathSeg {
			return false
		}
	}

	return true
}

// Middleware để xác thực JWT
func JWTAuthMiddleware(pubRoutes []string, tempRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		apiKey := c.GetHeader("API-KEY")
		if apiKey != "" {
			c.Set("api-key", apiKey)
		}
		// fmt.Printf("request: %s _ api-key: %s\n", c.Request.URL.Path, apiKey)
		// fmt.Printf("pubRoutes: %v\n", pubRoutes)
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
		// log.Printf("claims=%v, ignoreToken=%v, err=%v", claims, ignoreToken, err)
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

			// Kiểm tra nếu là TEMP token thì chỉ cho phép gọi tempRoutes
			if claims.Type == string(TempToken) {
				isTempRoute := false
				for _, route := range tempRoutes {
					if matchRoute(route, c.Request.URL.Path) {
						isTempRoute = true
						break
					}
				}
				if !isTempRoute {
					c.JSON(http.StatusForbidden, gin.H{
						"code":    403,
						"message": "Token không đúng",
					})
					c.Abort()
					return
				}
			}
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
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
