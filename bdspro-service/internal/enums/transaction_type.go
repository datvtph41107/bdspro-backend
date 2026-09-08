package enums

type DealTransactionType uint

const (
	DealTransactionTypeExpense   DealTransactionType = 10 // Chi phí
	DealTransactionTypeRevenue   DealTransactionType = 20 // Doanh thu
	DealTransactionTypeBuy       DealTransactionType = 30 // Mua
	DealTransactionTypeSell      DealTransactionType = 40 // Bán
	DealTransactionTypeRent      DealTransactionType = 50 // Số dư
	DealTransactionTypeNegotiate DealTransactionType = 60 // Góp vốn
)

var DealTransactionTypeMap = map[DealTransactionType]string{
	DealTransactionTypeExpense:   "Chi phí",
	DealTransactionTypeRevenue:   "Doanh thu",
	DealTransactionTypeBuy:       "Mua",
	DealTransactionTypeSell:      "Bán",
	DealTransactionTypeRent:      "Thuê",
	DealTransactionTypeNegotiate: "Góp vốn",
}
