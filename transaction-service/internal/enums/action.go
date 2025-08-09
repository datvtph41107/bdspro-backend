package enums

type EAction uint32

const (
	EActionTransfer EAction = 10
	EActionApprove  EAction = 20
	EActionReject   EAction = 30
	EActionSign     EAction = 40
)

var ActionNames = map[EAction]string{
	EActionTransfer: "Chuyển khoản",
	EActionApprove:  "Phê duyệt",
	EActionReject:   "Từ chối",
	EActionSign:     "Ký hợp đồng", // sign
}
