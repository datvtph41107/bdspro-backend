package org_router

import (
	shared_usecase "bdspro/internal/usecases/shared"
)

type OrgAssetRouter struct {
	UC *shared_usecase.AssetUsecase
}

func NewOrgAssetRouter(uc *shared_usecase.AssetUsecase) *OrgAssetRouter {
	return &OrgAssetRouter{UC: uc}
}

// func (r *OrgAssetRouter) RegisterRoutes(router *gin.RouterGroup, path string) {
// 	assetGroup := router.Group(path)
// 	{
// 		assetGroup.GET("/internal", r.SearchMe)
// 		assetGroup.GET("/share", r.SearchShare)
// 		assetGroup.GET("/split", r.SearchSplit)
// 		assetGroup.PUT("/archived", r.Archived)
// 		// assetGroup.POST("", r.CreateAsset)
// 		// assetGroup.POST("/from-product", r.CreateAssetFromProduct)

// 		// assetGroup.PUT(":id", r.UpdateAsset)
// 		assetGroup.DELETE(":id", r.DeleteAsset)
// 		assetGroup.GET("/:id/detail", r.GetDetail)
// 		assetGroup.POST("/merge", r.Merge)
// 		assetGroup.POST("/split", r.Split)
// 		// assetGroup.GET("/history/:id", r.History)
// 		// assetGroup.GET("/history-owner", r.HistoryOwner)
// 		assetGroup.GET("/publish/:profileId", r.GetPublish)
// 		// assetGroup.GET(":id/share", r.BaseRouter.GetByID)
// 	}
// }

// // // @Summary Tạo mới
// // // @Tags Tổ chức: Tài sản
// // // @Produce json
// // // @Param body body domain.Asset true "Body"
// // // @Security BearerAuth
// // // @Router /v1/bdspro/v2/org/asset [post]
// // func (r *OrgAssetRouter) CreateAsset(c *gin.Context) {
// // 	body := dto.AssetSaveRequest{}
// // 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// // 		_routes.RouteResult(c, nil, err)
// // 		return
// // 	}

// // 	result, err := r.UC.CreateAsset(c, &body)
// // 	_routes.RouteResult(c, result, err)
// // }

// // // @Summary Tạo từ sản phẩm
// // // @Tags Tổ chức: Tài sản
// // // @Produce json
// // // @Param body body dto.AssetSaveRequest true "Body"
// // // @Security BearerAuth
// // // @Router /v1/bdspro/v2/org/asset/from-product [post]
// // func (r *OrgAssetRouter) CreateAssetFromProduct(c *gin.Context) {
// // 	body := dto.AssetSaveRequest{}
// // 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// // 		_routes.RouteResult(c, nil, err)
// // 		return
// // 	}

// // 	result, err := r.UC.CreateAssetFromProduct(c, &body)
// // 	_routes.RouteResult(c, result, err)
// // }

// // @Summary Lấy danh sách tài sản của người dùng
// // @Tags Tổ chức: Tài sản
// // @Produce json
// // @Param query query dto.AssetSearchDTO true "Size"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/org/asset/search [get]
// func (r *OrgAssetRouter) SearchMe(c *gin.Context) {
// 	body := dto.AssetSearchDTO{}
// 	if err := _utils.ParseQuery2(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, total, err := r.UC.SearchAssetMe(c, body)
// 	_routes.RouteResult(c, _routes.ResponseDTO{
// 		Code:          0,
// 		Data:          result,
// 		TotalElements: &total,
// 	}, err)
// }

// func (r *OrgAssetRouter) SearchShare(c *gin.Context) {
// 	body := dto.AssetSearchDTO{}
// 	if err := _utils.ParseQuery2(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, total, err := r.UC.SearchAssetShare(c, body)
// 	_routes.RouteResult(c, _routes.ResponseDTO{
// 		Code:          0,
// 		Data:          result,
// 		TotalElements: &total,
// 	}, err)
// }

// // @Summary Lấy danh sách tài sản có số lượng con
// // @Tags Tổ chức: Tài sản
// // @Produce json
// // @Param body query dto.AssetSearchDTO true "body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/org/asset/split [get]
// func (r *OrgAssetRouter) SearchSplit(c *gin.Context) {
// 	body := dto.AssetSearchDTO{}
// 	if err := _utils.ParseQuery2(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, total, err := r.UC.GetAssetsWithChildCount(c, body)
// 	_routes.RouteResult(c, _routes.ResponseDTO{
// 		Code:          0,
// 		Data:          result,
// 		TotalElements: &total,
// 	}, err)
// }

