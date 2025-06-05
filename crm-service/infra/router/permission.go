package router

import (
	_routes "common/routes"
	"crm/infra/impl"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionRouter struct {
	UC *impl.PermissionUsecase
}

func NewPermissionRouter(UC *impl.PermissionUsecase) *PermissionRouter {
	return &PermissionRouter{
		UC: UC,
	}
}

func (r *PermissionRouter) PermissionRoutes(engine *gin.RouterGroup) {
	group := engine.Group("/permission")
	group.POST("/clear/:profileId", r.ClearProfileID)
}

func (r *PermissionRouter) ClearProfileID(c *gin.Context) {
	profileId, err := strconv.ParseUint(c.Param("profileId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = r.UC.ClearProfileID(c, profileId)
	_routes.RouteResult(c, &_routes.ResponseDTO{
		Code:    0,
		Message: "Reset thành công",
	}, err)
}
