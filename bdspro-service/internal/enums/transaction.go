package enums

type TransactionStatus int

const (
	TxNot        TransactionStatus = 10 // Chưa giao dịch
	TxProcessing TransactionStatus = 20 // Đang giao dịch
	TxCompleted  TransactionStatus = 30 // Hoàn tất giao dịch
	TxFail       TransactionStatus = 40 // Giao dịch thất bại
)

var TransactionStatusNames = map[TransactionStatus]string{
	TxNot:        "Chưa giao dịch",
	TxProcessing: "Đang giao dịch",
	TxCompleted:  "Hoàn tất giao dịch",
	TxFail:       "Giao dịch thất bại",
}
