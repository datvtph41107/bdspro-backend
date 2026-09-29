package mapper

import (
	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// MapToResolveResult chuyển đổi UnifiedResult → types.ResolveResult
// Đây là mapper trung gian giữa engine output và response types
func MapToResolveResult(
	parcelID uint64,
	parcelAreaSqm float64,
	primary *types.RankedLayer,
	secondary []*types.RankedLayer,
	conflicts []*types.ConflictDetail,
	riskScore uint32,
	riskLevel string,
	riskReasons []string,
	canBuild bool,
	buildStatus string,
	reason string,
	recommendation string,
	mode types.ResolveMode,
	processingTimeMs int64,
	layers []*types.LayerDetail,
	layerLandUseGroups map[uint64][]*types.LandUseGroupDetail,
) *types.ResolveResult {
	hasConflict := len(conflicts) > 0

	var criticalConflicts []*types.ConflictDetail
	var infoConflicts []*types.ConflictDetail
	for _, c := range conflicts {
		if c.Severity == types.SeverityHigh {
			criticalConflicts = append(criticalConflicts, c)
		} else {
			infoConflicts = append(infoConflicts, c)
		}
	}

	return &types.ResolveResult{
		Primary:            primary,
		Secondary:          secondary,
		HasConflict:        hasConflict,
		CriticalConflicts:  criticalConflicts,
		InfoConflicts:      infoConflicts,
		CanBuild:           canBuild,
		BuildStatus:        buildStatus,
		Reason:             reason,
		Mode:               mode,
		ProcessingTimeMs:   processingTimeMs,
		LayerLandUseGroups: layerLandUseGroups,
		ParcelAreaSqm:      parcelAreaSqm,
		RiskScore:          riskScore,
		RiskLevel:          riskLevel,
		RiskReasons:        riskReasons,
		Recommendation:     recommendation,
		Layers:             layers,
	}
}

// MapResolveResultToLegacy chuyển ResolveResult → map[string]any cho V1 API compatibility
func MapResolveResultToLegacy(result *types.ResolveResult) map[string]any {
	m := map[string]any{
		"has_conflict":   result.HasConflict,
		"can_build":      result.CanBuild,
		"build_status":   result.BuildStatus,
		"reason":         result.Reason,
		"risk_score":     result.RiskScore,
		"risk_level":     result.RiskLevel,
		"recommendation": result.Recommendation,
	}

	if result.Primary != nil {
		m["primary_plan"] = map[string]any{
			"layer_id":        result.Primary.LayerID,
			"layer_name":      result.Primary.LayerDisplayName,
			"land_use_code":   result.Primary.LandUseCode,
			"land_use_name":   result.Primary.LandUseName,
			"overlap_percent": result.Primary.OverlapPercent,
			"can_build":       result.Primary.CanBuild,
		}
	}

	if len(result.CriticalConflicts) > 0 {
		conflictList := make([]map[string]any, len(result.CriticalConflicts))
		for i, c := range result.CriticalConflicts {
			conflictList[i] = map[string]any{
				"type":        c.Type,
				"severity":    c.Severity,
				"description": c.Description,
			}
		}
		m["conflicts"] = conflictList
	}

	return m
}
