package user_router

import "github.com/gin-gonic/gin"

//	 userProductService.UserPostService = userPostService
//		userProductService.UserAssetService = userAssetService
type RouterManager struct {
	// uProductRouter          *UProductRouter
	// uProductChildRouter     *UProductChildRoute
	// uProductAccessRoute     *UProductAccessRoute
	// uAssetLegalRoute        *UAssetLegalRoute
	// uPostRoute              *UserPostRoute
	// uPostMediaRoute         *UPostMediaRoute
	// uProjectBuildRoute      *UserProjectBuildRouter
	userProductChildRouter   *UserProductChildRouter
	postMediaRouter          *PostMediaRouter
	uAssetRouter             *UserAssetRouter
	uAssetCostTypeRouter     *UserCostTypeRouter
	uAssetCostRouter         *UserAssetCostRouter
	uAssetExploitationRouter *UserAssetExploitationRouter
	uListRouter              *ListRouter
	uProductRouter           *UserProductRouter
	uPostRouter              *UserPostRouter
	uAssetLegalRouter        *UserAssetLegalRouter
	uAssetShareRouter        *AssetShareRouter
	uProductAccessRouter     *UserProductAccessRouter
}

func NewRouterManager(
	// uProductRouter *UProductRouter,
	// uProductChildRouter *UProductChildRoute,
	// uProductAccessRoute *UProductAccessRoute,
	// uAssetLegalRoute *UAssetLegalRoute,
	// uPostRoute *UserPostRoute,
	// uPostMediaRoute *UPostMediaRoute,
	// uProjectBuildRoute *UserProjectBuildRouter,
	// uAssetExploitationRoute *UAssetExploitationRoute,
	uListRouter *ListRouter,
	uAssetRouter *UserAssetRouter,
	uAssetCostTypeRouter *UserCostTypeRouter,
	uAssetCostRouter *UserAssetCostRouter,
	uAssetExploitationRouter *UserAssetExploitationRouter,
	uProductRouter *UserProductRouter,
	uPostRouter *UserPostRouter,
	postMediaRouter *PostMediaRouter,
	userProductChildRouter *UserProductChildRouter,
	uAssetLegalRouter *UserAssetLegalRouter,
	uAssetShareRouter *AssetShareRouter,
	uProductAccessRouter *UserProductAccessRouter,
) *RouterManager {
	return &RouterManager{
		// uProductRouter:          uProductRouter,
		// uProductChildRouter:     uProductChildRouter,
		// uProductAccessRoute:     uProductAccessRoute,
		// uAssetLegalRoute:        uAssetLegalRoute,
		// uPostRoute:              uPostRoute,
		// uPostMediaRoute:         uPostMediaRoute,
		// uProjectBuildRoute:      uProjectBuildRoute,
		// uAssetExploitationRoute: uAssetExploitationRoute,
		uAssetLegalRouter:        uAssetLegalRouter,
		userProductChildRouter:   userProductChildRouter,
		postMediaRouter:          postMediaRouter,
		uAssetRouter:             uAssetRouter,
		uListRouter:              uListRouter,
		uAssetCostTypeRouter:     uAssetCostTypeRouter,
		uAssetCostRouter:         uAssetCostRouter,
		uAssetExploitationRouter: uAssetExploitationRouter,
		uProductRouter:           uProductRouter,
		uPostRouter:              uPostRouter,
		uAssetShareRouter:        uAssetShareRouter,
		uProductAccessRouter:     uProductAccessRouter,
	}
}

// var publicRoutes = []string{
// 	"/bdspro/v2/user/list/property-type",
// 	"/bdspro/v2/user/list/property-type",
// 	"/bdspro/v2/user/post/global",
// 	"/bdspro/v2/user/product/market",
// 	"/bdspro/v2/user/posts/publish",
// }

func (rm *RouterManager) RegisterRouter(e *gin.Engine) {
	userRouter := e.Group("/bdspro/v2/user")
	// userRouter.Use(_jwt.JWTAuthMiddleware(publicRoutes))

	rm.uListRouter.RegisterRouter(userRouter)
	// rm.uAssetRouter.RegisterRoutes(userRouter, "/asset")
	rm.uAssetCostTypeRouter.RegisterRoutes(userRouter, "/assets/cost-type")
	rm.uAssetCostRouter.RegisterRoutes(userRouter, "/assets/cost")
	rm.uAssetExploitationRouter.RegisterRoutes(userRouter, "/assets/exploitation")
	rm.uProductRouter.RegisterRoutes(userRouter, "/product")
	rm.uPostRouter.RegisterRoutes(userRouter, "/post")
	rm.postMediaRouter.RegisterRoutes(userRouter, "/post/media")
	rm.userProductChildRouter.RegisterRoutes(userRouter, "/product/child")
	rm.uAssetLegalRouter.RegisterRoutes(userRouter, "/asset/legal")
	rm.uAssetShareRouter.RegisterRoutes(userRouter, "/asset/access")
	rm.uProductAccessRouter.RegisterRoutes(userRouter, "/product/access")
	// rm.uProductRouter.RegisterRoutes(userRouter)
	// rm.uProductChildRouter.RegisterRoute(userRouter)
	// rm.uProductAccessRoute.RegisterRoutes(userRouter)
	// rm.uAssetRoute.RegisterRoutes(userRouter)
	// rm.uAssetLegalRoute.RegisterRoutes(userRouter)
	// rm.uPostRoute.RegisterRoutes(userRouter)
	// rm.uPostMediaRoute.RegisterRoutes(userRouter)
	// rm.uProjectBuildRoute.RegisterRoutes(userRouter)
	// rm.uAssetExploitationRoute.RegisterRoutes(userRouter, "/asset-exploitation")
}
