package enums

type ETargetType = int

const (
	ETargetTypeProduct ETargetType = 10
	ETargetTypeAsset   ETargetType = 20
	ETargetTypePost    ETargetType = 30
)

var TargetTypeNames = map[ETargetType]string{
	ETargetTypeProduct: "Sản phẩm",
	ETargetTypeAsset:   "Tài sản",
	ETargetTypePost:    "Bài viết",
}
