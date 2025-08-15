package enums

type EOwnerType uint32

const (
	EOwnerUser    EOwnerType = 10
	EOwnerDeal    EOwnerType = 20
	EOwnerProduct EOwnerType = 30
)

var DomainNames = map[EOwnerType]string{
	EOwnerUser:    "Người dùng",
	EOwnerDeal:    "Thương vụ",
	EOwnerProduct: "Sản phẩm",
}
