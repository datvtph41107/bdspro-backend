package _jwt

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Middleware để xác thực JWT (Gin). Prefer JWTAuthHTTPMiddleware for net/http servers.
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

		req, status, message, silent := authenticateRequest(c.Request, pubRoutes, tempRoutes, cfg)
		if silent {
			c.Abort()
			return
		}
		if status != 0 {
			c.AbortWithStatusJSON(status, gin.H{
				"code":    status,
				"message": message,
			})
			return
		}

		c.Request = req
		if claims, ok := PrincipalFromRequestContext(req.Context()); ok && claims != nil {
			c.Set("claims", claims)
			c.Set("organizationId", claims.OrganizationId)
			c.Set("profileId", claims.ProfileId)
			c.Set("planId", claims.PlanId)
		}
		c.Next()
	}
}
