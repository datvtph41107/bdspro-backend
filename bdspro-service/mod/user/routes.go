package user

import (
	u_product "bdspro/mod/user/product"
	_jwt "common/jwt"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	product *u_product.UserProductRoute
}

func NewUserRoute(product *u_product.UserProductRoute) *UserRoute {
	return &UserRoute{
		product: product,
	}
}

func (r *UserRoute) RegisterUserRouter(e *gin.Engine) {
	userRouter := e.Group("/bdspro/user")
	userRouter.Use(_jwt.JWTAuthMiddleware(nil))
	r.product.RegisterRouteUserProduct(userRouter)
}
