package enums

type EffectiveState int32

const (
	EffectiveStateActive     EffectiveState = 10
	EffectiveStateRestricted EffectiveState = 20
	EffectiveStateFrozen     EffectiveState = 30
	EffectiveStateArchived   EffectiveState = 40
)

func (e EffectiveState) String() string {
	switch e {
	case EffectiveStateActive:
		return "ACTIVE"
	case EffectiveStateRestricted:
		return "RESTRICTED"
	case EffectiveStateFrozen:
		return "FROZEN"
	case EffectiveStateArchived:
		return "ARCHIVED"
	default:
		return "UNKNOWN"
	}
}

func (e EffectiveState) CanOperate() bool {
	return e == EffectiveStateActive
}

func (e EffectiveState) CanEdit() bool {
	return e == EffectiveStateActive || e == EffectiveStateRestricted
}

func (e EffectiveState) NeedsReview() bool {
	return e == EffectiveStateRestricted || e == EffectiveStateFrozen
}
