package public_router

import (
	"github.com/gin-gonic/gin"
)

type PublicRouter struct {
}

func NewPublicRouter() *PublicRouter {
	return &PublicRouter{}
}

func (r *PublicRouter) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/bdspro/v2")
	{
		// group.GET("/product", r.GetProduct)
		group.GET("/enums/:type", r.GetEnumByType)
	}
}
