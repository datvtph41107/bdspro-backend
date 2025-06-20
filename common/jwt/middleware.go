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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Xóa "Bearer " nếu có
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse và xác thực token
		claims, err := ParseJWT(tokenString)
		log.Printf("claims=%v, ignoreToken=%v", claims, ignoreToken)
		if err != nil && !ignoreToken {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
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

func JwtMiddleware(next http.Handler, publicPaths []string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Nếu path thuộc public, cho qua
		for _, p := range publicPaths {
			if strings.HasPrefix(path, p) {
				c.Next()
				return
			}
		}

		// Các path còn lại yêu cầu JWT
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			return
		}

		// Authorization: Bearer <token> → lấy ra phần token
		parts := strings.Split(authHeader, "Bearer ")

		var tokenString string
		if len(parts) != 2 {
			tokenString = parts[0]
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format"})
			// return
		} else {
			tokenString = strings.TrimSpace(parts[1])
		}

		// Parse và validate token
		// claims := jwt.MapClaims{}
		// token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// 	// TODO: thay secretKey bằng secret thực tế
		// 	return []byte("secretKey"), nil
		// })
		claims, err := ParseJWT(tokenString)
		log.Printf("err", err)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Nếu cần: inject thông tin user vào context Gin
		c.Set("userInfo", claims)
		c.Set("profileId", claims.ProfileId)
		c.Set("organizationId", claims.OrganizationId)

		c.Next() // Cho request đi tiếp

	})
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
		ctx = context.WithValue(ctx, authErrorKey, AuthError{Message: "missing token"})
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

// func respondWithError(w http.ResponseWriter, code int, message string) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(map[string]string{
// 		"error": message,
// 	})
// }

// func MiddlewareGrpcToken(publicPaths []string) func(ctx context.Context, r *http.ResponseWriter) metadata.MD {
// 	return func(ctx context.Context, r *http.ResponseWriter) metadata.MD {
// 		path := r.Q.Path

// 		md := metadata.MD{}
// 		// Nếu path thuộc public, cho qua
// 		for _, p := range publicPaths {
// 			if strings.HasPrefix(path, p) {
// 				return md
// 			}
// 		}

// 		// Các path còn lại yêu cầu JWT
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			respondWithError(r.ResponseWriter, http.StatusUnauthorized, "Missing Authorization header")
// 			return md
// 		}

// 		tokenString := r.Header.Get("Authorization")
// 		if tokenString != "" {
// 			// md.Append("authorization", tokenString)

// 			splitToken := strings.Split(tokenString, "Bearer ")
// 			var tokenStr string
// 			if len(splitToken) == 2 {
// 				tokenStr = strings.TrimSpace(splitToken[1])
// 			} else {
// 				tokenStr = strings.TrimSpace(splitToken[0])
// 			}
// 			claims, _ := ParseJWT(tokenStr)
// 			if claims != nil {
// 				md.Append("profileId", fmt.Sprintf("%d", claims.ProfileId))
// 				md.Append("planId", fmt.Sprintf("%d", *claims.PlanId))
// 			}

// 			// claims := jwt.MapClaims{}
// 			// token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
// 			// 	return []byte("secretKey"), nil
// 			// })
// 			// if err == nil && token.Valid {
// 			// 	if pid, ok := claims["profileId"].(float64); ok {
// 			// 		profileIdStr := fmt.Sprintf("%.0f", pid)
// 			// 		md.Append("profileid", profileIdStr)
// 			// 	}
// 			// }
// 		}

// 		md.Append("method", r.Method)
// 		log.Println("Metadata: ", md)
// 		return md
// 	}
// }

func ErrorToGatewayResponse(ctx context.Context, w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	log.Println("ErrorToGatewayResponse", err)

	// Nếu context có lỗi auth -> xử lý ở đây
	if ae, ok := ctx.Value(authErrorKey).(AuthError); ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{
			"error": ae.Message,
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
		"error": grpcErr.Message(),
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

// func ErrorToGatewayResponse(ctx context.Context, w http.ResponseWriter, err error) {
// 	st, ok := status.FromError(err)
// 	if !ok {
// 		http.Error(w, "Internal error", http.StatusInternalServerError)
// 		return
// 	}
// 	errorMap := make(map[string]string)

// 	for _, detail := range st.Details() {
// 		switch t := detail.(type) {
// 		case *pb_spec.FieldErrors:
// 			for k, v := range t.Errors {
// 				errorMap[k] = v
// 			}
// 		}
// 	}

// 	resp := map[string]interface{}{
// 		"code":    int(st.Code()),
// 		"message": st.Message(),
// 		"errors":  errorMap,
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(runtime.HTTPStatusFromCode(st.Code()))
// 	_ = json.NewEncoder(w).Encode(resp)

// 	// var specErr *_routes.Except
// 	// if errors.As(err, &specErr) {
// 	// 	w.Header().Set("Content-Type", "application/json")
// 	// 	w.WriteHeader(http.StatusOK) // hoặc `w.WriteHeader(specErr.Code)` nếu bạn muốn
// 	// 	json.NewEncoder(w).Encode(specErr)
// 	// } else {
// 	// 	w.WriteHeader(http.StatusInternalServerError)
// 	// 	_ = json.NewEncoder(w).Encode(map[string]interface{}{
// 	// 		"code":    500,
// 	// 		"message": err.Error(),
// 	// 	})
// 	// }
// }
