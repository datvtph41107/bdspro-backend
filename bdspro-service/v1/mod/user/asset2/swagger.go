package u_asset2

// @Summary Lấy danh sách tài sản của người dùng
// @Tags Tài sản
// @Produce json
// @Param query query u_asset.SearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/search [get]
func GetAsset() {}

// @Summary Lấy danh sách tài sản được chia sẻ
// @Tags Tài sản
// @Produce json
// @Param query query u_asset.SearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/internal [get]
func GetAssetInternal() {}

// @Summary Xem chi tiết
// @Tags Tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/{id} [get]
func DetailAsset() {}

// @Summary Tạo mới
// @Tags Tài sản
// @Produce json
// @Param body body models.Asset true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets [post]
func CreateAsset() {}

// @Summary Lấy theo ID
// @Tags Tài sản
// @Produce json
// @Param id path int true "ID"
// @Param body body models.Asset true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/{id} [put]
func UpdateAsset() {}

// @Summary Lấy theo ID
// @Tags Tài sản
// @Produce json
// @Param id path int true "ID"
// @Param body body u_asset.ArchivedDTO true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/archived/{id} [put]
func ArchivedAsset() {}

// @Summary Xóa
// @Tags Tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/{id} [delete]
func DeleteAsset() {}

// @Summary Gộp tài sản
// @Tags Tài sản
// @Produce json
// @Param body body u_asset.MergeDTO true "body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/merge [post]
func MergeAsset() {}

// @Summary Tách tài sản
// @Tags Tài sản
// @Produce json
// @Param body body u_asset.SplitDTOSwagger true "body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/split [post]
func SplitAsset() {}

// @Summary Lịch sử tách gộp tài sản danh cho tài sản
// @Tags Tài sản
// @Produce json
// @Param body query u_asset.HistoryDTO true "body"
// @Param id path uint64 true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/history/:id [get]
func HistoryAsset() {}

// @Summary Lịch sử tách gộp tài sản dành cho user
// @Tags Tài sản
// @Produce json
// @Param body query u_asset.HistoryDTO true "body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/assets/history-owner [get]
func HistoryAssetOwner() {}
