package u_asset_cost

// @Summary Lấy danh sách
// @Tags Chi phí tài sản
// @Produce json
// @Param query query u_asset_cost.SearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-cost [get]
func GetAsset() {}

// @Summary Xem chi tiết
// @Tags Chi phí tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-cost/{id} [get]
func DetailAsset() {}

// @Summary Tạo mới
// @Tags Chi phí tài sản
// @Produce json
// @Param body body models.AssetCost true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-cost [post]
func CreateAsset() {}

// @Summary Lấy theo ID
// @Tags Chi phí tài sản
// @Produce json
// @Param id path int true "ID"
// @Param body body models.AssetCost true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-cost/{id} [put]
func UpdateAsset() {}

// @Summary Xóa
// @Tags Chi phí tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-cost/{id} [delete]
func DeleteAsset() {}
