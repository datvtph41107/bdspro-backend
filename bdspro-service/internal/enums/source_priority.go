package enums

type EPriority uint32

const (
	PriorityLow    EPriority = 10
	PriorityNormal EPriority = 20
	PriorityHigh   EPriority = 30
)

var PriorityMap = map[EPriority]string{
	PriorityLow:    "Thấp",
	PriorityNormal: "Bình thường",
	PriorityHigh:   "Cao",
}
