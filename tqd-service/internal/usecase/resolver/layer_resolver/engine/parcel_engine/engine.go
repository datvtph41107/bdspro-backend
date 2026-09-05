package parcel_engine

import (
	"fmt"
	"hash/crc32"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"tqd/internal/dto"
	"tqd/internal/usecase/resolver/layer_resolver/rule"
	"tqd/internal/usecase/resolver/layer_resolver/types"
)

type ParcelEngine struct {
	config               *EngineConfig
	legalNormalizer      *rule.LegalStatusResolver
	conflictClassifier   *rule.ConflictClassifier
	semanticNormalizer   *rule.LandUseCodeResolver
	compareEngine        *rule.CompareEngine
	historicalRiskEngine *rule.HistoricalRiskEngine
	lifecycleEngine      *rule.LifecycleEngine
	auditEngine          *rule.AuditEngine
	timelineEngine       *rule.TimelineEngine
}

type rankedLayerInternal struct {
	candidate  *types.LayerCandidate
	legalScore float64
	score      float64
	rank       int
}

func NewParcelEngine(config *EngineConfig) *ParcelEngine {
	if config == nil {
		config = DefaultEngineConfig()
	}
	return &ParcelEngine{
		config:               config,
		legalNormalizer:      rule.NewLegalStatusResolver(config.LegalNormalization),
		conflictClassifier:   rule.NewConflictClassifier(config.ConflictClassification),
		semanticNormalizer:   rule.MakeLandUseCodeResolver(config.SemanticCanonicalMap, config.SemanticLabels),
		compareEngine:        rule.NewCompareEngine(config.CompareConfig, nil),
		historicalRiskEngine: rule.NewHistoricalRiskEngine(config.HistoricalRiskConfig, nil, nil),
	}
}

func (e *ParcelEngine) WithLifecycleEngine(le *rule.LifecycleEngine) *ParcelEngine {
	e.lifecycleEngine = le
	return e
}

func (e *ParcelEngine) WithAuditEngine(ae *rule.AuditEngine) *ParcelEngine {
	e.auditEngine = ae
	return e
}

func (e *ParcelEngine) WithTimelineEngine(te *rule.TimelineEngine) *ParcelEngine {
	e.timelineEngine = te
	return e
}

func (e *ParcelEngine) WithCompareRepo(repo rule.CompareRepository) *ParcelEngine {
	e.compareEngine = rule.NewCompareEngine(e.config.CompareConfig, repo)
	return e
}

func (e *ParcelEngine) WithHistoricalRiskRepo(repo rule.HistoricalRiskRepository) *ParcelEngine {
	e.historicalRiskEngine = rule.NewHistoricalRiskEngine(e.config.HistoricalRiskConfig, repo, e.timelineEngine)
	return e
}

func (e *ParcelEngine) GetConfig() *EngineConfig {
	return e.config
}

