package enums

type EOwnerType uint32

const (
	EOwnerUser         EOwnerType = 10
	EOwnerDealContract EOwnerType = 20
	EOwnerProduct      EOwnerType = 30
)

var DomainNames = map[EOwnerType]string{
	EOwnerUser:         "Người dùng",
	EOwnerDealContract: "Hợp đồng thương vụ",
	EOwnerProduct:      "Sản phẩm",
}
