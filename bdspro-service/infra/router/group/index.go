package group_router

import "github.com/gin-gonic/gin"

type GroupManager struct {
	groupProductRouter *GroupProductRouter
	groupAssetRouter   *GroupAssetRouter
}

func NewGroupRouter(
	groupProductRouter *GroupProductRouter,
	groupAssetRouter *GroupAssetRouter,
) *GroupManager {
	return &GroupManager{
		groupProductRouter: groupProductRouter,
		groupAssetRouter:   groupAssetRouter,
	}
}

func (rm *GroupManager) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/bdspro/v2/group")
	rm.groupProductRouter.RegisterRoutes(group)
	rm.groupAssetRouter.RegisterRoutes(group)
}