// ============================================================
// Process — điểm vào duy nhất của engine
//
// mode: "popup" | "overview" | "detail" | "compare"
// ============================================================
func (e *ParcelEngine) Process(
	rows []UnifiedRow,
	parcelInfo *dto.ParcelInfoResponse,
	legalDocsMap map[uint64][]*dto.LayerLegalDTO,
	mode string,
) *UnifiedResult {
	startTime := time.Now()
	resolveMode := types.ResolveMode(mode)

	result := &UnifiedResult{
		Metadata: &MetadataInfo{
			DataSource: "postgis real-time (unified)",
			QueriedAt:  time.Now().UTC().Format(time.RFC3339),
		},
	}

	if len(rows) == 0 {
		result.ParcelInfo = parcelInfo
		if parcelInfo != nil {
			result.ParcelID = parcelInfo.ParcelID
			result.ParcelAreaSqm = parcelInfo.AreaSqm
		}
		result.Risk = e.assessRisk(nil, nil)
		result.ResolutionStatus = string(types.StatusResolved)
		result.Metadata.ResponseTimeMs = time.Since(startTime).Milliseconds()
		return result
	}

	parcelArea := rows[0].ParcelAreaSqm
	result.ParcelID = rows[0].ParcelID
	result.ParcelAreaSqm = parcelArea
	result.ParcelInfo = parcelInfo

	// Bước 1: Build candidates từ rows
	candidates := e.buildCandidates(rows)

	// Bước 2: Score & rank candidates
	ranked := e.scoreAndRank(candidates)

	// Normalize legal cho primary
	var primaryNormalized *types.NormalizedLegal
	if len(ranked) > 0 {
		primaryNormalized = e.legalNormalizer.Normalize(ranked[0].candidate)
	}

	switch resolveMode {
	case types.ModePopup:
		// Popup: chỉ cần primary plan + risk cơ bản
		if len(ranked) == 0 {
			result.Risk = e.assessRisk(nil, nil)
			result.ResolutionStatus = string(types.StatusResolved)
			break
		}
		result.PrimaryPlanning = e.buildPrimaryPlanning(ranked[0], parcelArea, primaryNormalized)
		result.Risk = e.assessRisk(ranked, nil)
		result.ResolutionStatus = string(types.StatusResolved)

	case types.ModeOverview:
		// Overview: primary + land use groups + layers + risk + warnings
		if len(ranked) == 0 {
			result.Risk = e.assessRisk(nil, nil)
			result.ResolutionStatus = string(types.StatusResolved)
			break
		}
		layers := e.buildLayers(rows, legalDocsMap, parcelArea)
		warnings := e.buildWarnings(parcelArea, layers, nil)
		result.PrimaryPlanning = e.buildPrimaryPlanning(ranked[0], parcelArea, primaryNormalized)
		result.LandUseGroups = e.buildLandUseGroups(ranked, parcelArea)
		result.Layers = layers
		result.Risk = e.assessRisk(ranked, nil)
		result.Warnings = warnings
		result.ResolutionStatus = e.resolveStatus(nil, warnings)

	case types.ModeCompare:
		// Compare: conflicts + risk + primary + groups + secondary + compare diff
		if len(ranked) == 0 {
			result.Risk = e.assessRisk(nil, nil)
			result.ResolutionStatus = string(types.StatusResolved)
			break
		}
		conflicts := e.detectConflicts(ranked)
		result.PrimaryPlanning = e.buildPrimaryPlanning(ranked[0], parcelArea, primaryNormalized)
		result.LandUseGroups = e.buildLandUseGroups(ranked, parcelArea)
		result.SecondaryPlans = e.buildSecondaryPlans(ranked[1:])
		result.Conflicts = conflicts
		result.Risk = e.assessRisk(ranked, conflicts)
		result.CompareResult = e.buildCompareResult(rows)
		result.ResolutionStatus = e.resolveStatus(conflicts, nil)

	default: // ModeDetail (và fallback)
		// Detail: toàn bộ các bước
		if len(ranked) == 0 {
			result.Risk = e.assessRisk(nil, nil)
			result.ResolutionStatus = string(types.StatusResolved)
			break
		}
		conflicts := e.detectConflicts(ranked)
		layers := e.buildLayers(rows, legalDocsMap, parcelArea)
		warnings := e.buildWarnings(parcelArea, layers, conflicts)
		result.PrimaryPlanning = e.buildPrimaryPlanning(ranked[0], parcelArea, primaryNormalized)
		result.LandUseGroups = e.buildLandUseGroups(ranked, parcelArea)
		result.SecondaryPlans = e.buildSecondaryPlans(ranked[1:])
		result.Layers = layers
		result.Conflicts = conflicts
		result.Risk = e.assessRisk(ranked, conflicts)
		result.Warnings = warnings
		result.HistoricalRisk = e.buildHistoricalRisk(rows, result.Risk)
		result.ResolutionStatus = e.resolveStatus(conflicts, warnings)
	}

	result.Metadata.ResponseTimeMs = time.Since(startTime).Milliseconds()
	return result
}

func (e *ParcelEngine) resolveStatus(conflicts []*ConflictInfo, warnings []string) string {
	if len(conflicts) > 0 {
		return string(types.StatusResolvedWithConflict)
	}
	if len(warnings) > 0 {
		return string(types.StatusResolvedWithWarning)
	}
	return string(types.StatusResolved)
}

