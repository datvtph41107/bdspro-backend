package router

import (
	_routes "common/routes"
	_utils "common/utils"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/usecase"
	"crm/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RuleRouter struct {
	UC *usecase.RuleUsecase
}

func NewRuleRouter(UC *usecase.RuleUsecase) *RuleRouter {
	return &RuleRouter{UC: UC}
}

func (r *RuleRouter) RuleRoutes(engine *gin.RouterGroup) {
	api := engine.Group("/rule")
	api.GET("/:ownerId/:ownerType", r.Search)
	api.POST("", r.Create)
	api.PUT("/:id", r.Update)
	api.DELETE("/:id", r.Delete)
}

// @Summary Tìm kiếm quy tắc
// @Description Tìm kiếm quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param ownerId path int true "Owner ID"
// @Param ownerType path enums.EOwnerType true "Owner Type"
// @Param dto query dto.RuleSearchDTO true "Rule search dto"
// @Security BearerAuth
// @Router /rule/{ownerId}/{ownerType} [get]
func (r *RuleRouter) Search(c *gin.Context) {
	var dto dto.RuleSearchDTO
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	ownerId, ownerType, err := utils.GetOwnerFromContext(c)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	rules, total, err := r.UC.Search(c, ownerId, ownerType, dto)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          rules,
		TotalElements: &total,
	}, nil)
}

// @Summary Tạo quy tắc
// @Description Tạo quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param entity body domain.RuleEntity true "Rule entity"
// @Security BearerAuth
// @Router /rule [post]
func (r *RuleRouter) Create(c *gin.Context) {
	entity := &domain.RuleEntity{}
	if err := _utils.ParseBodyWithValidator(c, entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.Create(c, entity)
	_routes.RouteResult(c, result, err)
}

// @Summary Cập nhật quy tắc
// @Description Cập nhật quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param entity body domain.RuleEntity true "Rule entity"
// @Security BearerAuth
// @Router /rule/{id} [put]
func (r *RuleRouter) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	entity := &domain.RuleEntity{}
	if err := _utils.ParseBodyWithValidator(c, entity); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.Update(c, id, entity)
	_routes.RouteResult(c, result, err)
}

// @Summary Xóa quy tắc
// @Description Xóa quy tắc
// @Tags Quy tắc
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /rule/{id} [delete]
func (r *RuleRouter) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}
