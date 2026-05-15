package rule

import (
	"tqd/internal/modules/resolver/layer_resolver/types"
)

type AssessEngine struct {
	config *types.Config
}

func NewAssessEngine(config *types.Config) *AssessEngine {
	return &AssessEngine{config: config}
}

func (e *AssessEngine) Assess(primary *types.RankedLayer, hasConflict bool) (bool, string, string) {
	if !primary.CanBuild {
		return false, primary.BuildCondition, "prohibited"
	}
	if hasConflict {
		return true, primary.BuildCondition + " (Có mâu thuẫn với quy hoạch khác)", "conditional"
	}
	return true, primary.BuildCondition, "allowed"
}