// ============================================================
// STEP 1: BUILD CANDIDATES
// ============================================================
func (e *ParcelEngine) buildCandidates(rows []UnifiedRow) []*types.LayerCandidate {
	candidates := make([]*types.LayerCandidate, 0, len(rows))
	for _, row := range rows {
		// LandUseColor: đã COALESCE ở SQL (lu.color → label_color → group_color → #CCCCCC)
		// Không cần fallback thêm ở đây.
		landUseColor := row.LandUseColor
		if landUseColor == "" {
			landUseColor = row.LabelColor
		}
		if landUseColor == "" {
			landUseColor = row.GroupColor
		}

		c := &types.LayerCandidate{
			// Layer
			LayerID:          row.LayerID,
			LayerName:        row.LayerName,
			LayerDisplayName: row.LayerDisplayName,
			LayerType:        row.LayerType,
			LayerStatus:      row.LayerStatus,
			LegalStatus:      row.LayerLegalStatus,
			TrustValue:       row.LayerTrustValue,
			EffectiveDate:    parseTime(row.LayerEffectiveDate),
			ExpiryDate:       parseTime(row.LayerExpiryDate),

			// Label
			LabelID:          row.LabelID,
			LabelName:        row.LabelName,
			LabelDisplayName: row.LabelDisplayName,
			LabelColor:       row.LabelColor,

			// Land Use — PRIMARY từ qh_land_use
			LandUseCode:    row.LandUseCode,
			LandUseName:    row.LandUseName,
			LandUseColor:   landUseColor,
			GroupCode:      row.GroupCode,
			GroupName:      row.GroupName,
			GroupColor:     row.GroupColor,
			CanBuild:       row.CanBuild,
			Priority:       row.Priority,
			BuildCondition: row.BuildCondition,

			// Overlap
			OverlapAreaSqm: row.OverlapAreaSqm,
			OverlapPercent: row.OverlapPct,
			ParcelAreaSqm:  row.ParcelAreaSqm,

			// Spatial
			CenterLat: row.CenterLat,
			CenterLng: row.CenterLng,
			Geometry:  row.Geometry,
		}

		// Semantic normalization: chuẩn hóa LandUseCode theo canonical map
		e.semanticNormalizer.NormalizeCandidate(c)

		candidates = append(candidates, c)
	}
	return candidates
}

// ============================================================
// STEP 2: SCORE & RANK
// ============================================================
func (e *ParcelEngine) scoreAndRank(candidates []*types.LayerCandidate) []*rankedLayerInternal {
	ranked := make([]*rankedLayerInternal, len(candidates))
	for i, c := range candidates {
		normalized := e.legalNormalizer.Normalize(c)
		legalScore := normalized.Score
		priorityScore := float64(e.calculatePriority(c))
		areaScore := c.OverlapPercent

		totalScore := legalScore*e.config.LegalWeight +
			priorityScore*e.config.PriorityWeight +
			areaScore*e.config.AreaWeight

		ranked[i] = &rankedLayerInternal{
			candidate:  c,
			legalScore: legalScore,
			score:      totalScore,
		}
	}

	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].score == ranked[j].score {
			return ranked[i].candidate.OverlapPercent > ranked[j].candidate.OverlapPercent
		}
		return ranked[i].score > ranked[j].score
	})

	for i, r := range ranked {
		r.rank = i + 1
	}
	return ranked
}

func (e *ParcelEngine) calculatePriority(c *types.LayerCandidate) int {
	if c.Priority > 0 {
		return c.Priority
	}
	if p, ok := e.config.DefaultGroupPriority[c.GroupCode]; ok {
		return p
	}
	return e.config.DefaultPriority
}

