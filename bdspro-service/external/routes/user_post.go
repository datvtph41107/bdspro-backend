package routes

import (
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserPostRoute struct {
	UC *usecases.UPostUsecase
}

func NewUserPostRoute(uc *usecases.UPostUsecase) *UserPostRoute {
	return &UserPostRoute{
		UC: uc,
	}
}

func (r *UserPostRoute) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/posts")
	{
		// todo:
		group.POST("", r.CreatePost)
		group.POST("/new-expired", r.NewExpiredPost)
		group.GET("/:id", r.GetDetail)
		group.PUT("/:id/edit", r.UpdatePost)
		group.DELETE("/:id", r.Delete)
		group.PUT("/hidden", r.UpdateHidden)
		group.GET("/personal", r.GetPersonalPost)
		group.GET("/global", r.GetGlobalPost)
		// group.PUT("/action", r.Action)
		// group.GET("/share", r.Share)
		//
		// group.GET("/dashboard", r.GetGlobalPost) // thống kê lượt xem, click liên hệ
		group.GET("/up-vip", r.GetGlobalPost)
		group.GET("/marketing", r.GetGlobalPost) // quảng cáo
		group.GET("/export", r.GetGlobalPost)    // xuất excel
	}
}

func (r *UserPostRoute) CreatePost(c *gin.Context) {
	body := dto.PostSaveRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.CreatePost(c, &body)
	_routes.RouteResult(c, result, err)
}

func (r *UserPostRoute) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	body := dto.PostUpdateRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.UpdatePost(c, id, body)
	_routes.RouteResult(c, result, err)
}

func (r *UserPostRoute) UpdateHidden(c *gin.Context) {
	body := dto.UpdateHiddenRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.UpdateHidden(c, body)
	_routes.RouteResult(c, result, err)
}

// func (r *UserPostRoute) Action(c *gin.Context) {
// 	body := ActionRequest{}
// 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, err := r.UC.Action(c, body)
// 	_routes.RouteResult(c, result, err)
// }

func (r *UserPostRoute) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.Detail(c, id)
	_routes.RouteResult(c, result, err)
}

func (r *UserPostRoute) GetPersonalPost(c *gin.Context) {
	query := dto.PostSearchRequest{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.PersonalPost(c, query)
	_routes.RouteResult(c, result, err)
}

func (r *UserPostRoute) GetGlobalPost(c *gin.Context) {
	query := dto.PostSearchRequest{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.GlobalPost(c, query)
	_routes.RouteResult(c, *result, err)
}

// todo: logic thanh toán chưa rõ ràng
func (r *UserPostRoute) NewExpiredPost(c *gin.Context) {
	body := dto.PostExpiredRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.MakeNewExpired(c, body)
	_routes.RouteResult(c, result, err)
}

// todo: logic thanh toán chưa rõ ràng
func (r *UserPostRoute) Share(c *gin.Context) {
	body := dto.PostExpiredRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.MakeNewExpired(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *UserPostRoute) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}
