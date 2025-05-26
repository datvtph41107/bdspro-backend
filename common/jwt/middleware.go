package _jwt

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	pb_spec "pb/types/spec"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// JWTMiddleware kiểm tra JWT từ request header
// func JWTMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
// 			return
// 		}

// 		// Kiểm tra header có dạng "Bearer <token>"
// 		parts := strings.Split(authHeader, " ")
// 		if len(parts) != 2 || parts[0] != "Bearer" {
// 			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString := parts[1]
// 		claims, err := ValidateJWT(tokenString)
// 		if err != nil {
// 			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
// 			return
// 		}

// 		// Gán claims vào context để các handler khác có thể sử dụng
// 		r = r.WithContext(r.Context()) // Có thể sử dụng context để lưu userID nếu cần

// 		// fmt.Printf("Authenticated user ID: %d\n", claims.UserID)
// 		next.ServeHTTP(w, r)
// 	})
// }

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
		if pubRoutes != nil && isPublicRoute(pubRoutes, c.Request.URL.Path) && tokenString == "" {
			c.Next()
			return
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Xóa "Bearer " nếu có
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse và xác thực token
		claims, err := ParseJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// tx := c.MustGet("db").(*gorm.DB) // Lấy DB từ context
		// tx = tx.WithContext(context.WithValue(tx.Statement.Context, "profileId", claims.ProfileId))
		// c.Set("db", tx)

		// Gắn claims vào context để các handler sau có thể sử dụng
		c.Set("claims", claims)
		c.Set("profileId", claims.ProfileId)
		c.Set("planId", claims.PlanId)

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

		c.Next() // Cho request đi tiếp

	})
}

func ExtractMetadataFromRequest(ctx context.Context, r *http.Request) metadata.MD {
	md := metadata.MD{}

	tokenString := r.Header.Get("Authorization")
	if tokenString != "" {
		// md.Append("authorization", tokenString)

		splitToken := strings.Split(tokenString, "Bearer ")
		var tokenStr string
		if len(splitToken) == 2 {
			tokenStr = strings.TrimSpace(splitToken[1])
		} else {
			tokenStr = strings.TrimSpace(splitToken[0])
		}
		claims, _ := ParseJWT(tokenStr)
		if claims != nil {
			md.Append("profileId", fmt.Sprintf("%d", claims.ProfileId))
			md.Append("planId", fmt.Sprintf("%d", *claims.PlanId))
		}

		// claims := jwt.MapClaims{}
		// token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte("secretKey"), nil
		// })
		// if err == nil && token.Valid {
		// 	if pid, ok := claims["profileId"].(float64); ok {
		// 		profileIdStr := fmt.Sprintf("%.0f", pid)
		// 		md.Append("profileid", profileIdStr)
		// 	}
		// }
	}

	md.Append("method", r.Method)
	log.Println("Metadata: ", md)
	return md
}

func ErrorToGatewayResponse(ctx context.Context, w http.ResponseWriter, err error) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	errorMap := make(map[string]string)

	for _, detail := range st.Details() {
		switch t := detail.(type) {
		case *pb_spec.FieldErrors:
			for k, v := range t.Errors {
				errorMap[k] = v
			}
		}
	}

	resp := map[string]interface{}{
		"code":    int(st.Code()),
		"message": st.Message(),
		"errors":  errorMap,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(runtime.HTTPStatusFromCode(st.Code()))
	_ = json.NewEncoder(w).Encode(resp)

	// var specErr *_routes.Except
	// if errors.As(err, &specErr) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK) // hoặc `w.WriteHeader(specErr.Code)` nếu bạn muốn
	// 	json.NewEncoder(w).Encode(specErr)
	// } else {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	_ = json.NewEncoder(w).Encode(map[string]interface{}{
	// 		"code":    500,
	// 		"message": err.Error(),
	// 	})
	// }
}
