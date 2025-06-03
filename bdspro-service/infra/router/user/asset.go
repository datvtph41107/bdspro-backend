package user_router

import (
	"bdspro/internal/dto"
	user_dto "bdspro/internal/dto/user"
	"bdspro/internal/usecases"
	_routes "common/routes"
	_utils "common/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UAssetRoute struct {
	// crud2.BaseRouter[domain.Asset, *usecases.UAssetUsecase]
	UC *usecases.UAssetUsecase
}

func NewUAssetRoute(service *usecases.UAssetUsecase) *UAssetRoute {
	return &UAssetRoute{
		// BaseRouter: crud2.BaseRouter[domain.Asset, *usecases.UAssetUsecase]{
		// },
		UC: service,
	}
}

func (r *UAssetRoute) RegisterRoutes(router *gin.RouterGroup) {
	assetGroup := router.Group("/assets")
	{
		assetGroup.GET("/internal", r.SearchMe)
		assetGroup.GET("/share", r.SearchShare)
		assetGroup.GET("/split", r.SearchSplit)
		assetGroup.PUT("/archived", r.Archived)
		assetGroup.POST("", r.CreateAsset)
		assetGroup.POST("/from-product", r.CreateAssetFromProduct)

		assetGroup.PUT(":id", r.UpdateAsset)
		assetGroup.DELETE(":id", r.DeleteAsset)
		assetGroup.GET("/:id/detail", r.GetDetail)
		assetGroup.POST("/merge", r.Merge)
		assetGroup.POST("/split", r.Split)
		assetGroup.GET("/history/:id", r.History)
		assetGroup.GET("/history-owner", r.HistoryOwner)
		assetGroup.GET("/publish/:profileId", r.GetPublish)
		// assetGroup.GET(":id/share", r.BaseRouter.GetByID)
	}
}

func (r *UAssetRoute) CreateAsset(c *gin.Context) {
	body := user_dto.AssetSaveRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.CreateAsset(c, &body)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) CreateAssetFromProduct(c *gin.Context) {
	body := user_dto.AssetSaveRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.CreateAssetFromProduct(c, &body)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) SearchMe(c *gin.Context) {
	body := user_dto.AssetSearchDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.SearchAssetMe(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) SearchShare(c *gin.Context) {
	body := user_dto.AssetSearchDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.SearchAssetShare(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) SearchSplit(c *gin.Context) {
	body := user_dto.AssetSearchDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, total, err := r.UC.GetAssetsWithChildCount(c, body)
	_routes.RouteResult(c, _routes.ResponseDTO{
		Code:          0,
		Data:          result,
		TotalElements: &total,
	}, err)
}

func (r *UAssetRoute) Archived(c *gin.Context) {
	body := user_dto.AssetArchivedDTO{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	err := r.UC.ArchiveAsset(c, body.ID, body.Archived)
	_routes.RouteResult(c, nil, err)
}

func (r *UAssetRoute) Merge(c *gin.Context) {
	body := user_dto.MergeAssetDTO{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.MergeAsset(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) Split(c *gin.Context) {
	body := user_dto.SplitAssetDTO{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.SplitAsset(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) History(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	body := dto.AssetHistoryDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.History(c, id, body)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) HistoryOwner(c *gin.Context) {
	body := dto.AssetHistoryDTO{}
	if err := _utils.ParseQuery2(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.HistoryOwner(c, body)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) UpdateAsset(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	body := user_dto.AssetSaveRequest{}
	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
		_routes.RouteResult(c, nil, err)
		return
	}

	result, err := r.UC.UpdateAsset(c, id, &body)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) GetDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	result, err := r.UC.GetAssetDetail(c, id)
	_routes.RouteResult(c, result, err)
}

func (r *UAssetRoute) DeleteAsset(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	err := r.UC.DeleteAsset(c, id)
	_routes.RouteResult(c, nil, err)
}

func (r *UAssetRoute) GetPublish(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("profileId"), 10, 32)

	result, err := r.UC.GetPublish(c, id)
	_routes.RouteResult(c, result, err)
}
