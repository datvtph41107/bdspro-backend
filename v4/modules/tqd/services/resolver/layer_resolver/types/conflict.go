package types

// ConflictType — phân loại xung đột (6 loại theo SRS IV.3.2.3)
type ConflictType string

const (
	ConflictSemantic ConflictType = "semantic" // khác loại đất (GroupCode khác nhau)
	ConflictLegal    ConflictType = "legal"    // khác legal status canonical
	ConflictTemporal ConflictType = "temporal" // khác thời điểm hiệu lực
	ConflictSource   ConflictType = "source"   // khác nguồn dữ liệu
	ConflictGeometry ConflictType = "geometry" // chồng lấn hình học
	ConflictMixed    ConflictType = "mixed"    // kết hợp nhiều loại
)

// ClassifiedConflict — kết quả phân loại conflict chi tiết
type ClassifiedConflict struct {
	Type            ConflictType
	Severity        ConflictSeverity
	Confidence      float64  // 0.0 - 1.0
	Between         []uint64 // layer IDs
	Description     string
	ReasonCode      string
	Recommendation  string
	ActionRequired  string
	AffectedPercent float64
	IsPartial       bool     // partial region conflict?
	PartialRegions  []uint64 // region IDs bị ảnh hưởng 1 phần
}

// ConflictClassificationConfig — cấu hình cho Conflict Classifier
type ConflictClassificationConfig struct {
	SeverityThresholds SeverityThresholdConfig
	TypeWeights        map[ConflictType]float64
	GuidanceTemplates  map[ConflictType]*GuidanceTemplate
}

// SeverityThresholdConfig — ngưỡng phân loại severity
type SeverityThresholdConfig struct {
	High   float64
	Medium float64
}

// GuidanceTemplate — template mô tả và khuyến nghị cho mỗi loại conflict
type GuidanceTemplate struct {
	DescriptionTemplate string
	Recommendation      string
	ActionRequired      string
}

// DefaultSeverityThresholds — ngưỡng mặc định
func DefaultSeverityThresholds() SeverityThresholdConfig {
	return SeverityThresholdConfig{
		High:   70,
		Medium: 40,
	}
}

// DefaultConflictTypeWeights — trọng số mặc định cho từng loại conflict
func DefaultConflictTypeWeights() map[ConflictType]float64 {
	return map[ConflictType]float64{
		ConflictSemantic: 1.0,
		ConflictLegal:    0.8,
		ConflictTemporal: 0.6,
		ConflictSource:   0.4,
		ConflictGeometry: 0.5,
	}
}

// DefaultGuidanceTemplates — templates mặc định
func DefaultGuidanceTemplates() map[ConflictType]*GuidanceTemplate {
	return map[ConflictType]*GuidanceTemplate{
		ConflictSemantic: {
			DescriptionTemplate: "{a_landuse} ({a_layer}) xung đột với {b_landuse} ({b_layer})",
			Recommendation:      "Cần đối chiếu quy hoạch gốc để xác định mục đích sử dụng đất chính xác",
			ActionRequired:      "verify_land_use",
		},
		ConflictLegal: {
			DescriptionTemplate: "{a_layer} ({a_legal}) khác hiệu lực pháp lý với {b_layer} ({b_legal})",
			Recommendation:      "Ưu tiên layer có hiệu lực pháp lý cao hơn",
			ActionRequired:      "verify_legal_status",
		},
		ConflictTemporal: {
			DescriptionTemplate: "{a_layer} (hiệu lực {a_date}) và {b_layer} (hiệu lực {b_date}) khác thời điểm",
			Recommendation:      "Kiểm tra thời điểm hiệu lực của từng quy hoạch",
			ActionRequired:      "verify_timeline",
		},
		ConflictSource: {
			DescriptionTemplate: "{a_layer} (nguồn: {a_source}) và {b_layer} (nguồn: {b_source}) khác nguồn dữ liệu",
			Recommendation:      "Xác minh độ tin cậy của từng nguồn dữ liệu",
			ActionRequired:      "verify_source",
		},
		ConflictGeometry: {
			DescriptionTemplate: "{a_layer} và {b_layer} chồng lấn hình học",
			Recommendation:      "Kiểm tra ranh giới quy hoạch, có thể có sai số geometry",
			ActionRequired:      "verify_geometry",
		},
	}
}