// ============================================================
// STEP 3: DETECT CONFLICTS
// ============================================================
func (e *ParcelEngine) detectConflicts(ranked []*rankedLayerInternal) []*ConflictInfo {
	if len(ranked) < 2 {
		return nil
	}
	primary := ranked[0]
	var conflicts []*ConflictInfo
	for _, sec := range ranked[1:] {
		classified := e.conflictClassifier.Classify(primary.candidate, sec.candidate)
		if classified != nil {
			conflicts = append(conflicts, &ConflictInfo{
				Type:            string(classified.Type),
				Severity:        string(classified.Severity),
				Between:         classified.Between,
				Description:     classified.Description,
				Recommendation:  classified.Recommendation,
				ActionRequired:  classified.ActionRequired,
				AffectedPercent: classified.AffectedPercent,
			})
		}
	}
	return conflicts
}

// ============================================================
// STEP 4: ASSESS RISK
// ============================================================
func (e *ParcelEngine) assessRisk(ranked []*rankedLayerInternal, conflicts []*ConflictInfo) *RiskInfo {
	defaultThreshold := e.config.RiskThresholds[len(e.config.RiskThresholds)-1]
	risk := &RiskInfo{
		Level:          defaultThreshold.Level,
		Score:          0,
		CanBuild:       defaultThreshold.CanBuild,
		Reasons:        []string{},
		Recommendation: defaultThreshold.Recommendation,
	}

	if len(ranked) == 0 {
		risk.Reasons = append(risk.Reasons, e.config.NoPlanningImpactMsg)
		return risk
	}

	score := 0.0
	reasons := make([]string, 0, len(ranked))
	seenLayers := make(map[uint64]bool)

	for _, r := range ranked {
		c := r.candidate
		seenLayers[c.LayerID] = true

		rw := e.config.DefaultRiskWeight
		if cfg, ok := e.config.LandUseRiskWeights[c.LandUseCode]; ok {
			rw = cfg
		}
		score += c.OverlapPercent * rw.Multiplier
		reasons = append(reasons, fmt.Sprintf("%.1f%% diện tích thuộc %s - %s",
			c.OverlapPercent, c.LandUseName, rw.Label))
	}

	if len(seenLayers) > 1 {
		score += float64(len(seenLayers)) * e.config.MultiLayerPenalty
	}

	for _, c := range conflicts {
		switch c.Severity {
		case string(types.SeverityHigh):
			score += 15
		case string(types.SeverityMedium):
			score += 8
		case string(types.SeverityLow):
			score += 3
		}
	}
	if len(conflicts) > 0 {
		reasons = append(reasons, fmt.Sprintf("Phát hiện %d xung đột quy hoạch", len(conflicts)))
	}

	normalizedScore := uint32(math.Min(score, 100))

	for _, t := range e.config.RiskThresholds {
		if score >= t.MinScore {
			risk.Level = t.Level
			risk.Score = normalizedScore
			risk.CanBuild = t.CanBuild
			risk.Recommendation = t.Recommendation
			break
		}
	}

	if len(reasons) == 0 {
		reasons = append(reasons, e.config.NoSignificantImpactMsg)
	}
	risk.Reasons = reasons
	return risk
}

// ============================================================
// STEP 5: BUILD PRIMARY PLANNING
// ============================================================
func (e *ParcelEngine) buildPrimaryPlanning(
	primary *rankedLayerInternal,
	parcelArea float64,
	normalized *types.NormalizedLegal,
) *PrimaryPlanningInfo {
	c := primary.candidate
	ratio := 0.0
	if parcelArea > 0 {
		ratio = (c.OverlapAreaSqm / parcelArea) * 100
	}

	info := &PrimaryPlanningInfo{
		Code:        c.LandUseCode,
		Name:        c.LandUseName,
		Color:       c.LandUseColor,
		GroupCode:   c.GroupCode,
		GroupName:   c.GroupName,
		Ratio:       ratio,
		AreaSqm:     c.OverlapAreaSqm,
		LayerID:     c.LayerID,
		LayerName:   c.LayerDisplayName,
		LegalStatus: uint32(c.LegalStatus),
		CanBuild:    c.CanBuild,
	}

	if normalized != nil {
		info.LegalCanonical = string(normalized.Canonical)
		info.LegalConfidence = string(normalized.Confidence)
		info.LegalActiveState = string(normalized.ActiveState)
	}
	return info
}

