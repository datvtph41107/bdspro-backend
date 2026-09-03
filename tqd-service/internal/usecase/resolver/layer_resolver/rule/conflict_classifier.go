package rule

import (
	"fmt"
	"math"
	"strings"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// ConflictClassifier — Engine phân loại conflict 6 loại × 3 mức với confidence score
type ConflictClassifier struct {
	config *types.ConflictClassificationConfig
}

// NewConflictClassifier tạo ConflictClassifier với config, nếu config nil thì dùng default
func NewConflictClassifier(config *types.ConflictClassificationConfig) *ConflictClassifier {
	if config == nil {
		config = &types.ConflictClassificationConfig{
			SeverityThresholds: types.DefaultSeverityThresholds(),
			TypeWeights:        types.DefaultConflictTypeWeights(),
			GuidanceTemplates:  types.DefaultGuidanceTemplates(),
		}
	}
	if config.SeverityThresholds.High == 0 {
		config.SeverityThresholds = types.DefaultSeverityThresholds()
	}
	if config.TypeWeights == nil {
		config.TypeWeights = types.DefaultConflictTypeWeights()
	}
	if config.GuidanceTemplates == nil {
		config.GuidanceTemplates = types.DefaultGuidanceTemplates()
	}
	return &ConflictClassifier{config: config}
}

// Name trả về tên engine
func (c *ConflictClassifier) Name() string { return "ConflictClassifier" }

// Validate kiểm tra config hợp lệ
func (c *ConflictClassifier) Validate() error {
	if c.config.SeverityThresholds.High <= c.config.SeverityThresholds.Medium {
		return fmt.Errorf("ConflictClassifier: high threshold must be > medium threshold")
	}
	return nil
}

// Classify phân loại conflict giữa 2 layer candidate
// Trả về nil nếu không có conflict
func (c *ConflictClassifier) Classify(a, b *types.LayerCandidate) *types.ClassifiedConflict {
	// B1: Xác định các loại conflict tồn tại
	conflictTypes := c.detectConflictTypes(a, b)
	if len(conflictTypes) == 0 {
		return nil
	}

	// B2: Chọn loại conflict (mixed nếu >1)
	conflictType := c.resolveType(conflictTypes)

	// B3: Tính severity dựa trên overlap % và impact
	severity := c.calculateSeverity(conflictType, a, b)

	// B4: Tính confidence
	confidence := c.calculateConfidence(conflictType, a, b)

	// B5: Sinh description và recommendation từ template
	desc, rec, action := c.generateGuidance(conflictType, severity, a, b)

	return &types.ClassifiedConflict{
		Type:            conflictType,
		Severity:        severity,
		Confidence:      confidence,
		Between:         []uint64{a.LayerID, b.LayerID},
		Description:     desc,
		ReasonCode:      c.generateReasonCode(conflictTypes),
		Recommendation:  rec,
		ActionRequired:  action,
		AffectedPercent: math.Min(a.OverlapPercent, b.OverlapPercent),
	}
}

// detectConflictTypes xác định tất cả loại conflict giữa 2 layer
func (c *ConflictClassifier) detectConflictTypes(a, b *types.LayerCandidate) []types.ConflictType {
	var conflictTypes []types.ConflictType

	// 1. Semantic: khác GroupCode
	if a.GroupCode != "" && b.GroupCode != "" && a.GroupCode != b.GroupCode {
		conflictTypes = append(conflictTypes, types.ConflictSemantic)
	}

	// 2. Legal: khác LegalStatus
	if a.LegalStatus != b.LegalStatus {
		conflictTypes = append(conflictTypes, types.ConflictLegal)
	}

	// 3. Temporal: effective date ranges không overlap
	if c.isTemporalConflict(a, b) {
		conflictTypes = append(conflictTypes, types.ConflictTemporal)
	}

	// 4. Source: khác SourceType
	if a.SourceType != "" && b.SourceType != "" && a.SourceType != b.SourceType {
		conflictTypes = append(conflictTypes, types.ConflictSource)
	}

	// 5. Geometry: overlap area > 0 nhưng không full containment
	if a.OverlapAreaSqm > 0 && b.OverlapAreaSqm > 0 {
		conflictTypes = append(conflictTypes, types.ConflictGeometry)
	}

	return conflictTypes
}

// isTemporalConflict kiểm tra xem 2 layer có xung đột thời gian không
func (c *ConflictClassifier) isTemporalConflict(a, b *types.LayerCandidate) bool {
	if a.EffectiveDate == nil || b.EffectiveDate == nil {
		return false
	}
	// 2 layer có effective date cách nhau > 5 năm → temporal conflict
	diff := a.EffectiveDate.Sub(*b.EffectiveDate)
	if diff < 0 {
		diff = -diff
	}
	return diff.Hours() > 24*365*5 // 5 years
}

// resolveType chọn loại conflict (mixed nếu >1)
func (c *ConflictClassifier) resolveType(conflictTypes []types.ConflictType) types.ConflictType {
	if len(conflictTypes) == 1 {
		return conflictTypes[0]
	}
	return types.ConflictMixed
}

// calculateSeverity tính severity dựa trên overlap % và conflict type weight
func (c *ConflictClassifier) calculateSeverity(ct types.ConflictType, a, b *types.LayerCandidate) types.ConflictSeverity {
	overlapPct := math.Min(a.OverlapPercent, b.OverlapPercent)

	// Base score từ overlap
	baseScore := overlapPct

	// Bonus cho legal conflict (khác legal status nghiêm trọng)
	if ct == types.ConflictLegal || ct == types.ConflictMixed {
		legalDiff := math.Abs(float64(a.LegalStatus) - float64(b.LegalStatus))
		baseScore += legalDiff * 0.3
	}

	// Bonus cho semantic conflict
	if ct == types.ConflictSemantic {
		baseScore += 20
	}

	thresholds := c.config.SeverityThresholds
	switch {
	case baseScore >= thresholds.High:
		return types.SeverityHigh
	case baseScore >= thresholds.Medium:
		return types.SeverityMedium
	default:
		return types.SeverityLow
	}
}

// calculateConfidence tính confidence score (0.0 - 1.0)
func (c *ConflictClassifier) calculateConfidence(ct types.ConflictType, a, b *types.LayerCandidate) float64 {
	confidence := 0.5 // base

	// Tăng confidence nếu overlap rõ ràng
	if a.OverlapPercent > 10 && b.OverlapPercent > 10 {
		confidence += 0.2
	}

	// Tăng confidence nếu cả 2 layer có trust value cao
	if a.TrustValue > 0.7 && b.TrustValue > 0.7 {
		confidence += 0.15
	}

	// Tăng confidence nếu cả 2 có legal status rõ ràng
	if a.LegalStatus > 0 && b.LegalStatus > 0 {
		confidence += 0.1
	}

	// Giảm confidence nếu 1 trong 2 có trust value thấp
	if a.TrustValue < 0.3 || b.TrustValue < 0.3 {
		confidence -= 0.15
	}

	return math.Max(0.0, math.Min(1.0, confidence))
}

// generateGuidance sinh description, recommendation, action từ template
func (c *ConflictClassifier) generateGuidance(
	ct types.ConflictType, severity types.ConflictSeverity, a, b *types.LayerCandidate,
) (desc, rec, action string) {
	tmpl, ok := c.config.GuidanceTemplates[ct]
	if !ok {
		tmpl = &types.GuidanceTemplate{
			DescriptionTemplate: "{a_layer} xung đột với {b_layer}",
			Recommendation:      "Cần kiểm tra thêm",
			ActionRequired:      "verify",
		}
	}

	desc = tmpl.DescriptionTemplate
	desc = strings.ReplaceAll(desc, "{a_landuse}", a.LandUseName)
	desc = strings.ReplaceAll(desc, "{a_layer}", a.LayerDisplayName)
	desc = strings.ReplaceAll(desc, "{b_landuse}", b.LandUseName)
	desc = strings.ReplaceAll(desc, "{b_layer}", b.LayerDisplayName)
	desc = strings.ReplaceAll(desc, "{a_legal}", fmt.Sprintf("%d", a.LegalStatus))
	desc = strings.ReplaceAll(desc, "{b_legal}", fmt.Sprintf("%d", b.LegalStatus))
	desc = strings.ReplaceAll(desc, "{a_source}", a.SourceType)
	desc = strings.ReplaceAll(desc, "{b_source}", b.SourceType)

	if a.EffectiveDate != nil {
		desc = strings.ReplaceAll(desc, "{a_date}", a.EffectiveDate.Format("02/01/2006"))
	}
	if b.EffectiveDate != nil {
		desc = strings.ReplaceAll(desc, "{b_date}", b.EffectiveDate.Format("02/01/2006"))
	}

	rec = tmpl.Recommendation
	action = tmpl.ActionRequired

	return
}

// generateReasonCode sinh mã lý do từ danh sách conflict types
func (c *ConflictClassifier) generateReasonCode(types []types.ConflictType) string {
	parts := make([]string, len(types))
	for i, t := range types {
		parts[i] = string(t)
	}
	return strings.Join(parts, "+")
}
