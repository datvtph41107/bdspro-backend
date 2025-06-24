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

type PipelineRouter struct {
	UC *usecase.PipelineUsecase
}

func NewPipelineRouter(UC *usecase.PipelineUsecase) *PipelineRouter {
	return &PipelineRouter{UC: UC}
}

func (r *PipelineRouter) PipelineRoutes(engine *gin.RouterGroup) {
	api := engine.Group("/pipelines")
	api.GET("", r.Search)
	api.POST("", r.Create)
	api.PUT("/:id", r.Update)
	api.DELETE("/:id", r.Delete)
}

func (r *PipelineRouter) Search(c *gin.Context) {
	var dto dto.PipelineSearchDTO
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	// ownerId, err := strconv.ParseUint(c.Param("ownerId"), 10, 64)
	// if err != nil {
	// 	_routes.RouteResult(c, nil, err)
	// 	return
	// }

	// ownerTypeValue, err := strconv.ParseUint(c.Param("ownerType"), 10, 64)
	// if err != nil {
	// 	_routes.RouteResult(c, nil, err)
	// 	return
	// }
	// ownerType := enums.EOwnerType(ownerTypeValue)

	pipelines, total, err := r.UC.Search(c, dto, false)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          pipelines,
		TotalElements: &total,
	}, nil)
}

func (r *PipelineRouter) Create(c *gin.Context) {
	entity := &domain.PipelineEntity{}
	if err := _utils.ParseBodyWithValidator(c, entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.Create(c, entity)
	_routes.RouteResult(c, result, err)
}

func (r *PipelineRouter) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	entity := &dto.PipelineSaveDTO{}
	if err := _utils.ParseBodyWithValidator(c, entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.Update(c, id, entity)
	_routes.RouteResult(c, result, err)
}

func (r *PipelineRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}
