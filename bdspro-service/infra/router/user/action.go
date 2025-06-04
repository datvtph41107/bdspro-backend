package user_router

// type ActionRouter struct {
// 	// postService *u_post.PostService
// 	crud2.BaseRouter[domain.Action, *ActionService]
// }

// func NewActionRouter(service *ActionService) *ActionRouter {
// 	return &ActionRouter{
// 		// postService: service,
// 		BaseRouter: crud2.BaseRouter[domain.Action, *ActionService]{
// 			Service: service,
// 		},
// 	}
// }

// func (r *ActionRouter) RegisterRoutes(router *gin.RouterGroup) {
// 	group := router.Group("/action")
// 	{
// 		// todo:
// 		group.POST("", r.CreateAction)
// 		group.GET("/dashboard/post/:id", r.DashboardPost)
// 		group.GET("/dashboard/post/bar", r.DashboardPostBar)
// 	}
// }

// func (r *ActionRouter) CreateAction(c *gin.Context) {
// 	body := ActionRequest{}
// 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, err := r.Service.CreateAction(c, body)
// 	_routes.RouteResult(c, result, err)
// }

// func (r *ActionRouter) DashboardPost(c *gin.Context) {
// 	// body := DashboardRequest{}
// 	// if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 	// 	_routes.RouteResult(c, nil, err)
// 	// 	return
// 	// }

// 	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
// 	if err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}
// 	result, err := r.Service.DashboardPost(c, id)
// 	_routes.RouteResult(c, result, err)
// }

// func (r *ActionRouter) DashboardPostBar(c *gin.Context) {
// 	body := PostBarRequest{}
// 	if err := _utils.ParseQuery2(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}
// 	result, err := r.Service.DashboardPostBar(c, body)
// 	_routes.RouteResult(c, result, err)
// }
