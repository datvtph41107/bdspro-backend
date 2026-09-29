package shared

import (
	"bdspro/internal/common"
	"bdspro/internal/enums"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ICrudOwnerRouter[T any, D common.DTO] interface {
	RegisterRoutes(router *gin.RouterGroup, path string)
	Search(c *gin.Context)
	GetByID(c *gin.Context)
	Detail(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type CrudOwnerRouter[T any, D common.DTO] struct {
	UC common.IOwnerUsecase[T, D]
}

func (r *CrudOwnerRouter[T, D]) RegisterRoutes(router *gin.RouterGroup, path string) {
	group := router.Group(path)
	{
		group.GET("", r.Search)
		group.POST("", r.Create)
		group.PUT(":id", r.Update)
		group.DELETE(":id", r.Delete)
		group.GET(":id", r.GetByID)
		group.GET(":id/detail", r.Detail)
	}
}

func (r *CrudOwnerRouter[T, D]) Search(c *gin.Context) {
	var dto D
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := r.UC.Search(c, enums.EOwnerOfMember, dto)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

func (r *CrudOwnerRouter[T, D]) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.GetByID(c, id)
	_routes.RouteResult(c, result, err)
}

func (r *CrudOwnerRouter[T, D]) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	result, err := r.UC.Detail(c, id)
	_routes.RouteResult(c, result, err)
}

func (r *CrudOwnerRouter[T, D]) Create(c *gin.Context) {
	var dto T
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err := r.UC.Create(c, &dto)
	_routes.RouteResult(c, nil, err)
}

func (r *CrudOwnerRouter[T, D]) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	var dto T
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err = r.UC.Update(c, id, &dto)
	_routes.RouteResult(c, nil, err)
}

func (r *CrudOwnerRouter[T, D]) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}
	err = r.UC.Delete(c, id)
	_routes.RouteResult(c, nil, err)
}
