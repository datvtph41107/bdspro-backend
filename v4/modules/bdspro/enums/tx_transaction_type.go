package enums

type TxTransactionType uint

const (
	TxTransactionTypeExpense   TxTransactionType = 10 // Chi phí
	TxTransactionTypeRevenue   TxTransactionType = 20 // Doanh thu
	TxTransactionTypeBuy       TxTransactionType = 30 // Mua
	TxTransactionTypeSell      TxTransactionType = 40 // Bán
	TxTransactionTypeRent      TxTransactionType = 50 // Số dư
	TxTransactionTypeNegotiate TxTransactionType = 60 // Góp vốn
)

var TxTransactionTypeMap = map[TxTransactionType]string{
	TxTransactionTypeExpense:   "Chi phí",
	TxTransactionTypeRevenue:   "Doanh thu",
	TxTransactionTypeBuy:       "Mua",
	TxTransactionTypeSell:      "Bán",
	TxTransactionTypeRent:      "Thuê",
	TxTransactionTypeNegotiate: "Góp vốn",
}
