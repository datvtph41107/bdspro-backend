package enums

type TxOwnerType uint32

const (
	TxOwnerUser         TxOwnerType = 10
	TxOwnerDealContract TxOwnerType = 20
	TxOwnerProduct      TxOwnerType = 30
)

var TxDomainNames = map[TxOwnerType]string{
	TxOwnerUser:         "Người dùng",
	TxOwnerDealContract: "Hợp đồng thương vụ",
	TxOwnerProduct:      "Sản phẩm",
}
