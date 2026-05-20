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
	"tqd/internal/enums"
	"tqd/internal/modules/resolver/layer_resolver/types"
)

type ParcelEngine struct {
	config *EngineConfig
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
	return &ParcelEngine{config: config}
}

func (e *ParcelEngine) GetConfig() *EngineConfig {
	return e.config
}

// Process - xử lý toàn bộ nghiệp vụ, trả về UnifiedResult dùng cho cả V1 và V2
func (e *ParcelEngine) Process(rows []UnifiedRow, parcelInfo *dto.ParcelInfoResponse, legalDocsMap map[uint64][]*dto.LayerLegalDTO, mode string) *UnifiedResult {
	startTime := time.Now()

	result := &UnifiedResult{
		Metadata: &MetadataInfo{
			DataSource: "postgis real-time",
			QueriedAt:  time.Now().UTC().Format(time.RFC3339),
		},
	}

	if len(rows) == 0 {
		result.Metadata.ResponseTimeMs = time.Since(startTime).Milliseconds()
		return result
	}

	parcelArea := rows[0].ParcelAreaSqm
	result.ParcelID = rows[0].ParcelID
	result.ParcelAreaSqm = parcelArea
	result.ParcelInfo = parcelInfo

	// Bước 1: Build candidates từ rows
	candidates := e.buildCandidates(rows)

	// Bước 2: Score & Rank
	ranked := e.scoreAndRank(candidates)

	// Bước 3: Detect conflicts
	conflicts := e.detectConflicts(ranked)

	// Bước 4: Assess risk (single source of truth)
	risk := e.assessRisk(ranked, conflicts)

	// Bước 5: Build primary planning
	primaryPlanning := e.buildPrimaryPlanning(ranked[0], parcelArea)

	// Bước 6: Build land use groups
	landUseGroups := e.buildLandUseGroups(ranked, parcelArea)

	// Bước 7: Build raw layers (cho V1)
	layers := e.buildLayers(rows, legalDocsMap, parcelArea)

	// Bước 8: Build warnings
	warnings := e.buildWarnings(parcelArea, layers, conflicts)

	result.PrimaryPlanning = primaryPlanning
	result.LandUseGroups = landUseGroups
	result.Layers = layers
	result.Conflicts = conflicts
	result.Risk = risk
	result.Warnings = warnings
	result.Metadata.ResponseTimeMs = time.Since(startTime).Milliseconds()

	return result
}

// ============================================
// STEP 1: BUILD CANDIDATES
// ============================================

func (e *ParcelEngine) buildCandidates(rows []UnifiedRow) []*types.LayerCandidate {
	candidates := make([]*types.LayerCandidate, 0, len(rows))
	for _, row := range rows {
		c := &types.LayerCandidate{
			LayerID:          row.LayerID,
			LayerName:        row.LayerName,
			LayerDisplayName: row.LayerDisplayName,
			LayerType:        row.LayerType,
			LayerStatus:      row.LayerStatus,
			LegalStatus:      row.LayerLegalStatus,
			TrustValue:       row.LayerTrustValue,
			EffectiveDate:    parseTime(row.LayerEffectiveDate),
			ExpiryDate:       parseTime(row.LayerExpiryDate),
			Authority:        row.LayerAuthority,
			LabelID:          row.LabelID,
			LabelName:        row.LabelName,
			LabelDisplayName: row.LabelDisplayName,
			LabelColor:       row.LabelColor,
			LandUseCode:      row.RegionLandUseCode,
			LandUseName:      row.RegionLandUseName,
			LandUseColor:     row.LabelColor,
			GroupCode:        enums.GroupCode(row.LabelGroupCode).String(),
			GroupName:        row.GroupName,
			GroupColor:       row.GroupColor,
			CanBuild:         row.LabelCanBuild,
			Priority:         row.LabelPriority,
			BuildCondition:   row.LabelBuildCondition,
			OverlapAreaSqm:   row.OverlapAreaSqm,
			OverlapPercent:   row.OverlapPct,
			ParcelAreaSqm:    row.ParcelAreaSqm,
			CenterLat:        row.CenterLat,
			CenterLng:        row.CenterLng,
			Geometry:         row.Geometry,
		}
		candidates = append(candidates, c)
	}
	return candidates
}

