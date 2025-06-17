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

// @Summary Tạo đánh giá
// @Description Tạo đánh giá mới
// @Tags Rate
// @Accept json
// @Produce json
// @Param rate body user_model.Rate true "Thông tin đánh giá"
// @Security BearerAuth
// @Router /v1/user/rate [post]
func (r *RateRouter) Create(c *gin.Context) {
	var dto user_model.Rate
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.RateService.Create(c, &dto)
	_routes.RouteResult(c, result, err)
}

// @Summary Lấy danh sách đánh giá
// @Description Lấy danh sách đánh giá theo targetId và targetType
// @Tags Rate
// @Accept json
// @Produce json
// @Param targetId path string true "ID của target"
// @Param targetType path string true "Loại của target, 10: user, 20: group, 30: organization"
// @Param page query int false "Trang hiện tại"
// @Param size query int false "Số lượng trên mỗi trang"
// @Param score query int false "Điểm đánh giá, 1: 1 sao, 2: 2 sao, 3: 3 sao, 4: 4 sao, 5: 5 sao"
// @Param parentId query string false "ID của đánh giá cha"
// @Router /v1/user/rate/list/{targetId}/{targetType} [get]
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

// @Summary Cập nhật đánh giá
// @Description Cập nhật đánh giá theo id
// @Tags Rate
// @Accept json
// @Produce json
// @Param id path string true "ID của đánh giá"
// @Param rate body user_model.Rate true "Thông tin đánh giá"
// @Security BearerAuth
// @Router /v1/user/rate/{id} [put]
func (r *RateRouter) Update(c *gin.Context) {
	var dto user_model.Rate
	if err := _utils.ParseBodyWithValidator(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.RateService.Update(c, dto.ID, &dto)
	_routes.RouteResult(c, result, err)
}

// @Summary Xóa đánh giá
// @Description Xóa đánh giá theo id
// @Tags Rate
// @Accept json
// @Produce json
// @Param id path string true "ID của đánh giá"
// @Security BearerAuth
// @Router /v1/user/rate/{id} [delete]
func (r *RateRouter) Delete(c *gin.Context) {

}
