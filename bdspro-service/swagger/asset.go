package swagger_gen

// // @Summary Lấy danh sách tài sản của người dùng
// // @Tags Tài sản
// // @Produce json
// // @Param query query user_dto.AssetSearchDTO true "Size"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/search [get]
// func GetAsset() {}

// // @Summary Lấy danh sách tài sản được chia sẻ
// // @Tags Tài sản
// // @Produce json
// // @Param query query user_dto.AssetSearchDTO true "Size"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/internal [get]
// func GetAssetInternal() {}

// // @Summary Xem chi tiết
// // @Tags Tài sản
// // @Produce json
// // @Param id path int true "ID"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/{id} [get]
// func DetailAsset() {}

// // @Summary Tạo mới
// // @Tags Tài sản
// // @Produce json
// // @Param body body domain.Asset true "Body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets [post]
// func CreateAsset() {}

// // @Summary Lấy theo ID
// // @Tags Tài sản
// // @Produce json
// // @Param id path int true "ID"
// // @Param body body domain.Asset true "Body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/{id} [put]
// func UpdateAsset() {}

// // @Summary Lấy theo ID
// // @Tags Tài sản
// // @Produce json
// // @Param id path int true "ID"
// // @Param body body user_dto.AssetArchivedDTO true "Body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/archived/{id} [put]
// func ArchivedAsset() {}

// // @Summary Xóa
// // @Tags Tài sản
// // @Produce json
// // @Param id path int true "ID"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/{id} [delete]
// func DeleteAsset() {}

// // @Summary Gộp tài sản
// // @Tags Tài sản
// // @Produce json
// // @Param body body user_dto.MergeAssetDTO true "body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/merge [post]
// func MergeAsset() {}

// // @Summary Tách tài sản
// // @Tags Tài sản
// // @Produce json
// // @Param body body user_dto.SplitAssetDTO true "body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/split [post]
// func SplitAsset() {}

// // @Summary Lịch sử tách gộp tài sản danh cho tài sản
// // @Tags Tài sản
// // @Produce json
// // @Param body query dto.AssetHistoryDTO true "body"
// // @Param id path uint64 true "ID"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/history/:id [get]
// func HistoryAsset() {}

// // @Summary Lịch sử tách gộp tài sản dành cho user
// // @Tags Tài sản
// // @Produce json
// // @Param body query dto.AssetHistoryDTO true "body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/history-owner [get]
// func HistoryAssetOwner() {}

// // @Summary Lấy danh sách tài sản có số lượng con
// // @Tags Tài sản
// // @Produce json
// // @Param body query user_dto.AssetSearchDTO true "body"
// // @Security BearerAuth
// // @Router /v1/bdspro/v2/user/assets/split [get]
// func GetAssetsWithChildCount() {}
