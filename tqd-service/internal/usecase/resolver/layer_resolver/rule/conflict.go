package rule

import (
	"tqd/internal/usecase/resolver/layer_resolver/types"
)

type ConflictEngine struct {
	config *types.Config
}

func NewConflictEngine(config *types.Config) *ConflictEngine {
	return &ConflictEngine{config: config}
}

func (e *ConflictEngine) HasConflict(a, b *types.LayerCandidate) bool {
	if a.GroupCode == b.GroupCode {
		return false
	}
	if row, ok := e.config.ConflictMatrix[a.GroupCode]; ok {
		if val, ok := row[b.GroupCode]; ok {
			return val
		}
	}
	if row, ok := e.config.ConflictMatrix[b.GroupCode]; ok {
		if val, ok := row[a.GroupCode]; ok {
			return val
		}
	}
	return false
}
