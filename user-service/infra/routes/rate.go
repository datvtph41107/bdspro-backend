package routes

import (
	_routes "common/routes"
	_utils "common/utils"
	"strconv"
	"user/internal/dto"
	user_model "user/internal/models/user"
	user_service "user/internal/services/user"

	"github.com/gin-gonic/gin"
)

type RateRouter struct {
	RateService *user_service.RateService
}

func NewRateRouter(rateService *user_service.RateService) *RateRouter {
	return &RateRouter{RateService: rateService}
}

func RateRoutes(r *gin.Engine, router *RateRouter) {
	g := r.Group("/user/rate")
	{
		g.POST("", router.Create)
		g.GET("/list/:targetId/:targetType", router.List)
		g.PUT("/:id", router.Update)
		g.DELETE("/:id", router.Delete)
	}

}

func (r *RateRouter) Create(c *gin.Context) {
	var dto user_model.Rate
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.RateService.Create(c, &dto)
	_routes.RouteResult(c, result, err)
}

func (r *RateRouter) List(c *gin.Context) {
	targetID, _ := strconv.ParseUint(c.Param("targetId"), 10, 64)
	targetType, _ := strconv.ParseUint(c.Param("targetType"), 10, 8)

	var dto dto.RateSearch
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.RateService.List(c, targetID, uint8(targetType), dto)
	_routes.RouteResult(c, result, err)
}

func (r *RateRouter) Update(c *gin.Context) {
	var dto user_model.Rate
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.RateService.Update(c, dto.ID, &dto)
	_routes.RouteResult(c, result, err)
}

func (r *RateRouter) Delete(c *gin.Context) {

}
