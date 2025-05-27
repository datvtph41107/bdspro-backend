package data

type OrderType int8

const (
	OrderBill       OrderType = 1
	OrderMembership OrderType = 2
)

type OrderRequest struct {
	Amount      int
	Description string
	TargetId    *uint64
	Type        OrderType
}
