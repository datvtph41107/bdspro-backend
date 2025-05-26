package user

import (
	u_action "bdspro/v1/mod/user/action"
	u_apartment "bdspro/v1/mod/user/apartment"
	u_asset "bdspro/v1/mod/user/asset"
	u_asset_cost "bdspro/v1/mod/user/asset_cost"
	u_asset_cost_type "bdspro/v1/mod/user/asset_cost_type"
	u_asset_exploitation "bdspro/v1/mod/user/asset_exploitation"
	u_asset_income_type "bdspro/v1/mod/user/asset_income_type"
	u_asset_legal "bdspro/v1/mod/user/asset_legal"
	u_asset_share "bdspro/v1/mod/user/asset_share"
	u_build_room "bdspro/v1/mod/user/build_room"
	u_list "bdspro/v1/mod/user/list"
	u_post "bdspro/v1/mod/user/post"
	u_product "bdspro/v1/mod/user/product"
	u_product_access "bdspro/v1/mod/user/product_access"
	u_product_child "bdspro/v1/mod/user/product_child"
	u_project_build "bdspro/v1/mod/user/project_build"
	_jwt "common/jwt"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	product           *u_product.UserProductRoute
	list              *u_list.ListRoute
	productChild      *u_product_child.UserProductChildRoute
	productAccess     *u_product_access.UserProductAccessRoute
	post              *u_post.PostRouter
	action            *u_action.ActionRouter
	asset             *u_asset.AssetRouter
	assetCost         *u_asset_cost.AssetCostRouter
	assetCostType     *u_asset_cost_type.AssetCostTypeRouter
	assetExploitation *u_asset_exploitation.AssetExploitationRouter
	assetLegal        *u_asset_legal.AssetLegalRouter
	assetIncomeType   *u_asset_income_type.AssetIncomeTypeRouter
	assetShare        *u_asset_share.AssetShareRouter
	buildRoom         *u_build_room.UserBuildRoomRouter
	apartment         *u_apartment.UserApartmentRouter
	projectBuild      *u_project_build.UserProjectBuildRouter
}

func NewUserRoute(product *u_product.UserProductRoute,
	list *u_list.ListRoute,
	productChild *u_product_child.UserProductChildRoute,
	productAccess *u_product_access.UserProductAccessRoute,
	post *u_post.PostRouter,
	action *u_action.ActionRouter,
	asset *u_asset.AssetRouter,
	assetCost *u_asset_cost.AssetCostRouter,
	assetCostType *u_asset_cost_type.AssetCostTypeRouter,
	assetExploitation *u_asset_exploitation.AssetExploitationRouter,
	assetLegal *u_asset_legal.AssetLegalRouter,
	assetIncomeType *u_asset_income_type.AssetIncomeTypeRouter,
	assetShare *u_asset_share.AssetShareRouter,
	buildRoom *u_build_room.UserBuildRoomRouter,
	apartment *u_apartment.UserApartmentRouter,
	projectBuild *u_project_build.UserProjectBuildRouter,
) *UserRoute {
	return &UserRoute{
		product:           product,
		list:              list,
		productChild:      productChild,
		productAccess:     productAccess,
		post:              post,
		action:            action,
		asset:             asset,
		assetCost:         assetCost,
		assetCostType:     assetCostType,
		assetExploitation: assetExploitation,
		assetLegal:        assetLegal,
		assetIncomeType:   assetIncomeType,
		assetShare:        assetShare,
		buildRoom:         buildRoom,
		apartment:         apartment,
		projectBuild:      projectBuild,
	}
}

var publicRoutes = []string{
	"/bdspro/v1/user/list/property-type",
	"/bdspro/v1/user/list/property-type",
	"/bdspro/v1/user/post/global",
	"/bdspro/v1/user/product/market",
}

func (r *UserRoute) RegisterUserRouter(e *gin.Engine) {
	userRouter := e.Group("/bdspro/v1/user")
	userRouter.Use(_jwt.JWTAuthMiddleware(publicRoutes))
	r.product.RegisterRouteUserProduct(userRouter)
	r.list.RegisterRouteListRoute(userRouter)
	r.productChild.RegisterRoute(userRouter)
	r.productAccess.RegisterRoutes(userRouter)
	r.post.RegisterRoutes(userRouter)
	r.action.RegisterRoutes(userRouter)
	r.asset.RegisterRoutes(userRouter)
	r.assetCost.RegisterRoutes(userRouter, "/asset-cost")
	r.assetCostType.RegisterRoutes(userRouter, "/type/asset-cost")
	r.assetIncomeType.RegisterRoutes(userRouter, "/type/asset-income")
	r.assetExploitation.RegisterRoutes(userRouter, "/asset-exploitation")
	r.assetLegal.RegisterRoutes(userRouter, "/asset-legal")
	r.assetShare.RegisterRoutes(userRouter, "/asset/share")
	r.apartment.RegisterRoutes(userRouter)
	r.buildRoom.RegisterRoutes(userRouter)
	r.projectBuild.RegisterRoutes(userRouter)
}
