package enums

type EAssetStatus uint

const (
	EAssetStatusOwning  EAssetStatus = 10 // đang sở hữu
	EAssetStatusRenting EAssetStatus = 20 // đang cho thuê
	EAssetStatusSold    EAssetStatus = 30 // đã bán
	EAssetStatusSelling EAssetStatus = 40 // đang bán
	EAssetStatusNotSold EAssetStatus = 50 // đang bán
)

var EAssetStatusNames = map[EAssetStatus]string{
	EAssetStatusOwning:  "Đang sở hữu",
	EAssetStatusRenting: "Đang cho thuê",
	EAssetStatusSold:    "Đã bán",
	EAssetStatusSelling: "Đang bán",
	EAssetStatusNotSold: "Đang bán",
}

func (e EAssetStatus) IsValid() bool {
	return e == EAssetStatusOwning ||
		e == EAssetStatusRenting ||
		e == EAssetStatusSold ||
		e == EAssetStatusSelling ||
		e == EAssetStatusNotSold
}
