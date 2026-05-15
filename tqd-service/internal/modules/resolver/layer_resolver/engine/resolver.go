package engine

import (
	"sort"
	"time"
	"tqd/internal/modules/resolver/layer_resolver/rule"
	"tqd/internal/modules/resolver/layer_resolver/types"
)

type ResolverEngine struct {
	config         *types.Config
	priorityEngine *rule.PriorityEngine
	conflictEngine *rule.ConflictEngine
	assessEngine   *rule.AssessEngine
}

func NewResolverEngine(config *types.Config) *ResolverEngine {
	return &ResolverEngine{
		config:         config,
		priorityEngine: rule.NewPriorityEngine(config),
		conflictEngine: rule.NewConflictEngine(config),
		assessEngine:   rule.NewAssessEngine(config),
	}
}

func (e *ResolverEngine) Resolve(candidates []*types.LayerCandidate, mode types.ResolveMode) *types.ResolveResult {
	startTime := time.Now()

	if len(candidates) == 0 {
		return &types.ResolveResult{
			CanBuild:    false,
			BuildStatus: "unknown",
			Reason:      "Không có dữ liệu quy hoạch",
			Mode:        mode,
		}
	}

	ranked := e.scoreAndRank(candidates)

	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].OverlapPercent > ranked[j].OverlapPercent
		}
		return ranked[i].Score > ranked[j].Score
	})

	for i, r := range ranked {
		r.Rank = i + 1
	}

	primary := ranked[0]
	secondary := ranked[1:]

	criticalConflicts := e.detectConflicts(primary, secondary)
	hasConflict := len(criticalConflicts) > 0

	canBuild, reason, status := e.assessEngine.Assess(primary, hasConflict)

	result := &types.ResolveResult{
		Primary:           primary,
		Secondary:         secondary,
		HasConflict:       hasConflict,
		CriticalConflicts: criticalConflicts,
		CanBuild:          canBuild,
		BuildStatus:       status,
		Reason:            reason,
		Mode:              mode,
		ProcessingTimeMs:  time.Since(startTime).Milliseconds(),
	}

	if mode == types.ModePopup {
		result.Secondary = nil
	}

	return result
}

func (e *ResolverEngine) scoreAndRank(candidates []*types.LayerCandidate) []*types.RankedLayer {
	ranked := make([]*types.RankedLayer, len(candidates))

	for i, c := range candidates {
		legalScore := e.config.LegalScores[c.LegalStatus]
		priorityScore := float64(e.priorityEngine.Calculate(c))
		areaScore := c.OverlapPercent

		totalScore := legalScore*e.config.LegalWeight +
			priorityScore*e.config.PriorityWeight +
			areaScore*e.config.AreaWeight

		ranked[i] = &types.RankedLayer{
			EvaluatedLayer: &types.EvaluatedLayer{
				LayerCandidate: c,
				LegalWeight:    legalScore,
			},
			Score: totalScore,
		}
	}
	return ranked
}

func (e *ResolverEngine) detectConflicts(primary *types.RankedLayer, secondary []*types.RankedLayer) []*types.ConflictDetail {
	var conflicts []*types.ConflictDetail

	for _, sec := range secondary {
		if e.conflictEngine.HasConflict(primary.LayerCandidate, sec.LayerCandidate) {
			conflicts = append(conflicts, &types.ConflictDetail{
				Type:            "land_use_mismatch",
				Severity:        types.SeverityHigh,
				Between:         []uint64{primary.LayerID, sec.LayerID},
				Description:     primary.LandUseName + " vs " + sec.LandUseName,
				Recommendation:  "Cần đối chiếu hồ sơ pháp lý gốc",
				ActionRequired:  "verify",
				AffectedPercent: sec.OverlapPercent,
			})
		}
	}
	return conflicts
}