func (e *ParcelEngine) buildSecondaryPlans(secondary []*rankedLayerInternal) []*SecondaryPlanInfo {
	if len(secondary) == 0 {
		return nil
	}
	plans := make([]*SecondaryPlanInfo, 0, len(secondary))
	for _, r := range secondary {
		c := r.candidate
		plans = append(plans, &SecondaryPlanInfo{
			LayerID:        c.LayerID,
			LayerName:      c.LayerDisplayName,
			LandUseCode:    c.LandUseCode,
			LandUseName:    c.LandUseName,
			LandUseColor:   c.LandUseColor,
			GroupCode:      c.GroupCode,
			OverlapPercent: c.OverlapPercent,
			OverlapAreaSqm: c.OverlapAreaSqm,
		})
	}
	return plans
}

// ============================================================
// STEP 6: BUILD LAND USE GROUPS
// ============================================================
func (e *ParcelEngine) buildLandUseGroups(ranked []*rankedLayerInternal, parcelArea float64) []*LandUseGroupInfo {
	groupMap := make(map[string]*layerGroupAccum)

	for _, r := range ranked {
		c := r.candidate
		if c.GroupCode == "" || c.GroupCode == "unknown" {
			continue
		}

		acc, ok := groupMap[c.GroupCode]
		if !ok {
			acc = &layerGroupAccum{
				code:     c.GroupCode,
				name:     c.GroupName,
				color:    c.GroupColor,
				canBuild: c.CanBuild,
				details:  make(map[string]*LandUseDetail),
			}
			groupMap[c.GroupCode] = acc
		}

		acc.ratio += c.OverlapPercent
		acc.areaSqm += c.OverlapAreaSqm

		if d, exists := acc.details[c.LandUseCode]; exists {
			d.Ratio += c.OverlapPercent
			d.AreaSqm += c.OverlapAreaSqm
		} else {
			// WarnLevel: lấy từ UnifiedRow thông qua candidate source
			// Candidate không có WarnLevel trực tiếp (không phải field của LayerCandidate)
			// → warnLevel sẽ được carry riêng từ UnifiedRow khi buildLayers
			acc.details[c.LandUseCode] = &LandUseDetail{
				LandUseCode: c.LandUseCode,
				LandUseName: c.LandUseName,
				Ratio:       c.OverlapPercent,
				AreaSqm:     c.OverlapAreaSqm,
			}
		}
	}

	result := make([]*LandUseGroupInfo, 0, len(groupMap))
	for _, acc := range groupMap {
		details := make([]*LandUseDetail, 0, len(acc.details))
		for _, d := range acc.details {
			details = append(details, d)
		}
		sort.Slice(details, func(i, j int) bool {
			return details[i].Ratio > details[j].Ratio
		})
		result = append(result, &LandUseGroupInfo{
			Code:     acc.code,
			Name:     acc.name,
			Color:    acc.color,
			CanBuild: acc.canBuild,
			Ratio:    acc.ratio,
			AreaSqm:  acc.areaSqm,
			Details:  details,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Ratio > result[j].Ratio
	})
	return result
}

// ============================================================
// STEP 7: BUILD LAYERS (detail mode)
// ============================================================
func (e *ParcelEngine) buildLayers(
	rows []UnifiedRow,
	legalDocsMap map[uint64][]*dto.LayerLegalDTO,
	parcelArea float64,
) []*LayerDetail {
	layerMap := make(map[uint64]*LayerDetail)
	// Preserve insertion order
	layerOrder := make([]uint64, 0)

	for _, r := range rows {
		layer, ok := layerMap[r.LayerID]
		if !ok {
			// Documents
			var docs []*DocumentInfo
			if legals, exists := legalDocsMap[r.LayerID]; exists {
				docs = make([]*DocumentInfo, 0, len(legals))
				for _, l := range legals {
					docs = append(docs, &DocumentInfo{
						ID:       l.ID,
						Name:     l.Name,
						FileUrl:  l.FileUrl,
						FileType: l.FileType,
					})
				}
			}

			// Authority
			var authority *AuthorityInfo
			if r.AuthorityID > 0 {
				authority = &AuthorityInfo{
					ID:          r.AuthorityID,
					Name:        r.AuthorityName,
					Code:        r.AuthorityCode,
					Description: r.AuthorityDesc,
				}
			}

			layer = &LayerDetail{
				ID:              r.LayerID,
				Name:            r.LayerName,
				DisplayName:     r.LayerDisplayName,
				Type:            r.LayerType,
				Status:          r.LayerStatus,
				LegalStatus:     r.LayerLegalStatus,
				LegalStatusName: getLegalStatusName(r.LayerLegalStatus),
				LayerAvatar:     r.LayerAvatar,
				EffectiveDate:   r.LayerEffectiveDate,
				ExpiryDate:      r.LayerExpiryDate,
				Authority:       authority,
				Documents:       docs,
				Labels:          []*LabelDetail{},
			}
			layerMap[r.LayerID] = layer
			layerOrder = append(layerOrder, r.LayerID)
		}

		// Group regions theo land use code dưới label
		labelKey := buildLabelKey(r.LayerID, r.LandUseCode)
		var label *LabelDetail
		for _, existing := range layer.Labels {
			if buildLabelKey(r.LayerID, existing.Code) == labelKey {
				label = existing
				break
			}
		}
		if label == nil {
			label = &LabelDetail{
				ID:      stableLabelID(r.LayerID, r.LandUseCode),
				Code:    r.LandUseCode,
				Name:    r.LandUseName,
				Color:   r.LandUseColor,
				Regions: []*RegionDetail{},
			}
			layer.Labels = append(layer.Labels, label)
		}

		// Region — WarnLevel và CanBuild từ UnifiedRow trực tiếp (không qua candidate)
		region := &RegionDetail{
			ID:             r.RegionID,
			Name:           r.RegionName,
			DisplayName:    r.RegionDisplayName,
			LandUseCode:    r.LandUseCode,
			LandUseName:    r.LandUseName,
			LandUseColor:   r.LandUseColor,
			OverlapAreaSqm: r.OverlapAreaSqm,
			OverlapPct:     r.OverlapPct,
			CenterLat:      r.CenterLat,
			CenterLng:      r.CenterLng,
			Geometry:       r.Geometry,
			WarnLevel:      r.WarnLevel, // lu.warn_level — align QuickLayers
			CanBuild:       r.CanBuild,  // lu.can_build — align QuickLayers
		}
		label.Regions = append(label.Regions, region)
		label.AreaSqm += r.OverlapAreaSqm
		layer.TotalArea += r.OverlapAreaSqm
	}

	// Sort và tính percentage
	layers := make([]*LayerDetail, 0, len(layerOrder))
	for _, id := range layerOrder {
		layer := layerMap[id]
		if parcelArea > 0 {
			layer.TotalPct = math.Min((layer.TotalArea/parcelArea)*100, 100)
		}
		for _, label := range layer.Labels {
			if parcelArea > 0 {
				label.Ratio = math.Min((label.AreaSqm/parcelArea)*100, 100)
			}
			sort.Slice(label.Regions, func(i, j int) bool {
				return label.Regions[i].OverlapAreaSqm > label.Regions[j].OverlapAreaSqm
			})
		}
		sort.Slice(layer.Labels, func(i, j int) bool {
			return layer.Labels[i].AreaSqm > layer.Labels[j].AreaSqm
		})
		layers = append(layers, layer)
	}

	sort.Slice(layers, func(i, j int) bool {
		return layers[i].TotalArea > layers[j].TotalArea
	})
	return layers
}

func getLegalStatusName(status uint32) string {
	switch status {
	case 10:
		return "Đã ban hành"
	case 20:
		return "Đang hiệu lực"
	case 30:
		return "Hết hiệu lực"
	case 40:
		return "Đang điều chỉnh"
	default:
		return fmt.Sprintf("Không xác định (%d)", status)
	}
}

// ============================================================
// STEP 8: BUILD WARNINGS
// ============================================================
func (e *ParcelEngine) buildWarnings(
	parcelArea float64,
	layers []*LayerDetail,
	conflicts []*ConflictInfo,
) []string {
	var warnings []string
	if len(layers) == 0 {
		return warnings
	}

	totalOverlap := 0.0
	for _, l := range layers {
		totalOverlap += l.TotalArea
	}
	if parcelArea > 0 && totalOverlap > parcelArea {
		warnings = append(warnings, e.config.WarningOverlapExceeded)
	}
	if len(layers) > 1 {
		warnings = append(warnings, e.config.WarningMultiLayer)
	}
	if len(conflicts) > 0 {
		warnings = append(warnings, e.config.WarningConflict)
	}
	return warnings
}

// ============================================================
// COMPARE & HISTORICAL RISK
// ============================================================
func (e *ParcelEngine) buildCompareResult(rows []UnifiedRow) *CompareResultInfo {
	type layerSnap struct {
		layerID   uint64
		layerName string
		landUse   string
		legal     uint32
	}

	seen := make(map[uint64]*layerSnap)
	var changes []*CompareChangeInfo

	for _, row := range rows {
		snap := &layerSnap{
			layerID:   row.LayerID,
			layerName: row.LayerDisplayName,
			landUse:   row.LandUseName,
			legal:     row.LayerLegalStatus,
		}
		if prev, exists := seen[row.LayerID]; exists {
			if prev.landUse != snap.landUse || prev.legal != snap.legal {
				changes = append(changes, &CompareChangeInfo{
					LayerID:    row.LayerID,
					LayerName:  row.LayerDisplayName,
					ChangeType: "modified",
					OldValue:   fmt.Sprintf("%s (legal:%d)", prev.landUse, prev.legal),
					NewValue:   fmt.Sprintf("%s (legal:%d)", snap.landUse, snap.legal),
				})
			}
		}
		seen[row.LayerID] = snap
	}

	return &CompareResultInfo{
		TotalLayers: len(seen),
		Changes:     changes,
	}
}

func (e *ParcelEngine) buildHistoricalRisk(rows []UnifiedRow, risk *RiskInfo) *HistoricalRiskInfo {
	if risk == nil {
		return nil
	}
	now := time.Now()
	activeCount := 0
	expiringCount := 0

	for _, row := range rows {
		if row.LayerEffectiveDate != "" {
			if effDate := parseTime(row.LayerEffectiveDate); effDate != nil && !effDate.After(now) {
				activeCount++
			}
		}
		if row.LayerExpiryDate != "" {
			if expDate := parseTime(row.LayerExpiryDate); expDate != nil &&
				expDate.Before(now.AddDate(0, 6, 0)) && expDate.After(now) {
				expiringCount++
			}
		}
	}

	trend := "stable"
	trendDetail := "Rủi ro ổn định"
	if expiringCount > 0 {
		trend = "worsening"
		trendDetail = fmt.Sprintf("%d layer sắp hết hạn trong 6 tháng tới", expiringCount)
	}

	return &HistoricalRiskInfo{
		Trend:          trend,
		TrendDetail:    trendDetail,
		ActiveLayers:   activeCount,
		ExpiringLayers: expiringCount,
		AnalyzedAt:     now.Format(time.RFC3339),
	}
}

// ============================================================
// HELPERS
// ============================================================

func buildLabelKey(layerID uint64, landUseCode string) string {
	return fmt.Sprintf("%d:%s", layerID, strings.TrimSpace(landUseCode))
}

func stableLabelID(layerID uint64, landUseCode string) uint64 {
	key := buildLabelKey(layerID, landUseCode)
	return uint64(crc32.ChecksumIEEE([]byte(key)))
}

// hexToRGB parse "#RRGGBB" → (r, g, b uint8, ok bool)
func hexToRGB(hex string) (r, g, b uint8, ok bool) {
	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(hex) != 6 {
		return 0, 0, 0, false
	}
	v, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return uint8((v >> 16) & 0xFF), uint8((v >> 8) & 0xFF), uint8(v & 0xFF), true
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}
