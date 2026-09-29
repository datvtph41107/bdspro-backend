package enums

type TransactionType uint

const (
	TransactionTypeExpense   TransactionType = 10 // Chi phí
	TransactionTypeRevenue   TransactionType = 20 // Doanh thu
	TransactionTypeBuy       TransactionType = 30 // Mua
	TransactionTypeSell      TransactionType = 40 // Bán
	TransactionTypeRent      TransactionType = 50 // Số dư
	TransactionTypeNegotiate TransactionType = 60 // Góp vốn
)

var TransactionTypeMap = map[TransactionType]string{
	TransactionTypeExpense:   "Chi phí",
	TransactionTypeRevenue:   "Doanh thu",
	TransactionTypeBuy:       "Mua",
	TransactionTypeSell:      "Bán",
	TransactionTypeRent:      "Thuê",
	TransactionTypeNegotiate: "Góp vốn",
}
