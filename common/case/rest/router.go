package rest

import (
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// type IRestRouter[T IRestModel, DTO ISearchDTO] interface {
// }

type RestRouter[M IRestModel, DTO any] struct {
	Service IRestService[M, DTO]
}

func (r *RestRouter[M, DTO]) RegisterRoutes(router *gin.RouterGroup, path string) {
	assetGroup := router.Group(path)
	{
		assetGroup.GET("", r.Search)
		assetGroup.POST("", r.Create)
		assetGroup.PUT(":id", r.Update)
		assetGroup.DELETE(":id", r.Delete)
		assetGroup.GET(":id", r.GetByID)
	}
}

func (r *RestRouter[M, DTO]) Create(c *gin.Context) {
	entity := r.Service.GenModel()
	// var e := *entity
	// e = entity.(T)

	if err := _utils.ParseBodyWithValidator(c, &entity.M); err != nil {
		_routes.RouteResult(c, nil, err)

		return
	}

	result, err := r.Service.Create(c, &entity)
	_routes.RouteResult(c, result, err)
}

func (r *RestRouter[T, DTO]) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	entity := r.Service.GenModel()
	if err := _utils.ParseBodyWithValidator(c, &entity.M); err != nil {
		_routes.RouteResult(c, nil, err)

		return
	}

	result, err := r.Service.Update(c, id, &entity)
	_routes.RouteResult(c, result, err)
}

func (r *RestRouter[T, DTO]) Delete(c *gin.Context) {
	// success, err := r.Service.Delete(c)
	// _routes.RouteResult(c, success, err)
}

func (r *RestRouter[T, DTO]) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "ID không hợp lệ",
		})
		return
	}

	entity, err := r.Service.GetByID(c, id)
	_routes.RouteResult(c, entity, err)
}

func (r *RestRouter[T, DTO]) Search(c *gin.Context) {
	var dto SearchDTO[DTO]
	if err := _utils.ParseQuery2(c, &dto); err != nil {
		_routes.RouteResult(c, nil, &_routes.Except{
			Code:    400,
			Message: "Dữ liệu không hợp lệ",
		})

		return
	}

	entities, err := r.Service.Search(c, dto)
	_routes.RouteResult(c, entities, err)
}
