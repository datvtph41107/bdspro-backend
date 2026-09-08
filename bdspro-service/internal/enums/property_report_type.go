package enums

// EReportType — loại báo cáo chính, map với UI "Chọn loại báo cáo"
type EReportType int16

const (
	EReportTypeWrongInfo     EReportType = 10 // Sai thông tin
	EReportTypeWrongLocation EReportType = 20 // Sai vị trí
	EReportTypeDuplicate     EReportType = 30 // Trùng / Nhầm
	EReportTypeSensitive     EReportType = 40 // Nội dung nhạy cảm
	EReportTypeSpam          EReportType = 50 // Spam / Rác
	EReportTypeOther         EReportType = 60 // Khác
)

func (e EReportType) IsValid() bool {
	switch e {
	case EReportTypeWrongInfo,
		EReportTypeWrongLocation,
		EReportTypeDuplicate,
		EReportTypeSensitive,
		EReportTypeSpam,
		EReportTypeOther:
		return true
	}
	return false
}

// ─────────────────────────────────────────────────────────────────────────────
// EReportSubIssue — các sub-issue theo từng loại báo cáo
// Dùng để validate và map từ proto repeated string
// ─────────────────────────────────────────────────────────────────────────────

// Sub-issues for EReportTypeWrongInfo
const (
	SubIssueWrongName    = "wrong_name"    // Tên BĐS/Alias
	SubIssueWrongType    = "wrong_type"    // Loại BĐS
	SubIssueWrongArea    = "wrong_area"    // Diện tích
	SubIssueWrongProject = "wrong_project" // Dự án
	SubIssueWrongFloor   = "wrong_floor"   // Hệ tầng
	SubIssueWrongAmenity = "wrong_amenity" // Tiện ích
	SubIssueWrongLegal   = "wrong_legal"   // Pháp lý
	SubIssueWrongImage   = "wrong_image"   // Hình ảnh
	SubIssueWrongDesc    = "wrong_desc"    // Mô tả
)

// Sub-issues for EReportTypeWrongLocation
const (
	SubIssueWrongAddress = "wrong_address" // Sai địa chỉ cụ thể
	SubIssueWrongCoords  = "wrong_coords"  // Sai tọa độ
	SubIssueWrongPolygon = "wrong_polygon" // Sai ranh thửa/polygon
	SubIssueWrongMapRef  = "wrong_map_ref" // Không phù hợp thực tế
)

// Sub-issues for EReportTypeDuplicate
const (
	SubIssueDupOtherBDS   = "dup_other_bds"  // Trùng BĐS khác
	SubIssueDupWrongInfo2 = "dup_wrong_info" // Nhầm thông tin
)

// Sub-issues for EReportTypeSensitive
const (
	SubIssueSensitiveCertImg   = "sensitive_cert_img"   // Ảnh có giấy tờ
	SubIssueSensitivePersonImg = "sensitive_person_img" // Ảnh có mặt người
	SubIssueSensitiveContent   = "sensitive_content"    // Nội dung nhạy cảm
	SubIssueSensitiveDesc2     = "sensitive_desc"       // Văn bản mô tả nhạy cảm
)

// Sub-issues for EReportTypeSpam (wrong info type reuse)
const (
	SubIssueSpamBDSRepeat = "spam_repeat" // BDS bị lặp lại nhiều lần
	SubIssueSpamSystem    = "spam_system" // Lỗi hệ thống
)
