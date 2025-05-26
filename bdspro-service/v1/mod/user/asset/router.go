package u_asset

import (
	"bdspro/v1/models"
	"common/case/crud2"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AssetRouter struct {
	crud2.BaseRouter[models.Asset, *UserAssetService]
}

func NewAssetRouter(service *UserAssetService) *AssetRouter {
	return &AssetRouter{
		BaseRouter: crud2.BaseRouter[models.Asset, *UserAssetService]{
			Service: service,
		},
	}
}

func (r *AssetRouter) RegisterRoutes(router *gin.RouterGroup) {
	assetGroup := router.Group("/assets")
	{
		assetGroup.GET("/internal", r.SearchMe)
		assetGroup.GET("/share", r.SearchShare)
		assetGroup.PUT("/archived", r.Archived)
		assetGroup.POST("", r.CreateAsset)
		assetGroup.POST("/from-product", r.CreateAssetFromProduct)

		assetGroup.PUT(":id", r.UpdateAsset)
		assetGroup.DELETE(":id", r.BaseRouter.Delete)
		assetGroup.GET(":id/detail", r.BaseRouter.GetByID)
		assetGroup.POST("/merge", r.Merge)
		assetGroup.POST("/split", r.Split)
		assetGroup.GET("/history/:id", r.History)
		assetGroup.GET("/history-owner", r.HistoryOwner)
		// assetGroup.GET(":id/share", r.BaseRouter.GetByID)
	}
}

func (r *AssetRouter) CreateAsset(c *gin.Context) {
	body := models.Asset{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.CreateAsset(c, &body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) CreateAssetFromProduct(c *gin.Context) {
	body := models.Asset{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.CreateAssetFromProduct(c, &body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) SearchMe(c *gin.Context) {
	body := SearchDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.SearchAssetMe(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) SearchShare(c *gin.Context) {
	body := SearchDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.SearchAssetShare(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) Archived(c *gin.Context) {
	body := ArchivedDTO{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err := r.Service.ArchiveAsset(c, body.ID, body.Archived)
	_routes.RouteResult(c, nil, err)
}

func (r *AssetRouter) Merge(c *gin.Context) {
	body := MergeDTO{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.MergeAsset(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) Split(c *gin.Context) {
	body := SplitDTO{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.SplitAsset(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) History(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	body := HistoryDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.History(c, id, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) HistoryOwner(c *gin.Context) {
	body := HistoryDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.HistoryOwner(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *AssetRouter) UpdateAsset(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	body := models.Asset{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.Service.UpdateAsset(c, id, &body)
	_routes.RouteResult(c, result, err)
}
