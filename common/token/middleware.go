package token

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		if isPublicRoute(pubRoutes, c.Request.URL.Path) && tokenString == "" {
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

		c.Next()
	}
}
