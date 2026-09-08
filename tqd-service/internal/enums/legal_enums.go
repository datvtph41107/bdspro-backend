package enums

// ============================================================
// LEGAL ENUMS — Tách riêng để dùng chung giữa các package
// ============================================================
// LOGICAL:
// - Surface: Cần enum dùng chung cho nhiều package
// - Root: Tránh circular import
// - Mechanism: Package enums chứa tất cả enum definitions

// DocumentType là uint32
type DocumentType uint32

const (
	DocTypePlanningDecision DocumentType = 10
	DocTypeLandCertificate  DocumentType = 20
	DocTypeLaw              DocumentType = 30
	DocTypeDecree           DocumentType = 40
	DocTypeRegulation       DocumentType = 50
	DocTypeCircular         DocumentType = 60
	DocTypeOther            DocumentType = 99
)

// DocumentTypeString map cho convert
var DocumentTypeString = map[DocumentType]string{
	DocTypePlanningDecision: "PLANNING_DECISION",
	DocTypeLandCertificate:  "LAND_CERTIFICATE",
	DocTypeLaw:              "LAW",
	DocTypeDecree:           "DECREE",
	DocTypeRegulation:       "REGULATION",
	DocTypeCircular:         "CIRCULAR",
	DocTypeOther:            "OTHER",
}

// DocumentTypeName map cho display
var DocumentTypeName = map[DocumentType]string{
	DocTypePlanningDecision: "Quyết định phê duyệt quy hoạch",
	DocTypeLandCertificate:  "Giấy chứng nhận QSDĐ",
	DocTypeLaw:              "Luật",
	DocTypeDecree:           "Nghị định",
	DocTypeRegulation:       "Quy chuẩn",
	DocTypeCircular:         "Thông tư",
	DocTypeOther:            "Văn bản pháp lý",
}

// DocumentStatus là uint32
type DocumentStatus uint32

const (
	DocStatusActive     DocumentStatus = 10
	DocStatusSuperseded DocumentStatus = 20
	DocStatusExpired    DocumentStatus = 30
	DocStatusDraft      DocumentStatus = 40
)

// DocumentStatusString map cho convert
var DocumentStatusString = map[DocumentStatus]string{
	DocStatusActive:     "ACTIVE",
	DocStatusSuperseded: "SUPERSEDED",
	DocStatusExpired:    "EXPIRED",
	DocStatusDraft:      "DRAFT",
}

// DocumentStatusLabel map cho display
var DocumentStatusLabel = map[DocumentStatus]string{
	DocStatusActive:     "Đang có hiệu lực",
	DocStatusSuperseded: "Đã bị thay thế",
	DocStatusExpired:    "Hết hiệu lực",
	DocStatusDraft:      "Dự thảo",
}

// DocumentStatusColor map cho UI
var DocumentStatusColor = map[DocumentStatus]string{
	DocStatusActive:     "#22C55E",
	DocStatusSuperseded: "#F59E0B",
	DocStatusExpired:    "#EF4444",
	DocStatusDraft:      "#8B5CF6",
}

// BuildStatus là uint32
type BuildStatus uint32

const (
	BuildStatusAllowed    BuildStatus = 10
	BuildStatusForbidden  BuildStatus = 20
	BuildStatusRestricted BuildStatus = 30
	BuildStatusUnknown    BuildStatus = 99
)

// BuildStatusString map cho convert
var BuildStatusString = map[BuildStatus]string{
	BuildStatusAllowed:    "ALLOWED",
	BuildStatusForbidden:  "FORBIDDEN",
	BuildStatusRestricted: "RESTRICTED",
	BuildStatusUnknown:    "UNKNOWN",
}

// BuildStatusLabel map cho display
var BuildStatusLabel = map[BuildStatus]string{
	BuildStatusAllowed:    "Được phép xây dựng",
	BuildStatusForbidden:  "Không được phép xây dựng",
	BuildStatusRestricted: "Xây dựng có điều kiện",
	BuildStatusUnknown:    "Chưa xác định",
}

// AlertLevel là uint32
type AlertLevel uint32

const (
	AlertLevelNone   AlertLevel = 0
	AlertLevelLow    AlertLevel = 10
	AlertLevelMedium AlertLevel = 20
	AlertLevelHigh   AlertLevel = 30
)

// AlertLevelString map cho convert
var AlertLevelString = map[AlertLevel]string{
	AlertLevelNone:   "NONE",
	AlertLevelLow:    "LOW",
	AlertLevelMedium: "MEDIUM",
	AlertLevelHigh:   "HIGH",
}
