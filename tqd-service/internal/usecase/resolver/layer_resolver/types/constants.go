package types

// =====================================================
// SINGLE SOURCE OF TRUTH - TOÀN BỘ GIÁ TRỊ MẶC ĐỊNH
// Thay đổi ở đây → ảnh hưởng toàn bộ hệ thống
// =====================================================

// ---------- Scoring weights ----------
var (
	DefaultLegalWeight    = 0.60
	DefaultPriorityWeight = 0.30
	DefaultAreaWeight     = 0.10
)

// ---------- Legal scores ----------
var DefaultLegalScores = map[uint32]float64{
	1000: 100,
	900:  70,
	800:  30,
	700:  15,
	400:  10,
	300:  10,
}

// ---------- Risk weights ----------
var DefaultLandUseRiskWeights = map[string]*LandUseRiskWeight{
	"DGT": {Multiplier: 2.0, Label: "không được phép xây dựng"},
	"DHT": {Multiplier: 2.0, Label: "không được phép xây dựng"},
	"DCX": {Multiplier: 1.0, Label: "hạn chế xây dựng"},
	"DSH": {Multiplier: 1.0, Label: "hạn chế xây dựng"},
	"ODT": {Multiplier: 0.35, Label: "phù hợp mục đích ở"},
	"ONT": {Multiplier: 0.35, Label: "phù hợp mục đích ở"},
}

var DefaultRiskWeightValue = &LandUseRiskWeight{Multiplier: 0.8, Label: "cần kiểm tra thêm"}

var DefaultMultiLayerPenaltyValue = 5.0

// ---------- Risk thresholds ----------
var DefaultRiskThresholds = []*RiskThreshold{
	{MinScore: 85, Level: "CRITICAL", CanBuild: false, Recommendation: "Rủi ro rất cao, không nên giao dịch trước khi xác minh quy hoạch."},
	{MinScore: 60, Level: "HIGH", CanBuild: false, Recommendation: "Rủi ro cao, cần kiểm tra kỹ quy hoạch trước khi quyết định."},
	{MinScore: 30, Level: "MEDIUM", CanBuild: true, Recommendation: "Có rủi ro quy hoạch, nên kiểm tra thêm hồ sơ pháp lý."},
	{MinScore: 0, Level: "LOW", CanBuild: true, Recommendation: "Rủi ro thấp, có thể tiếp tục xem xét."},
}

// ---------- Priority ----------
var DefaultPriorityValue = 50

// ---------- Conflict labels ----------
const (
	DefaultConflictType           = "land_use_mismatch"
	DefaultConflictSeverity       = "high"
	DefaultConflictRecommendation = "Cần đối chiếu hồ sơ pháp lý gốc"
	DefaultConflictActionRequired = "verify"
)

// ---------- Warning messages ----------
const (
	DefaultWarningOverlapExceeded = "Phát hiện chồng lấn quy hoạch vượt quá diện tích thửa đất"
	DefaultWarningMultiLayer      = "Thửa đất đang chịu ảnh hưởng từ nhiều layer quy hoạch"
	DefaultWarningConflict        = "Phát hiện xung đột giữa các quy hoạch"
)

// ---------- Assess messages ----------
const (
	DefaultNoPlanningImpactMsg    = "Không phát hiện ảnh hưởng quy hoạch"
	DefaultNoSignificantImpactMsg = "Không phát hiện ảnh hưởng quy hoạch đáng kể"
)

// ---------- Build status labels ----------
const (
	DefaultStatusProhibited     = "prohibited"
	DefaultStatusConditional    = "conditional"
	DefaultStatusAllowed        = "allowed"
	DefaultConflictStatusSuffix = " (Có mâu thuẫn với quy hoạch khác)"
)

// ---------- Land type names ----------
var DefaultLandTypeNames = map[uint64]string{
	1:  "Đất ở đô thị",
	2:  "Đất ở nông thôn",
	3:  "Đất nông nghiệp",
	4:  "Đất lâm nghiệp",
	5:  "Đất thương mại",
	6:  "Đất dịch vụ",
	7:  "Đất công nghiệp",
	8:  "Đất cây xanh",
	9:  "Đất giao thông",
	10: "Đất hạ tầng",
	11: "Đất công cộng",
	12: "Đất hỗn hợp",
}

// ---------- Summary labels ----------
const DefaultConflictSummarySuffix = " - Có xung đột quy hoạch"
