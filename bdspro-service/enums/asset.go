package enums

type AssetStatus int

const (
	AssetOwning  AssetStatus = 1
	AssetRenting AssetStatus = 2
	AssetSold    AssetStatus = 3
)

func (s AssetStatus) String() string {
	switch s {
	case AssetOwning:
		return "Đang sở hữu"
	case AssetRenting:
		return "Đang cho thuê"
	case AssetSold:
		return "Đã bán"
	default:
		return ""
	}
}
