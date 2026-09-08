package types

// CanonicalLegalStatus — kết quả chuẩn hóa trạng thái pháp lý
type CanonicalLegalStatus string

const (
	LegalEffective  CanonicalLegalStatus = "effective"
	LegalApproved   CanonicalLegalStatus = "approved"
	LegalDraft      CanonicalLegalStatus = "draft"
	LegalExpired    CanonicalLegalStatus = "expired"
	LegalReplaced   CanonicalLegalStatus = "replaced"
	LegalReferenced CanonicalLegalStatus = "referenced"
	LegalUnknown    CanonicalLegalStatus = "unknown"
)

// ActiveState — trạng thái hiệu lực thực tế của layer
type ActiveState string

const (
	ActiveStateActive    ActiveState = "active"
	ActiveStateNotActive ActiveState = "not_active"
	ActiveStateAmbiguous ActiveState = "ambiguous"
)

// LegalConfidence — độ tin cậy của kết quả chuẩn hóa
type LegalConfidence string

const (
	ConfidenceHigh     LegalConfidence = "high"
	ConfidenceMedium   LegalConfidence = "medium"
	ConfidenceLow      LegalConfidence = "low"
	ConfidenceDegraded LegalConfidence = "degraded"
	ConfidenceError    LegalConfidence = "error"
)

// NormalizedLegal — kết quả đầu ra của Legal Normalization Engine
type NormalizedLegal struct {
	RawStatus    uint32               // input gốc từ DB
	Canonical    CanonicalLegalStatus // chuẩn hóa về canonical taxonomy
	ActiveState  ActiveState          // active / not_active / ambiguous
	Confidence   LegalConfidence      // độ tin cậy của kết quả
	Score        float64              // điểm pháp lý (0-100), đã weighted by confidence
	WarningFlags []string             // cảnh báo nếu có ambiguity
	ReasonCodes  []string             // mã lý do cho trace/audit
	TracePayload map[string]any       // trace payload cho audit
}

// LegalNormalizationConfig — cấu hình cho Legal Normalization Engine
type LegalNormalizationConfig struct {
	CanonicalMap     map[uint32]CanonicalLegalStatus // raw status → canonical
	ConfidenceWeight ConfidenceWeightConfig
	ActiveRules      []ActiveRuleConfig
}

// ConfidenceWeightConfig — trọng số tính confidence
type ConfidenceWeightConfig struct {
	HasLegalStatus      int     // +1 nếu có LegalStatus > 0
	HasEffectiveDate    int     // +1 nếu có EffectiveDate
	HasAuthority        int     // +1 nếu có Authority
	HasLegalDoc         int     // +1 nếu có LegalDocRef
	TrustValueThreshold float64 // ngưỡng trust value để +1
}

// ActiveRuleConfig — rule đánh giá active state
type ActiveRuleConfig struct {
	Name      string
	Condition string // mô tả điều kiện
	Result    string // "active" | "not_active" | "ambiguous"
}

// DefaultCanonicalMap — ánh xạ mặc định raw → canonical
func DefaultCanonicalMap() map[uint32]CanonicalLegalStatus {
	return map[uint32]CanonicalLegalStatus{
		1000: LegalEffective,
		900:  LegalApproved,
		800:  LegalDraft,
		700:  LegalReferenced,
		400:  LegalExpired,
		300:  LegalReplaced,
		0:    LegalUnknown,
	}
}

// DefaultConfidenceWeight — trọng số confidence mặc định
func DefaultConfidenceWeight() ConfidenceWeightConfig {
	return ConfidenceWeightConfig{
		HasLegalStatus:      1,
		HasEffectiveDate:    1,
		HasAuthority:        1,
		HasLegalDoc:         1,
		TrustValueThreshold: 0.8,
	}
}

// DefaultActiveRules — rules đánh giá active state mặc định
func DefaultActiveRules() []ActiveRuleConfig {
	return []ActiveRuleConfig{
		{
			Name:      "effective_within_range",
			Condition: "effective_date <= now AND (expiry_date IS NULL OR expiry_date > now)",
			Result:    "active",
		},
		{
			Name:      "expired",
			Condition: "expiry_date < now",
			Result:    "not_active",
		},
		{
			Name:      "future_effective",
			Condition: "effective_date > now",
			Result:    "not_active",
		},
	}
}
