package u_post

import (
	u_post_dto "bdspro/v1/mod/user/post/dto"
	"bdspro/v1/models"
	"common/case/crud2"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// import (
// 	"bdspro/v1/models"
// 	"common/case/crud"

// 	"github.com/gin-gonic/gin"
// )

type PostRouter struct {
	UserPostService *UserPostService
	crud2.BaseRouter[models.Post, *UserPostService]
}

func NewPostRouter(service *UserPostService) *PostRouter {
	return &PostRouter{
		UserPostService: service,
		BaseRouter: crud2.BaseRouter[models.Post, *UserPostService]{
			Service: service,
		},
	}
}

func (r *PostRouter) RegisterRoutes(router *gin.RouterGroup) {
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

func (r *PostRouter) CreatePost(c *gin.Context) {
	body := u_post_dto.PostSaveRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UserPostService.CreatePost(c, &body)
	_routes.RouteResult(c, result, err)
}

func (r *PostRouter) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	body := u_post_dto.PostUpdateRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UserPostService.UpdatePost(c, id, body)
	_routes.RouteResult(c, result, err)
}

func (r *PostRouter) UpdateHidden(c *gin.Context) {
	body := u_post_dto.UpdateHiddenRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UserPostService.UpdateHidden(c, body)
	_routes.RouteResult(c, result, err)
}

// func (r *PostRouter) Action(c *gin.Context) {
// 	body := ActionRequest{}
// 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, err := r.UserPostService.Action(c, body)
// 	_routes.RouteResult(c, result, err)
// }

func (r *PostRouter) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UserPostService.Detail(c, id)
	_routes.RouteResult(c, result, err)
}

func (r *PostRouter) GetPersonalPost(c *gin.Context) {
	query := u_post_dto.SearchRequest{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UserPostService.PersonalPost(c, query)
	_routes.RouteResult(c, result, err)
}

func (r *PostRouter) GetGlobalPost(c *gin.Context) {
	query := u_post_dto.SearchRequest{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UserPostService.GlobalPost(c, query)
	_routes.RouteResult(c, *result, err)
}

// todo: logic thanh toán chưa rõ ràng
func (r *PostRouter) NewExpiredPost(c *gin.Context) {
	body := u_post_dto.PostExpiredRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UserPostService.MakeNewExpired(c, body)
	_routes.RouteResult(c, result, err)
}

// todo: logic thanh toán chưa rõ ràng
func (r *PostRouter) Share(c *gin.Context) {
	body := u_post_dto.PostExpiredRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UserPostService.MakeNewExpired(c, body)
	_routes.RouteResult(c, result, err)
}
