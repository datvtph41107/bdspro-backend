package user

import (
	u_action "bdspro/v1/mod/user/action"
	u_list "bdspro/v1/mod/user/list"
	u_post "bdspro/v1/mod/user/post"
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
	post          *u_post.PostRouter
	action        *u_action.ActionRouter
}

func NewUserRoute(product *u_product.UserProductRoute,
	list *u_list.ListRoute,
	productChild *u_product_child.UserProductChildRoute,
	productAccess *u_product_access.UserProductAccessRoute,
	post *u_post.PostRouter,
	action *u_action.ActionRouter,
) *UserRoute {
	return &UserRoute{
		product:       product,
		list:          list,
		productChild:  productChild,
		productAccess: productAccess,
		post:          post,
		action:        action,
	}
}

var publicRoutes = []string{
	"/bdspro/user/list/property-type",
	"/bdspro/user/list/property-type",
	"/bdspro/user/post/global",
	"/bdspro/user/product/market",
}

func (r *UserRoute) RegisterUserRouter(e *gin.Engine) {
	userRouter := e.Group("/bdspro/user")
	userRouter.Use(_jwt.JWTAuthMiddleware(publicRoutes))
	r.product.RegisterRouteUserProduct(userRouter)
	r.list.RegisterRouteListRoute(userRouter)
	r.productChild.RegisterRoute(userRouter)
	r.productAccess.RegisterRoutes(userRouter)
	r.post.RegisterRoutes(userRouter)
	r.action.RegisterRoutes(userRouter)
}
