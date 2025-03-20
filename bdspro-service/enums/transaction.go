package enums

type TransactionStatus int

const (
	TxNot        TransactionStatus = 1 // Chưa giao dịch
	TxProcessing TransactionStatus = 2 // Đang giao dịch
	TxCompleted  TransactionStatus = 3 // Hoàn tất giao dịch
	TxFail       TransactionStatus = 4 // Giao dịch thất bại
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
