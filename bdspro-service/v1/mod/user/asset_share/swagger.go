package u_asset_share

// @Summary Lấy danh sách
// @Tags Chia sẻ tài sản
// @Produce json
// @Param query query u_asset.SearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share [get]
func GetAsset() {}

// @Summary Xem chi tiết
// @Tags Chia sẻ tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share/{id} [get]
func DetailAsset() {}

// @Summary Tạo mới
// @Tags Chia sẻ tài sản
// @Produce json
// @Param body body models.AssetShare true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share [post]
func CreateAsset() {}

// @Summary Lấy theo ID
// @Tags Chia sẻ tài sản
// @Produce json
// @Param id path int true "ID"
// @Param body body models.AssetShare true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share/{id} [put]
func UpdateAsset() {}

// @Summary Xóa
// @Tags Chia sẻ tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/asset/share/{id} [delete]
func DeleteAsset() {}
