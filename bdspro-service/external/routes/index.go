package routes

import (
	_jwt "common/jwt"

	"github.com/gin-gonic/gin"
)

//	 userProductService.UserPostService = userPostService
//		userProductService.UserAssetService = userAssetService
type RouterManager struct {
	uProductRouter      *UProductRouter
	uProductChildRouter *UProductChildRoute
	uProductAccessRoute *UProductAccessRoute
}

func NewRouterManager(uProductRouter *UProductRouter,
	uProductChildRouter *UProductChildRoute,
	uProductAccessRoute *UProductAccessRoute,
) *RouterManager {
	return &RouterManager{
		uProductRouter:      uProductRouter,
		uProductChildRouter: uProductChildRouter,
		uProductAccessRoute: uProductAccessRoute,
	}
}

var publicRoutes = []string{
	"/bdspro/v2/user/list/property-type",
	"/bdspro/v2/user/list/property-type",
	"/bdspro/v2/user/post/global",
	"/bdspro/v2/user/product/market",
}

func (rm *RouterManager) RegisterRouter(e *gin.Engine) {
	userRouter := e.Group("/bdspro/v2/user")
	userRouter.Use(_jwt.JWTAuthMiddleware(publicRoutes))

	rm.uProductRouter.RegisterRoutes(userRouter)
	rm.uProductChildRouter.RegisterRoute(userRouter)
	rm.uProductAccessRoute.RegisterRoutes(userRouter)
}
