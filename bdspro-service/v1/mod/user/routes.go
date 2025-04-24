package user

import (
	u_list "bdspro/v1/mod/user/list"
	u_product "bdspro/v1/mod/user/product"
	_jwt "common/jwt"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	product *u_product.UserProductRoute
	list    *u_list.ListRoute
}

func NewUserRoute(product *u_product.UserProductRoute,
	list *u_list.ListRoute,
) *UserRoute {
	return &UserRoute{
		product: product,
		list:    list,
	}
}

var publicRoutes = []string{"/bdspro/v1/user/list/property-type", "/bdspro/v1/user/product/market"}

func (r *UserRoute) RegisterUserRouter(e *gin.Engine) {
	userRouter := e.Group("/bdspro/v1/user")
	userRouter.Use(_jwt.JWTAuthMiddleware(publicRoutes))
	r.product.RegisterRouteUserProduct(userRouter)
	r.list.RegisterRouteListRoute(userRouter)
}
