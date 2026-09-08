package rule

import (
	"fmt"
	"time"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// LegalStatusResolver — Engine chuẩn hóa trạng thái pháp lý
// Pipeline 8 bước: raw → canonical → active → confidence → warnings → reason codes → score → trace
type LegalStatusResolver struct {
	config       *types.LegalNormalizationConfig
	canonicalMap map[uint32]types.CanonicalLegalStatus
}

// NewLegalStatusResolver tạo LegalStatusResolver với config, nếu config nil thì dùng default
func NewLegalStatusResolver(config *types.LegalNormalizationConfig) *LegalStatusResolver {
	if config == nil {
		config = &types.LegalNormalizationConfig{
			CanonicalMap:     types.DefaultCanonicalMap(),
			ConfidenceWeight: types.DefaultConfidenceWeight(),
			ActiveRules:      types.DefaultActiveRules(),
		}
	}
	if config.CanonicalMap == nil {
		config.CanonicalMap = types.DefaultCanonicalMap()
	}
	return &LegalStatusResolver{
		config:       config,
		canonicalMap: config.CanonicalMap,
	}
}

// Name trả về tên engine
func (n *LegalStatusResolver) Name() string { return "LegalStatusResolver" }

// Validate kiểm tra config hợp lệ
func (n *LegalStatusResolver) Validate() error {
	if len(n.canonicalMap) == 0 {
		return fmt.Errorf("LegalStatusResolver: canonical_map is empty")
	}
	return nil
}

// Normalize thực hiện pipeline 8 bước chuẩn hóa
func (n *LegalStatusResolver) Normalize(candidate *types.LayerCandidate) *types.NormalizedLegal {
	result := &types.NormalizedLegal{
		RawStatus:    candidate.LegalStatus,
		TracePayload: make(map[string]any),
	}

	// B1: Map raw → canonical taxonomy
	result.Canonical = n.mapToCanonical(candidate.LegalStatus)

	// B2: Đánh giá active state dựa trên EffectiveDate, ExpiryDate
	result.ActiveState = n.evaluateActiveState(candidate)

	// B3: Tính confidence dựa trên độ đầy đủ của metadata
	result.Confidence = n.calculateConfidence(candidate)

	// B4: Sinh warning flags nếu có ambiguity
	result.WarningFlags = n.generateWarnings(result)

	// B5: Sinh reason codes
	result.ReasonCodes = n.generateReasonCodes(candidate, result)

	// B6: Tính legal score (weighted by confidence)
	result.Score = n.calculateScore(result)

	// B7: Build trace payload
	result.TracePayload = n.buildTrace(candidate, result)

	return result
}

// mapToCanonical ánh xạ raw status → canonical taxonomy
func (n *LegalStatusResolver) mapToCanonical(raw uint32) types.CanonicalLegalStatus {
	if canonical, ok := n.canonicalMap[raw]; ok {
		return canonical
	}
	return types.LegalUnknown
}

// evaluateActiveState đánh giá active state dựa trên ngày hiệu lực và hết hạn
func (n *LegalStatusResolver) evaluateActiveState(c *types.LayerCandidate) types.ActiveState {
	now := time.Now()

	// Rule 1: Có EffectiveDate <= now AND (ExpiryDate == nil OR ExpiryDate > now) → active
	if c.EffectiveDate != nil && !c.EffectiveDate.After(now) {
		if c.ExpiryDate == nil || c.ExpiryDate.After(now) {
			return types.ActiveStateActive
		}
	}

	// Rule 2: ExpiryDate < now → not_active
	if c.ExpiryDate != nil && c.ExpiryDate.Before(now) {
		return types.ActiveStateNotActive
	}

	// Rule 3: EffectiveDate > now (chưa đến ngày hiệu lực) → not_active
	if c.EffectiveDate != nil && c.EffectiveDate.After(now) {
		return types.ActiveStateNotActive
	}

	// Rule 4: Thiếu thông tin → ambiguous
	return types.ActiveStateAmbiguous
}

// calculateConfidence tính độ tin cậy dựa trên metadata completeness
func (n *LegalStatusResolver) calculateConfidence(c *types.LayerCandidate) types.LegalConfidence {
	w := n.config.ConfidenceWeight
	score := 0

	if c.LegalStatus > 0 {
		score += w.HasLegalStatus
	}
	if c.EffectiveDate != nil {
		score += w.HasEffectiveDate
	}
	if c.Authority != "" {
		score += w.HasAuthority
	}
	if c.LegalDocRef != "" {
		score += w.HasLegalDoc
	}
	if float64(c.TrustValue) >= w.TrustValueThreshold {
		score++
	}

	switch {
	case score >= 4:
		return types.ConfidenceHigh
	case score >= 3:
		return types.ConfidenceMedium
	case score >= 2:
		return types.ConfidenceLow
	case score >= 1:
		return types.ConfidenceDegraded
	default:
		return types.ConfidenceError
	}
}

// generateWarnings sinh warning flags
func (n *LegalStatusResolver) generateWarnings(nl *types.NormalizedLegal) []string {
	var warnings []string

	if nl.Canonical == types.LegalUnknown {
		warnings = append(warnings, "unknown_legal_status: raw status không khớp canonical taxonomy")
	}
	if nl.ActiveState == types.ActiveStateAmbiguous {
		warnings = append(warnings, "ambiguous_active_state: thiếu EffectiveDate hoặc ExpiryDate")
	}
	if nl.Confidence == types.ConfidenceDegraded || nl.Confidence == types.ConfidenceError {
		warnings = append(warnings, "low_confidence: metadata không đầy đủ, kết quả có thể không chính xác")
	}

	return warnings
}

// generateReasonCodes sinh mã lý do
func (n *LegalStatusResolver) generateReasonCodes(c *types.LayerCandidate, nl *types.NormalizedLegal) []string {
	var codes []string

	codes = append(codes, fmt.Sprintf("raw_status=%d", nl.RawStatus))
	codes = append(codes, fmt.Sprintf("canonical=%s", nl.Canonical))
	codes = append(codes, fmt.Sprintf("active=%s", nl.ActiveState))
	codes = append(codes, fmt.Sprintf("confidence=%s", nl.Confidence))

	if c.EffectiveDate != nil {
		codes = append(codes, fmt.Sprintf("effective_date=%s", c.EffectiveDate.Format("2006-01-02")))
	}
	if c.ExpiryDate != nil {
		codes = append(codes, fmt.Sprintf("expiry_date=%s", c.ExpiryDate.Format("2006-01-02")))
	}

	return codes
}

// calculateScore tính legal score đã weighted by confidence
func (n *LegalStatusResolver) calculateScore(nl *types.NormalizedLegal) float64 {
	// Base score từ canonical status
	baseScore := n.canonicalBaseScore(nl.Canonical)

	// Weight by confidence
	confidenceMultiplier := n.confidenceMultiplier(nl.Confidence)

	// Weight by active state
	activeMultiplier := n.activeMultiplier(nl.ActiveState)

	return baseScore * confidenceMultiplier * activeMultiplier
}

// canonicalBaseScore trả về base score cho canonical status
func (n *LegalStatusResolver) canonicalBaseScore(canonical types.CanonicalLegalStatus) float64 {
	switch canonical {
	case types.LegalEffective:
		return 100.0
	case types.LegalApproved:
		return 85.0
	case types.LegalDraft:
		return 50.0
	case types.LegalReferenced:
		return 40.0
	case types.LegalExpired:
		return 20.0
	case types.LegalReplaced:
		return 10.0
	default:
		return 30.0
	}
}

// confidenceMultiplier trả về hệ số nhân cho confidence level
func (n *LegalStatusResolver) confidenceMultiplier(confidence types.LegalConfidence) float64 {
	switch confidence {
	case types.ConfidenceHigh:
		return 1.0
	case types.ConfidenceMedium:
		return 0.9
	case types.ConfidenceLow:
		return 0.7
	case types.ConfidenceDegraded:
		return 0.5
	default:
		return 0.3
	}
}

// activeMultiplier trả về hệ số nhân cho active state
func (n *LegalStatusResolver) activeMultiplier(active types.ActiveState) float64 {
	switch active {
	case types.ActiveStateActive:
		return 1.0
	case types.ActiveStateNotActive:
		return 0.4
	default:
		return 0.6
	}
}

// buildTrace xây dựng trace payload
func (n *LegalStatusResolver) buildTrace(c *types.LayerCandidate, nl *types.NormalizedLegal) map[string]any {
	return map[string]any{
		"engine":        n.Name(),
		"layer_id":      c.LayerID,
		"layer_name":    c.LayerName,
		"raw_status":    nl.RawStatus,
		"canonical":     nl.Canonical,
		"active_state":  nl.ActiveState,
		"confidence":    nl.Confidence,
		"score":         nl.Score,
		"warning_flags": nl.WarningFlags,
		"reason_codes":  nl.ReasonCodes,
	}
}
