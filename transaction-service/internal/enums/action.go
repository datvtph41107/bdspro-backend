package enums

type EAction uint32

const (
	EActionDeposite EAction = 10
	EActionSign     EAction = 20
	EActionPay      EAction = 30
	EActionHandOver EAction = 40
	EActionCancel   EAction = 50
)

var ActionNames = map[EAction]string{
	EActionDeposite: "Đặt cọc",
	EActionSign:     "Ký hợp đồng",
	EActionPay:      "Thanh toán",
	EActionHandOver: "Bàn giao",
	EActionCancel:   "Hủy",
}
