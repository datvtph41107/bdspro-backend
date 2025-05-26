package enums

type AssetStatus uint

const (
	AssetStatusOwning  AssetStatus = 10 // đang sở hữu
	AssetStatusRenting AssetStatus = 20 // đang cho thuê
	AssetStatusSold    AssetStatus = 30 // đã bán
	AssetStatusSelling AssetStatus = 40 // đang bán
	AssetStatusNotSold AssetStatus = 50 // đang bán
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
