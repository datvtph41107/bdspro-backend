package router

import (
	_routes "common/routes"
	_utils "common/utils"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type HistoryRouter struct {
	service *usecase.HistoryUsecase
}

func NewHistoryRouter(service *usecase.HistoryUsecase) *HistoryRouter {
	return &HistoryRouter{service: service}
}

func (r *HistoryRouter) HistoryRoutes(engine *gin.RouterGroup) {
	api := engine.Group("/histories")
	api.GET("/:id/:ownerType", r.Search)
}

// @Summary Lấy danh sách lịch sử
// @Description Lấy danh sách lịch sử
// @Tags Lịch sử
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param ownerType path enums.EOwnerType true "Owner Type"
// @Param dto query dto.HistorySearchDTO true "Thông tin tìm kiếm"
// @Security BearerAuth
// @Router /histories/{id}/{ownerType} [get]
func (r *HistoryRouter) Search(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	ownerTypeInt, err := strconv.ParseInt(c.Param("ownerType"), 10, 64)
	ownerType := enums.EOwnerType(ownerTypeInt)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	dto := dto.HistorySearchDTO{}
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, total, err := r.service.Search(c, id, ownerType, dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Message:       "Lấy danh sách lịch sử thành công",
		Data:          result,
		TotalElements: &total,
	}, err)
}
