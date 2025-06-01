package user_router

import (
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"

	"github.com/gin-gonic/gin"
)

type UPostMediaRoute struct {
	UC *usecases.UPostMediaUsecase
}

func NewUPostMediaRoute(uc *usecases.UPostMediaUsecase) *UPostMediaRoute {
	return &UPostMediaRoute{
		UC: uc,
	}
}

func (r *UPostMediaRoute) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/post-media")
	{
		group.GET("", r.Search)
	}
}

func (s *UPostMediaRoute) Search(c *gin.Context) {
	dto := dto.PostMediaGetDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := s.UC.GetMediaList(c, dto)

	_routes.RouteResult(c, _routes.ResponseDTO{
		Code: 0,
		Data: result,
	}, err)
}
