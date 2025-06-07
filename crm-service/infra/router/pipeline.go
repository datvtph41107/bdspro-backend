package router

import (
	_routes "common/routes"
	_utils "common/utils"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
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
	api.GET("/:ownerId/:ownerType", r.Search)
	api.POST("", r.Create)
	api.PUT("/:id", r.Update)
	api.DELETE("/:id", r.Delete)
}

// @Summary Tìm kiếm quy trình
// @Description Tìm kiếm quy trình
// @Tags Quy trình
// @Accept json
// @Produce json
// @Param ownerId path int true "Owner ID"
// @Param ownerType path enums.EOwnerType true "Owner Type"
// @Param dto query dto.PipelineSearchDTO true "Pipeline search dto"
// @Security BearerAuth
// @Router /pipeline/{ownerId}/{ownerType} [get]
func (r *PipelineRouter) Search(c *gin.Context) {
	var dto dto.PipelineSearchDTO
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	ownerId, err := strconv.ParseUint(c.Param("ownerId"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	ownerTypeValue, err := strconv.ParseUint(c.Param("ownerType"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	ownerType := enums.EOwnerType(ownerTypeValue)

	pipelines, total, err := r.UC.Search(c, ownerId, ownerType, dto)
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

// @Summary Tạo quy trình
// @Description Tạo quy trình
// @Tags Quy trình
// @Accept json
// @Produce json
// @Param entity body domain.PipelineEntity true "Pipeline entity"
// @Security BearerAuth
// @Router /pipeline [post]
func (r *PipelineRouter) Create(c *gin.Context) {
	entity := &domain.PipelineEntity{}
	if err := _utils.ParseBodyWithValidator(c, entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.Create(c, 0, entity)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật quy trình
// @Description Cập nhật quy trình
// @Tags Quy trình
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param entity body domain.PipelineEntity true "Pipeline entity"
// @Security BearerAuth
// @Router /pipeline/{id} [put]
func (r *PipelineRouter) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	entity := &domain.PipelineEntity{}
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
