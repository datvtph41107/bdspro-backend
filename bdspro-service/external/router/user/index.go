package user_router

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
	uAssetRoute         *UAssetRoute
	uAssetLegalRoute    *UAssetLegalRoute
	uPostRoute          *UserPostRoute
	uPostMediaRoute     *UPostMediaRoute
}

func NewRouterManager(uProductRouter *UProductRouter,
	uProductChildRouter *UProductChildRoute,
	uProductAccessRoute *UProductAccessRoute,
	uAssetRoute *UAssetRoute,
	uAssetLegalRoute *UAssetLegalRoute,
	uPostRoute *UserPostRoute,
	uPostMediaRoute *UPostMediaRoute,
) *RouterManager {
	return &RouterManager{
		uProductRouter:      uProductRouter,
		uProductChildRouter: uProductChildRouter,
		uProductAccessRoute: uProductAccessRoute,
		uAssetRoute:         uAssetRoute,
		uAssetLegalRoute:    uAssetLegalRoute,
		uPostRoute:          uPostRoute,
		uPostMediaRoute:     uPostMediaRoute,
	}
}

var publicRoutes = []string{
	"/bdspro/v2/user/list/property-type",
	"/bdspro/v2/user/list/property-type",
	"/bdspro/v2/user/post/global",
	"/bdspro/v2/user/product/market",
	"/bdspro/v2/user/posts/publish",
}

func (rm *RouterManager) RegisterRouter(e *gin.Engine) {
	userRouter := e.Group("/bdspro/v2/user")
	userRouter.Use(_jwt.JWTAuthMiddleware(publicRoutes))

	rm.uProductRouter.RegisterRoutes(userRouter)
	rm.uProductChildRouter.RegisterRoute(userRouter)
	rm.uProductAccessRoute.RegisterRoutes(userRouter)
	rm.uAssetRoute.RegisterRoutes(userRouter)
	rm.uAssetLegalRoute.RegisterRoutes(userRouter)
	rm.uPostRoute.RegisterRoutes(userRouter)
	rm.uPostMediaRoute.RegisterRoutes(userRouter)
}
