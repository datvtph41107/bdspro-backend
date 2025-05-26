package swagger_gen

import (
	pb_bdspro "pb/types/bdspro"
)

// @Summary Lấy danh sách
// @Tags Thu - Chi tài sản
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param query query pb_bdspro.AssetCostDTO true "Size"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-cost [get]
func GetAssetCost() {}

// @Summary Xem chi tiết
// @Tags Thu - Chi tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-cost/{id} [get]
func DetailAssetCost() {}

// @Summary Tạo mới
// @Tags Thu - Chi tài sản
// @Produce json
// @Param body body pb_bdspro.AssetCostDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-cost [post]
func CreateAssetCost(pb_bdspro.AssetCostDTO) {}

// @Summary Lấy theo ID
// @Tags Thu - Chi tài sản
// @Produce json
// @Param id path int true "ID"
// @Param body body pb_bdspro.AssetCostDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-cost/{id} [put]
func UpdateAssetCost() {}

// @Summary Xóa
// @Tags Thu - Chi tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-cost/{id} [delete]
func DeleteAssetCost() {}
