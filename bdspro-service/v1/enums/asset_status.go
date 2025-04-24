package enums

type AssetStatus uint

const (
	AssetStatusOwning  AssetStatus = 1 // đang sở hữu
	AssetStatusRenting AssetStatus = 2 // đang cho thuê
	AssetStatusSold    AssetStatus = 3 // đã bán
	AssetStatusSelling AssetStatus = 4 // đang bán
	AssetStatusNotSold AssetStatus = 4 // đang bán
)

// bỏ đã sở hữu

func (s AssetStatus) String() string {
	switch s {
	case AssetStatusOwning:
		return "Đang sở hữu"
	case AssetStatusRenting:
		return "Đang cho thuê"
	case AssetStatusSold:
		return "Đã bán"
	default:
		return ""
	}
}
