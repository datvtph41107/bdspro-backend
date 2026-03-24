package enums

import "fmt"

type ImpactLevel int32

const (
	ImpactLevelNone     ImpactLevel = 0
	ImpactLevelMinor    ImpactLevel = 10
	ImpactLevelModerate ImpactLevel = 20
	ImpactLevelMajor    ImpactLevel = 30
	ImpactLevelCritical ImpactLevel = 40
)

func (e ImpactLevel) String() string {
	switch e {
	case ImpactLevelNone:
		return "NO_IMPACT"
	case ImpactLevelMinor:
		return "MINOR_IMPACT"
	case ImpactLevelModerate:
		return "MODERATE_IMPACT"
	case ImpactLevelMajor:
		return "MAJOR_IMPACT"
	case ImpactLevelCritical:
		return "CRITICAL_IMPACT"
	default:
		return "UNKNOWN"
	}
}

func (e ImpactLevel) IsValid() bool {
	return e >= ImpactLevelNone && e <= ImpactLevelCritical
}

func (e ImpactLevel) RequiresConfirmation() bool {
	return e >= ImpactLevelModerate
}

func (e ImpactLevel) RequiresBlock() bool {
	return e >= ImpactLevelCritical
}

func (e ImpactLevel) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, e.String())), nil
}
