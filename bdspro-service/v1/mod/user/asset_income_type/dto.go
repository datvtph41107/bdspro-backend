package u_asset_income_type

type SearchDTO struct {
	Page    int    `form:"page"`
	Size    int    `form:"size"`
	AssetID uint64 `form:"assetId"`
}