// // @Summary Lấy theo ID
// // @Tags Tổ chức: Tài sản
// // @Produce json
// // @Param id path int true "ID"
// // @Param body body dto.AssetArchivedDTO true "Body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/org/asset/archived/{id} [put]
// func (r *OrgAssetRouter) Archived(c *gin.Context) {
// 	body := dto.AssetArchivedDTO{}
// 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	err := r.UC.ArchiveAsset(c, body.ID, body.Archived)
// 	_routes.RouteResult(c, nil, err)
// }

// // @Summary Gộp tài sản
// // @Tags Tổ chức: Tài sản
// // @Produce json
// // @Param body body dto.MergeAssetDTO true "body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/org/asset/merge [post]
// func (r *OrgAssetRouter) Merge(c *gin.Context) {
// 	body := dto.MergeAssetDTO{}
// 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, err := r.UC.MergeAsset(c, body)
// 	_routes.RouteResult(c, result, err)
// }

// // @Summary Tách tài sản
// // @Tags Tổ chức: Tài sản
// // @Produce json
// // @Param body body dto.SplitAssetDTO true "body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/org/asset/split [post]
// func (r *OrgAssetRouter) Split(c *gin.Context) {
// 	body := dto.SplitAssetDTO{}
// 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// 		_routes.RouteResult(c, nil, err)
// 		return
// 	}

// 	result, err := r.UC.SplitAsset(c, body)
// 	_routes.RouteResult(c, result, err)
// }

// // // @Summary Lịch sử tách gộp tài sản danh cho tài sản
// // // @Tags Tổ chức: Tài sản
// // // @Produce json
// // // @Param body query dto.AssetHistoryDTO true "body"
// // // @Param id path uint64 true "ID"
// // // @Security BearerAuth
// // // @Router /v1/bdspro/v2/org/asset/history/:id [get]
// // func (r *OrgAssetRouter) History(c *gin.Context) {
// // 	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

// // 	body := dto.AssetHistoryDTO{}
// // 	if err := _utils.ParseQuery2(c, &body); err != nil {
// // 		_routes.RouteResult(c, nil, err)
// // 		return
// // 	}

// // 	result, err := r.UC.History(c, id, body)
// // 	_routes.RouteResult(c, result, err)
// // }

// // // @Summary Lịch sử tách gộp tài sản dành cho user
// // // @Tags Tổ chức: Tài sản
// // // @Produce json
// // // @Param body query dto.AssetHistoryDTO true "body"
// // // @Security BearerAuth
// // // @Router /v1/bdspro/v2/org/asset/history-owner [get]
// // func (r *OrgAssetRouter) HistoryOwner(c *gin.Context) {
// // 	body := dto.AssetHistoryDTO{}
// // 	if err := _utils.ParseQuery2(c, &body); err != nil {
// // 		_routes.RouteResult(c, nil, err)
// // 		return
// // 	}

// // 	result, err := r.UC.HistoryOwner(c, body)
// // 	_routes.RouteResult(c, result, err)
// // }

// // func (r *OrgAssetRouter) UpdateAsset(c *gin.Context) {
// // 	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

// // 	body := dto.AssetSaveRequest{}
// // 	if err := _utils.ParseBodyWithValidator(c, &body); err != nil {
// // 		_routes.RouteResult(c, nil, err)
// // 		return
// // 	}

// // 	result, err := r.UC.UpdateAsset(c, id, &body)
// // 	_routes.RouteResult(c, result, err)
// // }

// // @Summary Xem chi tiết
// // @Tags Tổ chức: Tài sản
// // @Produce json
// // @Param id path int true "ID"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/org/asset/{id} [get]
// func (r *OrgAssetRouter) GetDetail(c *gin.Context) {
// 	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

// 	result, err := r.UC.GetAssetDetail(c, id)
// 	_routes.RouteResult(c, result, err)
// }

// // @Summary Xóa
// // @Tags Tổ chức: Tài sản
// // @Produce json
// // @Param id path int true "ID"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/org/asset/{id} [delete]
// func (r *OrgAssetRouter) DeleteAsset(c *gin.Context) {
// 	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

// 	err := r.UC.DeleteAsset(c, id)
// 	_routes.RouteResult(c, nil, err)
// }

// // @Summary Lấy danh sách tài sản đã đăng
// // @Tags Tổ chức: Tài sản
// // @Produce json
// // @Param profileId path int true "ID"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/org/asset/publish/{profileId} [get]
// func (r *OrgAssetRouter) GetPublish(c *gin.Context) {
// 	id, _ := strconv.ParseUint(c.Param("profileId"), 10, 32)

// 	result, err := r.UC.GetPublish(c, id)
// 	_routes.RouteResult(c, result, err)
// }
