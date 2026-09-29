package enums

type ECostType uint

const (
	ECostTypeExpense ECostType = 10
	ECostTypeRevenue ECostType = 20
)

var ECostTypeMap = map[ECostType]string{
	ECostTypeExpense: "Doanh thu",
	ECostTypeRevenue: "Chi phí",
}
