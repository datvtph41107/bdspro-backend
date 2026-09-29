package dto

// LinkProductAssetRequest DTO cho request liên kết nhiều product với asset
type LinkProductAssetRequest struct {
	ProductIDs []uint64 `json:"productIds" binding:"required,min=1"` // Nhiều sản phẩm
	AssetID    uint64   `json:"assetId" binding:"required"`          // 1 tài sản
}

// UnlinkProductAssetRequest DTO cho request hủy liên kết product với asset
type UnlinkProductAssetRequest struct {
	ProductID uint64 `json:"productId" binding:"required"`
	AssetID   uint64 `json:"assetId" binding:"required"`
}

// ProductAssetListResponse DTO cho response danh sách asset theo product
type ProductAssetListResponse struct {
	AssetIDs []uint64 `json:"assetIds"`
}

// AssetProductListResponse DTO cho response danh sách product theo asset
type AssetProductListResponse struct {
	ProductIDs []uint64 `json:"productIds"`
}
