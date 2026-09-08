package enums

// ─────────────────────────────────────────────────────────────────────────────
// EReportStatus — trạng thái xử lý báo cáo
// ─────────────────────────────────────────────────────────────────────────────

type EReportStatus int16

const (
	EReportStatusPending   EReportStatus = 10 // Chờ xử lý
	EReportStatusReviewing EReportStatus = 20 // Đang xem xét
	EReportStatusApproved  EReportStatus = 30 // Đã duyệt — đã merge propose
	EReportStatusRejected  EReportStatus = 40 // Từ chối
	EReportStatusClosed    EReportStatus = 50 // Đóng (không cần xử lý)
)
