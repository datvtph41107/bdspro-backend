package types

import "time"

// LifecycleState — trạng thái trong vòng đời layer
type LifecycleState string

const (
	StateDraft     LifecycleState = "draft"
	StateApproved  LifecycleState = "approved"
	StateEffective LifecycleState = "effective"
	StateExpired   LifecycleState = "expired"
	StateReplaced  LifecycleState = "replaced"
	StateSuspended LifecycleState = "suspended"
)

// LifecycleTransition — một bước chuyển trạng thái
type LifecycleTransition struct {
	From      LifecycleState
	To        LifecycleState
	Timestamp time.Time
	Actor     string
	Reason    string
	LegalDoc  string
}

// LayerVersion — một phiên bản trong version chain
type LayerVersion struct {
	Version     int
	LayerID     uint64
	ParentID    *uint64
	FamilyID    uint64
	State       LifecycleState
	Transitions []LifecycleTransition
	CreatedAt   time.Time
	CreatedBy   string
	IsLatest    bool
}

// LayerLineage — toàn bộ cây phả hệ của 1 layer
type LayerLineage struct {
	RootLayerID  uint64
	FamilyID     uint64
	Versions     []*LayerVersion
	CurrentID    uint64
	ReplacedByID *uint64
	Children     []*LayerLineage
}

// LifecycleConfig — cấu hình cho Lifecycle Engine
type LifecycleConfig struct {
	AllowedTransitions   map[LifecycleState][]LifecycleState
	AutoIncrementVersion bool
}

// DefaultAllowedTransitions — các chuyển trạng thái hợp lệ
func DefaultAllowedTransitions() map[LifecycleState][]LifecycleState {
	return map[LifecycleState][]LifecycleState{
		StateDraft:     {StateApproved, StateSuspended},
		StateApproved:  {StateEffective, StateDraft, StateSuspended},
		StateEffective: {StateExpired, StateReplaced, StateSuspended},
		StateExpired:   {},
		StateReplaced:  {},
		StateSuspended: {StateDraft, StateApproved, StateEffective},
	}
}
