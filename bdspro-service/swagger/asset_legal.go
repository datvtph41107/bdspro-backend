
package swagger_gen

import (
	pb_bdspro "pb/types/bdspro"
)

// @Summary Lấy danh sách
// @Tags AssetLegal
// @Produce json
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Param query query pb_bdspro.AssetLegalDTO true "Size"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-legal [get]
func GetAssetLegal() {}

// @Summary Xem chi tiết
// @Tags AssetLegal
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-legal/{id} [get]
func DetailAssetLegal() {}

// @Summary Tạo mới
// @Tags AssetLegal
// @Produce json
// @Param body body pb_bdspro.AssetLegalDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-legal [post]
func CreateAssetLegal(pb_bdspro.AssetLegalDTO) {}

// @Summary Lấy theo ID
// @Tags AssetLegal
// @Produce json
// @Param id path int true "ID"
// @Param body body pb_bdspro.AssetLegalDTO true "Body"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-legal/{id} [put]
func UpdateAssetLegal() {}

// @Summary Xóa
// @Tags AssetLegal
// @Produce json
// @Param id path int true "ID"
// @Security BearerAuth
// @Router /v2/bdspro/v2/asset-legal/{id} [delete]
func DeleteAssetLegal() {}

