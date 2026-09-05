package rule

import (
	"tqd/internal/usecase/resolver/layer_resolver/types"
)

type PriorityEngine struct {
	config *types.Config
}

func NewPriorityEngine(config *types.Config) *PriorityEngine {
	return &PriorityEngine{config: config}
}

func (e *PriorityEngine) Calculate(layer *types.LayerCandidate) int {
	if layer.Priority > 0 {
		return layer.Priority
	}
	if p, ok := e.config.DefaultGroupPriority[layer.GroupCode]; ok {
		return p
	}
	return 50
}
