package u_post

import (
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
	postService *PostService
	crud2.BaseRouter[models.Post, *PostService]
}

func NewPostRouter(service *PostService) *PostRouter {
	return &PostRouter{
		postService: service,
		BaseRouter: crud2.BaseRouter[models.Post, *PostService]{
			Service: service,
		},
	}
}

func (r *PostRouter) RegisterRoutes(router *gin.RouterGroup) {
	group := router.Group("/post")
	{
		// todo:
		group.POST("", r.CreatePost)
		group.POST("/new-expired", r.NewExpiredPost)
		group.GET("/:id", r.GetDetail)
		group.PUT("/:id", r.UpdatePost)
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
	body := PostSaveRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.postService.CreatePost(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *PostRouter) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	body := PostUpdateRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.postService.UpdatePost(c, id, body)
	_routes.RouteResult(c, result, err)
}

func (r *PostRouter) UpdateHidden(c *gin.Context) {
	body := UpdateHiddenRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.postService.UpdateHidden(c, body)
	_routes.RouteResult(c, result, err)
}

// func (r *PostRouter) Action(c *gin.Context) {
// 	body := ActionRequest{}
// 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, err := r.postService.Action(c, body)
// 	_routes.RouteResult(c, result, err)
// }

func (r *PostRouter) GetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.postService.Detail(c, id)
	_routes.RouteResult(c, result, err)
}

func (r *PostRouter) GetPersonalPost(c *gin.Context) {
	query := SearchRequest{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.postService.PersonalPost(c, query)
	_routes.RouteResult(c, *result, err)
}

func (r *PostRouter) GetGlobalPost(c *gin.Context) {
	query := SearchRequest{}
	if err := _utils.ParseQuery2(c, &query); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.postService.GlobalPost(c, query)
	_routes.RouteResult(c, *result, err)
}

// todo: logic thanh toán chưa rõ ràng
func (r *PostRouter) NewExpiredPost(c *gin.Context) {
	body := PostExpiredRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.postService.MakeNewExpired(c, body)
	_routes.RouteResult(c, result, err)
}

// todo: logic thanh toán chưa rõ ràng
func (r *PostRouter) Share(c *gin.Context) {
	body := PostExpiredRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.postService.MakeNewExpired(c, body)
	_routes.RouteResult(c, result, err)
}
