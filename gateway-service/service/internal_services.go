package service

import (
	"os"

	"common/tcp/connect"

	"github.com/gin-gonic/gin"
)

// RegisterInternalServicesRoute đăng ký GET /internal/services —
// các microservice gọi khi start để biết upstream nào đang up rồi dial lại.
func RegisterInternalServicesRoute(r *gin.Engine) {
	r.GET("/internal/services", func(c *gin.Context) {
		env := os.Getenv("ENV_RUNTIME")
		c.JSON(200, connect.ListFromViper(env))
	})
}
