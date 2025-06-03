package org_router

import (
	"github.com/gin-gonic/gin"
)

type RouterManager struct {
	orgAccessRoute   *OrgAccessRoute
	orgProductRouter *OrgProductRouter
}

func NewOrgRouterManager(orgAccessRoute *OrgAccessRoute, orgProductRouter *OrgProductRouter) *RouterManager {
	return &RouterManager{
		orgAccessRoute:   orgAccessRoute,
		orgProductRouter: orgProductRouter,
	}
}

func (rm *RouterManager) RegisterRoutes(router *gin.Engine) {
	orgGroup := router.Group("/bdspro/v2/org")
	rm.orgAccessRoute.RegisterRoutes(orgGroup)
	rm.orgProductRouter.RegisterRoutes(orgGroup)
}
