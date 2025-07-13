package enums

type TransactionType uint

const (
	TransactionTypeExpense TransactionType = 10 // Chi phí
	TransactionTypeRevenue TransactionType = 20 // Doanh thu
)

var TransactionTypeMap = map[TransactionType]string{
	TransactionTypeExpense: "Chi phí",
	TransactionTypeRevenue: "Doanh thu",
} 