package u_asset_income_type

// @Summary Lấy danh sách loại thu nhập
// @Tags Danh mục loại
// @Produce json
// @Param query query u_asset_income_type.SearchDTO true "Size"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-income [get]
func GetAsset() {}

// @Summary Xem chi tiết loại thu nhập
// @Tags Danh mục loại
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-income/{id} [get]
func DetailAsset() {}

// @Summary Tạo mới loại thu nhập
// @Tags Danh mục loại
// @Produce json
// @Param body body models.AssetIncomeType true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-income [post]
func CreateAsset() {}

// @Summary Lấy theo ID loại thu nhập
// @Tags Danh mục loại
// @Produce json
// @Param id path int true "ID"
// @Param body body models.AssetIncomeType true "Body"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-income/{id} [put]
func UpdateAsset() {}

// @Summary Xóa loại thu nhập
// @Tags Danh mục loại
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v1/bdspro/v1/user/type/asset-income/{id} [delete]
func DeleteAsset() {}
