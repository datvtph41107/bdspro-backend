package enums

type ECostType uint

const (
	ECostTypeIn  ECostType = 10
	ECostTypeOut ECostType = 20
)

var CostTypeMap = map[ECostType]string{
	ECostTypeIn:  "Loại thu",
	ECostTypeOut: "Loại chi",
}
