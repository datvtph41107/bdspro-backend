package router

import (
	_routes "common/routes"
	_utils "common/utils"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/usecase"
	"strconv"

	"github.com/gin-gonic/gin"
)

type StageRouter struct {
	UC *usecase.StageUsecase
}

func NewStageRouter(service *usecase.StageUsecase) *StageRouter {
	return &StageRouter{UC: service}
}

func (r *StageRouter) StageRoutes(engine *gin.RouterGroup) {
	api := engine.Group("/stages")
	api.GET("/:pipelineId", r.Search)
	api.POST("", r.Create)
	api.PUT("/:id", r.Update)
	api.DELETE("/:id", r.Delete)
}

// @Summary Tìm kiếm giai đoạn theo pipeline
// @Description Tìm kiếm giai đoạn theo pipeline
// @Tags Giai đoạn
// @Accept json
// @Produce json
// @Param pipelineId path int true "Pipeline ID"
// @Param dto query dto.StageSearchDTO true "Stage search dto"
// @Security BearerAuth
// @Router /stages/{pipelineId} [get]
func (r *StageRouter) Search(c *gin.Context) {
	pipelineId, err := strconv.ParseUint(c.Param("pipelineId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	var dto dto.StageSearchDTO
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	stages, total, err := r.UC.Search(c, pipelineId, dto)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          stages,
		TotalElements: &total,
	}, nil)
}

// @Summary Tạo giai đoạn
// @Description Tạo giai đoạn
// @Tags Giai đoạn
// @Accept json
// @Produce json
// @Param entity body domain.StageEntity true "Stage entity"
// @Security BearerAuth
// @Router /stages [post]
func (r *StageRouter) Create(c *gin.Context) {
	entity := &domain.StageEntity{}
	if err := _utils.ParseBodyWithValidator(c, entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.Create(c, entity)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật giai đoạn
// @Description Cập nhật giai đoạn
// @Tags Giai đoạn
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param entity body domain.StageEntity true "Stage entity"
// @Security BearerAuth
// @Router /stages/{id} [put]
func (r *StageRouter) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	entity := &domain.StageEntity{}
	if err := _utils.ParseBodyWithValidator(c, entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.Update(c, id, entity)
	_routes.RouteResult(c, result, err)
}

// @Summary Xóa giai đoạn
// @Description Xóa giai đoạn
// @Tags Giai đoạn
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /stages/{id} [delete]
func (r *StageRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}
