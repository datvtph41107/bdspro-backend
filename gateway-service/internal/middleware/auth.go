package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware checks JWT for non-public routes.
func JWTAuthMiddleware(publicRoutes, tempRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		// Check if path is public
		for _, route := range publicRoutes {
			if strings.HasPrefix(path, route) {
				c.Next()
				return
			}
		}
		for _, route := range tempRoutes {
			if strings.HasPrefix(path, route) {
				c.Next()
				return
			}
		}
		// TODO: actual JWT validation using cfg.Auth.JWTSecret
		// For now, placeholder
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		// Validate token...
		c.Next()
	}
}
