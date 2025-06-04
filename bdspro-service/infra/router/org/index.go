package org_router

import (
	"github.com/gin-gonic/gin"
)

type RouterManager struct {
	// orgAccessRoute *OrgAccessRoute
	// orgProductRouter *OrgProductRouter
	orgAssetRouter             *OrgAssetRouter
	orgCostTypeRouter          *OrgCostTypeRouter
	orgAssetCostRouter         *OrgAssetCostRouter
	orgAssetExploitationRouter *OrgAssetExploitationRouter
	orgProductRouter           *OrgProductRouter
	orgPostRouter              *OrgPostRouter
}

func NewOrgRouterManager(
	// orgAccessRoute *OrgAccessRoute,
	// orgProductRouter *OrgProductRouter,
	orgAssetRouter *OrgAssetRouter,
	orgCostTypeRouter *OrgCostTypeRouter,
	orgAssetCostRouter *OrgAssetCostRouter,
	orgAssetExploitationRouter *OrgAssetExploitationRouter,
	orgProductRouter *OrgProductRouter,
	orgPostRouter *OrgPostRouter,
) *RouterManager {
	return &RouterManager{
		// orgAccessRoute: orgAccessRoute,
		// orgProductRouter: orgProductRouter,
		orgAssetRouter:             orgAssetRouter,
		orgCostTypeRouter:          orgCostTypeRouter,
		orgAssetCostRouter:         orgAssetCostRouter,
		orgAssetExploitationRouter: orgAssetExploitationRouter,
		orgProductRouter:           orgProductRouter,
		orgPostRouter:              orgPostRouter,
	}
}

func (rm *RouterManager) RegisterRoutes(e *gin.Engine) {
	orgGroup := e.Group("/bdspro/v2/org")
	// rm.orgAccessRoute.RegisterRoutes(orgGroup)
	// rm.orgProductRouter.RegisterRoutes(orgGroup)
	rm.orgAssetRouter.RegisterRoutes(orgGroup, "/asset")
	rm.orgCostTypeRouter.RegisterRoutes(orgGroup, "/assets/cost-type")
	rm.orgAssetCostRouter.RegisterRoutes(orgGroup, "/assets/cost")
	rm.orgAssetExploitationRouter.RegisterRoutes(orgGroup, "/assets/exploitation")
	rm.orgProductRouter.RegisterRoutes(orgGroup, "/product")
	rm.orgPostRouter.RegisterRoutes(orgGroup, "/post")
}
