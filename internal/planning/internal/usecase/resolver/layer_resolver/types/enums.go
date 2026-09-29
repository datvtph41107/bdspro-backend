package types

// ResolveMode - Chế độ xử lý
type ResolveMode string

const (
	ModePopup    ResolveMode = "popup"
	ModeOverview ResolveMode = "overview"
	ModeDetail   ResolveMode = "detail"
	ModeCompare  ResolveMode = "compare"
)

func (m ResolveMode) IsValid() bool {
	switch m {
	case ModePopup, ModeOverview, ModeDetail, ModeCompare:
		return true
	default:
		return false
	}
}

type ConflictSeverity string

const (
	SeverityHigh   ConflictSeverity = "high"
	SeverityMedium ConflictSeverity = "medium"
	SeverityLow    ConflictSeverity = "low"
)
