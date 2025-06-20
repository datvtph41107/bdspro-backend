package enums

type EPriority int

const (
	PriorityLow    EPriority = 10
	PriorityMedium EPriority = 20
	PriorityHigh   EPriority = 30
)

var PriorityMap = map[EPriority]string{
	PriorityLow:    "Thấp",
	PriorityMedium: "Trung bình",
	PriorityHigh:   "Cao",
}
