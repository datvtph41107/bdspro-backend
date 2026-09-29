package enums

type EDomainAccess int32

const (
	EDomainProduct EDomainAccess = 10 // Product
	EDomainAsset   EDomainAccess = 20 // Asset
)

var EDomainAccessMap = map[EDomainAccess]string{
	EDomainProduct: "Sản phẩm",
	EDomainAsset:   "Tài sản",
}
