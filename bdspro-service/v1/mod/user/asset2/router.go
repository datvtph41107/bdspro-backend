package u_asset2

import (
	"bdspro/v1/models"
	"common/case/rest"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetRouter struct {
	rest.RestRouter[models.Asset, SearchDTO]
	AssetService *AssetService
}

func NewAssetRouter2(service *AssetService) *AssetRouter {
	return &AssetRouter{
		RestRouter: rest.RestRouter[models.Asset, SearchDTO]{
			Service: service,
		},
		AssetService: service,
	}
}

func (r *AssetRouter) RegisterRoutes(router *gin.RouterGroup) {
	assetGroup := router.Group("/assets")
	{
		assetGroup.GET("/search", r.Search)
		assetGroup.GET("/internal", r.SearchInternal)
		assetGroup.PUT("/archived", r.Archived)
		assetGroup.POST("", r.CreateAsset)
		assetGroup.PUT(":id", r.RestRouter.Update)
		assetGroup.DELETE(":id", r.RestRouter.Delete)
		assetGroup.GET(":id/detail", r.RestRouter.GetByID)
		assetGroup.POST("/merge", r.Merge)
		assetGroup.POST("/split", r.Split)
		assetGroup.GET("/history/:id", r.History)
		assetGroup.GET("/history-owner", r.HistoryOwner)
		// assetGroup.GET(":id/share", r.BaseRouter.GetByID)
	}
}

func (r *AssetRouter) CreateAsset(c *gin.Context) {
	body := rest.RestModel[models.Asset]{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.AssetService.CreateAsset(c, &body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) Search(c *gin.Context) {
	body := SearchDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.AssetService.SearchAssets(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) SearchInternal(c *gin.Context) {
	body := SearchDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.AssetService.SearchAssetInternal(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) Archived(c *gin.Context) {
	body := ArchivedDTO{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err := r.AssetService.ArchiveAsset(c, body.ID, body.Archived)
	_routes.RouteResult(c, nil, err)
}

func (r *AssetRouter) Merge(c *gin.Context) {
	body := MergeDTO{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.AssetService.MergeAsset(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) Split(c *gin.Context) {
	body := SplitDTO{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.AssetService.SplitAsset(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) History(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	body := HistoryDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.AssetService.History(c, id, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) HistoryOwner(c *gin.Context) {
	body := HistoryDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.AssetService.HistoryOwner(c, body)
	_routes.RouteResult(c, result, err)
}
