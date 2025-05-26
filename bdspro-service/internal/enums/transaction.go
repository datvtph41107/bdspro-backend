package enums

type TransactionStatus int

const (
	TxNot        TransactionStatus = 10 // Chưa giao dịch
	TxProcessing TransactionStatus = 20 // Đang giao dịch
	TxCompleted  TransactionStatus = 30 // Hoàn tất giao dịch
	TxFail       TransactionStatus = 40 // Giao dịch thất bại
)

func (s TransactionStatus) String() string {
	switch s {
	case TxNot:
		return "Chưa giao dịch"
	case TxProcessing:
		return "Đang giao dịch"
	case TxCompleted:
		return "Hoàn tất giao dịch"
	case TxFail:
		return "Giao dịch thất bại"
	default:
		return ""
	}
}
