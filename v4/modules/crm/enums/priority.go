package enums

type EPriority int

const (
	PriorityLow      EPriority = 10
	PriorityMedium   EPriority = 20
	PriorityHigh     EPriority = 30
	PriorityCritical EPriority = 40
)

var PriorityMap = map[EPriority]string{
	PriorityLow:      "Thấp",
	PriorityMedium:   "Trung bình",
	PriorityHigh:     "Cao",
	PriorityCritical: "Nghiêm trọng",
}

func (p EPriority) IsValid() bool {
	return p == PriorityLow || p == PriorityMedium || p == PriorityHigh || p == PriorityCritical
}
