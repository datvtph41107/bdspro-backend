package rule

import (
	"fmt"
	"math"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

type AssessEngine struct {
	config *types.Config
}

func NewAssessEngine(config *types.Config) *AssessEngine {
	return &AssessEngine{config: config}
}

func (e *AssessEngine) Assess(primary *types.RankedLayer, secondary []*types.RankedLayer, hasConflict bool) (canBuild bool, reason string, status string, riskScore uint32, riskLevel string, riskReasons []string, recommendation string) {
	allLayers := append([]*types.RankedLayer{primary}, secondary...)

	score := 0.0
	reasons := make([]string, 0, len(allLayers))

	for _, layer := range allLayers {
		pct := layer.OverlapPercent
		code := layer.LandUseCode
		name := layer.LandUseName

		rw := e.config.DefaultRiskWeight
		if cfg, ok := e.config.LandUseRiskWeights[code]; ok {
			rw = cfg
		}

		score += pct * rw.Multiplier
		reasons = append(reasons, fmt.Sprintf("%.1f%% diện tích thuộc %s - %s", pct, name, rw.Label))
	}

	if len(allLayers) > 1 {
		uniqueLayers := make(map[uint64]bool)
		for _, l := range allLayers {
			uniqueLayers[l.LayerID] = true
		}
		if len(uniqueLayers) > 1 {
			score += float64(len(uniqueLayers)) * e.config.MultiLayerPenalty
		}
	}

	normalizedScore := uint32(math.Min(score, 100))

	for _, t := range e.config.RiskThresholds {
		if score >= t.MinScore {
			riskLevel = t.Level
			canBuild = t.CanBuild
			recommendation = t.Recommendation
			break
		}
	}

	if len(reasons) == 0 {
		reasons = append(reasons, e.config.NoSignificantImpactMsg)
	}

	if !primary.CanBuild {
		status = e.config.StatusProhibited
		reason = primary.BuildCondition
	} else if hasConflict {
		status = e.config.StatusConditional
		reason = primary.BuildCondition + e.config.ConflictStatusSuffix
	} else {
		status = e.config.StatusAllowed
		reason = primary.BuildCondition
	}

	return canBuild, reason, status, normalizedScore, riskLevel, reasons, recommendation
}

func (e *AssessEngine) calculateParcelArea(primary *types.RankedLayer) float64 {
	if primary != nil {
		return primary.ParcelAreaSqm
	}
	return 0
}