// ============================================
// STEP 2: SCORE & RANK
// ============================================

func (e *ParcelEngine) scoreAndRank(candidates []*types.LayerCandidate) []*rankedLayerInternal {
	ranked := make([]*rankedLayerInternal, len(candidates))
	for i, c := range candidates {
		legalScore := e.config.LegalScores[c.LegalStatus]
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

// ============================================
// STEP 3: DETECT CONFLICTS
// ============================================

func (e *ParcelEngine) detectConflicts(ranked []*rankedLayerInternal) []*ConflictInfo {
	if len(ranked) < 2 {
		return nil
	}

	primary := ranked[0]
	var conflicts []*ConflictInfo

	for _, sec := range ranked[1:] {
		if e.hasConflict(primary.candidate, sec.candidate) {
			conflicts = append(conflicts, &ConflictInfo{
				Type:            e.config.ConflictType,
				Severity:        e.config.ConflictSeverity,
				Between:         []uint64{primary.candidate.LayerID, sec.candidate.LayerID},
				Description:     primary.candidate.LandUseName + " vs " + sec.candidate.LandUseName,
				Recommendation:  e.config.ConflictRecommendation,
				ActionRequired:  e.config.ConflictActionRequired,
				AffectedPercent: sec.candidate.OverlapPercent,
			})
		}
	}
	return conflicts
}

func (e *ParcelEngine) hasConflict(a, b *types.LayerCandidate) bool {
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

// ============================================
// STEP 4: ASSESS RISK (SINGLE SOURCE OF TRUTH)
// ============================================

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
		pct := r.candidate.OverlapPercent
		code := r.candidate.LandUseCode
		name := r.candidate.LandUseName
		seenLayers[r.candidate.LayerID] = true

		rw := e.config.DefaultRiskWeight
		if cfg, ok := e.config.LandUseRiskWeights[code]; ok {
			rw = cfg
		}

		score += pct * rw.Multiplier
		reasons = append(reasons, fmt.Sprintf("%.1f%% diện tích thuộc %s - %s", pct, name, rw.Label))
	}

	if len(seenLayers) > 1 {
		score += float64(len(seenLayers)) * e.config.MultiLayerPenalty
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

// ============================================
// STEP 5: BUILD PRIMARY PLANNING
// ============================================

func (e *ParcelEngine) buildPrimaryPlanning(primary *rankedLayerInternal, parcelArea float64) *PrimaryPlanningInfo {
	c := primary.candidate
	ratio := 0.0
	if parcelArea > 0 {
		ratio = (c.OverlapAreaSqm / parcelArea) * 100
	}

	return &PrimaryPlanningInfo{
		Code:        c.LandUseCode,
		Name:        c.LandUseName,
		Color:       c.LandUseColor,
		GroupCode:   c.GroupCode,
		GroupName:   c.GroupName,
		Ratio:       ratio,
		AreaSqm:     c.OverlapAreaSqm,
		LayerID:     c.LayerID,
		LayerName:   c.LayerDisplayName,
		LegalStatus: c.LegalStatus,
		CanBuild:    c.CanBuild,
	}
}

// ============================================
// STEP 6: BUILD LAND USE GROUPS
// ============================================

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

		detailKey := c.LandUseCode
		if d, exists := acc.details[detailKey]; exists {
			d.Ratio += c.OverlapPercent
			d.AreaSqm += c.OverlapAreaSqm
		} else {
			acc.details[detailKey] = &LandUseDetail{
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

// ============================================
// STEP 7: BUILD RAW LAYERS (CHO V1)
// ============================================

func (e *ParcelEngine) buildLayers(rows []UnifiedRow, legalDocsMap map[uint64][]*dto.LayerLegalDTO, parcelArea float64) []*LayerDetail {
	layerMap := make(map[uint64]*LayerDetail)

	for _, r := range rows {
		layer, ok := layerMap[r.LayerID]
		if !ok {
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

			var authority *AuthorityInfo
			if r.AuthorityIssuringID > 0 {
				authority = &AuthorityInfo{
					ID:          r.AuthorityIssuringID,
					Name:        r.AuthorityIssuringName,
					Code:        r.AuthorityIssuringCode,
					Description: r.AuthorityIssuringDesc,
				}
			}

			layer = &LayerDetail{
				ID:            r.LayerID,
				Name:          r.LayerName,
				DisplayName:   r.LayerDisplayName,
				Type:          r.LayerType,
				Status:        r.LayerStatus,
				LegalStatus:   r.LayerLegalStatus,
				EffectiveDate: r.LayerEffectiveDate,
				ExpiryDate:    r.LayerExpiryDate,
				Authority:     authority,
				Documents:     docs,
				Labels:        []*LabelDetail{},
				Geometry:      r.Geometry,
			}
			layerMap[r.LayerID] = layer
		}

		labelKey := buildLabelKey(r.LayerID, r.RegionLandUseCode)
		var label *LabelDetail
		for _, existing := range layer.Labels {
			if buildLabelKey(r.LayerID, existing.Code) == labelKey {
				label = existing
				break
			}
		}

		if label == nil {
			label = &LabelDetail{
				ID:      stableLabelID(r.LayerID, r.RegionLandUseCode),
				Code:    r.RegionLandUseCode,
				Name:    r.RegionLandUseName,
				Color:   r.RegionLandUseColor,
				Regions: []*RegionDetail{},
			}
			layer.Labels = append(layer.Labels, label)
		}

		region := &RegionDetail{
			ID:             r.RegionID,
			Name:           r.RegionName,
			DisplayName:    r.RegionDisplayName,
			LandUseCode:    r.RegionLandUseCode,
			LandUseName:    r.RegionLandUseName,
			Color:          regionRGB(r),
			OverlapAreaSqm: r.OverlapAreaSqm,
			OverlapPct:     r.OverlapPct,
			CenterLat:      r.CenterLat,
			CenterLng:      r.CenterLng,
			Geometry:       r.Geometry,
		}
		label.Regions = append(label.Regions, region)
		label.AreaSqm += r.OverlapAreaSqm
		layer.TotalArea += r.OverlapAreaSqm
	}

	layers := make([]*LayerDetail, 0, len(layerMap))
	for _, layer := range layerMap {
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

// ============================================
// STEP 8: BUILD WARNINGS
// ============================================

func (e *ParcelEngine) buildWarnings(parcelArea float64, layers []*LayerDetail, conflicts []*ConflictInfo) []string {
	var warnings []string

	if len(layers) == 0 {
		return warnings
	}

	totalOverlapRaw := 0.0
	for _, l := range layers {
		totalOverlapRaw += l.TotalArea
	}

	if parcelArea > 0 && totalOverlapRaw > parcelArea {
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

// ============================================
// HELPER FUNCTIONS
// ============================================

func buildLabelKey(layerID uint64, landUseCode string) string {
	return fmt.Sprintf("%d:%s", layerID, strings.TrimSpace(landUseCode))
}

func stableLabelID(layerID uint64, landUseCode string) uint64 {
	key := buildLabelKey(layerID, landUseCode)
	return uint64(crc32.ChecksumIEEE([]byte(key)))
}

func regionRGB(r UnifiedRow) *dto.RGB {
	if r.RegionLandUseColor != "" {
		if rgb, ok := hexToRGB(r.RegionLandUseColor); ok {
			return rgb
		}
	}
	return &dto.RGB{
		Red:   r.RegionColorRed,
		Green: r.RegionColorGreen,
		Blue:  r.RegionColorBlue,
	}
}

func hexToRGB(hex string) (*dto.RGB, bool) {
	hex = strings.TrimPrefix(strings.TrimSpace(hex), "#")
	if len(hex) != 6 {
		return nil, false
	}
	v, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return nil, false
	}
	return &dto.RGB{
		Red:   uint8((v >> 16) & 0xFF),
		Green: uint8((v >> 8) & 0xFF),
		Blue:  uint8(v & 0xFF),
	}, true
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
