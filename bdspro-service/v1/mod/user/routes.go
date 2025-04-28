package user

import (
	u_list "bdspro/v1/mod/user/list"
	u_product "bdspro/v1/mod/user/product"
	u_product_access "bdspro/v1/mod/user/product_access"
	u_product_child "bdspro/v1/mod/user/product_child"
	_jwt "common/jwt"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	product       *u_product.UserProductRoute
	list          *u_list.ListRoute
	productChild  *u_product_child.UserProductChildRoute
	productAccess *u_product_access.UserProductAccessRoute
}

func NewUserRoute(product *u_product.UserProductRoute,
	list *u_list.ListRoute,
	productChild *u_product_child.UserProductChildRoute,
	productAccess *u_product_access.UserProductAccessRoute,
) *UserRoute {
	return &UserRoute{
		product:       product,
		list:          list,
		productChild:  productChild,
		productAccess: productAccess,
	}
}

var publicRoutes = []string{"/bdspro/v1/user/list/property-type", "/bdspro/v1/user/product/market"}

func (r *UserRoute) RegisterUserRouter(e *gin.Engine) {
	userRouter := e.Group("/bdspro/v1/user")
	userRouter.Use(_jwt.JWTAuthMiddleware(publicRoutes))
	r.product.RegisterRouteUserProduct(userRouter)
	r.list.RegisterRouteListRoute(userRouter)
	r.productChild.RegisterRoute(userRouter)
	r.productAccess.RegisterRoutes(userRouter)
}
