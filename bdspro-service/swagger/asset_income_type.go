package swagger_gen

import (
	pb_bdspro "pb/types/bdspro"
)

// @Summary Lấy danh sách
// @Tags Thu - Chi tài sản
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param query query pb_bdspro.AssetIncomeTypeDTO true "Size"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-income-type [get]
func GetAssetIncomeType() {}

// @Summary Xem chi tiết
// @Tags Thu - Chi tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-income-type/{id} [get]
func DetailAssetIncomeType() {}

// @Summary Tạo mới
// @Tags Thu - Chi tài sản
// @Produce json
// @Param body body pb_bdspro.AssetIncomeTypeDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-income-type [post]
func CreateAssetIncomeType(pb_bdspro.AssetIncomeTypeDTO) {}

// @Summary Lấy theo ID
// @Tags Thu - Chi tài sản
// @Produce json
// @Param id path int true "ID"
// @Param body body pb_bdspro.AssetIncomeTypeDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-income-type/{id} [put]
func UpdateAssetIncomeType() {}

// @Summary Xóa
// @Tags Thu - Chi tài sản
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-income-type/{id} [delete]
func DeleteAssetIncomeType() {}
