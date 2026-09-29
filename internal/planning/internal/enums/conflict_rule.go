package enums

type ConflictSeverity uint32

const (
	ConflictSeverityUnknown ConflictSeverity = 0
	ConflictSeverityHigh    ConflictSeverity = 1
	ConflictSeverityMedium  ConflictSeverity = 2
	ConflictSeverityLow     ConflictSeverity = 3
)

func (s ConflictSeverity) String() string {
	switch s {
	case ConflictSeverityHigh:
		return "high"
	case ConflictSeverityMedium:
		return "medium"
	case ConflictSeverityLow:
		return "low"
	default:
		return "unknown"
	}
}

func (s ConflictSeverity) DisplayName() string {
	switch s {
	case ConflictSeverityHigh:
		return "Nghiêm trọng"
	case ConflictSeverityMedium:
		return "Trung bình"
	case ConflictSeverityLow:
		return "Nhẹ"
	default:
		return "Không xác định"
	}
}
