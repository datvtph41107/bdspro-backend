package group_router

import "github.com/gin-gonic/gin"

type GroupManager struct {
	groupProductRouter *GroupProductRouter
}

func NewGroupRouter(groupProductRouter *GroupProductRouter) *GroupManager {
	return &GroupManager{
		groupProductRouter: groupProductRouter,
	}
}

func (rm *GroupManager) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/bdspro/v2/group")
	rm.groupProductRouter.RegisterRoutes(group)
}
