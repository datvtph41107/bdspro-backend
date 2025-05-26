package u_asset_legal

// @Summary Lấy danh sách
// @Tags Thu - chi
// @Produce json
// @Param query query u_asset_cost.SearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-legal [get]
func GetAsset() {}

// @Summary Xem chi tiết
// @Tags Thu - chi
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-legal/{id} [get]
func DetailAsset() {}

// @Summary Tạo mới
// @Tags Thu - chi
// @Produce json
// @Param body body models.AssetLegal true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-legal [post]
func CreateAsset() {}

// @Summary Lấy theo ID
// @Tags Thu - chi
// @Produce json
// @Param id path int true "ID"
// @Param body body models.AssetLegal true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-legal/{id} [put]
func UpdateAsset() {}

// @Summary Xóa
// @Tags Thu - chi
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset-legal/{id} [delete]
func DeleteAsset() {}
