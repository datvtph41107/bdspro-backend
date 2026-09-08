package enums

type TxAction uint32

const (
	TxActionDeposite  TxAction = 10
	TxActionSign      TxAction = 20
	TxActionPay       TxAction = 30
	TxActionHandOver  TxAction = 40
	TxActionCancel    TxAction = 50
	TxActionNegotiate TxAction = 60
)

var TxActionNames = map[TxAction]string{
	TxActionDeposite:  "Đặt cọc",
	TxActionSign:      "Ký hợp đồng",
	TxActionPay:       "Thanh toán",
	TxActionHandOver:  "Bàn giao",
	TxActionCancel:    "Hủy",
	TxActionNegotiate: "Đàm phán",
}

// IsValid kiểm tra giá trị TxAction có hợp lệ không
func (a TxAction) IsValid() bool {
	switch a {
	case TxActionDeposite,
		TxActionSign,
		TxActionPay,
		TxActionHandOver,
		TxActionCancel,
		TxActionNegotiate:
		return true
	default:
		return false
	}
}
