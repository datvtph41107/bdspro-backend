package enums

type TxCostType uint

const (
	TxCostTypeExpense TxCostType = 10
	TxCostTypeRevenue TxCostType = 20
)

var TxCostTypeMap = map[TxCostType]string{
	TxCostTypeExpense: "Chi phí",
	TxCostTypeRevenue: "Doanh thu",
}

