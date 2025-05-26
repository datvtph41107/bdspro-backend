package u_asset_cost_type

// @Summary Lấy danh sách loại chi phí
// @Tags Danh mục loại
// @Produce json
// @Param query query u_asset_cost_type.SearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-cost [get]
func GetAsset() {}

// @Summary Xem chi tiết loại chi phí
// @Tags Danh mục loại
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-cost/{id} [get]
func DetailAsset() {}

// @Summary Tạo mới loại chi phí
// @Tags Danh mục loại
// @Produce json
// @Param body body models.AssetCostType true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-cost [post]
func CreateAsset() {}

// @Summary Lấy theo ID loại chi phí
// @Tags Danh mục loại
// @Produce json
// @Param id path int true "ID"
// @Param body body models.AssetCostType true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-cost/{id} [put]
func UpdateAsset() {}

// @Summary Xóa loại chi phí
// @Tags Danh mục loại
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-cost/{id} [delete]
func DeleteAsset() {}
