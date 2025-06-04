package user_router

import (
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostMediaRouter struct {
	UC *usecases.PostMediaUsecase
}

func NewPostMediaRoute(uc *usecases.PostMediaUsecase) *PostMediaRouter {
	return &PostMediaRouter{
		UC: uc,
	}
}

func (r *PostMediaRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
	group := router.Group(path)
	{
		group.GET("/:postId", r.Search)
	}
}

// @Summary Lấy danh sách media của bài đăng
// @Tags User: Tin đăng
// @Produce json
// @Param postId query string true "ID của bài đăng"
// @Security BearerAuth
// @Router /v1/bdspro/v2/user/post/media/{postId} [get]
func (s *PostMediaRouter) Search(c *gin.Context) {
	dto := dto.PostMediaGetDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	postId, err := strconv.ParseUint(c.Param("postId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := s.UC.GetMediaList(c, postId, dto)

	_routes.RouteResult(c, _routes.ResponseDTO{
		Code: 0,
		Data: result,
	}, err)
}
