package enums

type ReportStatus uint32

const (
	ReportStatusPending  ReportStatus = 10 // Chờ xử lý
	ReportStatusApproved ReportStatus = 20 // Đã duyệt
	ReportStatusRejected ReportStatus = 30 // Từ chối
	ReportStatusResolved ReportStatus = 40 // Đã giải quyết
)

var ReportStatusMap = map[ReportStatus]string{
	ReportStatusPending:  "Chờ xử lý",
	ReportStatusApproved: "Đã duyệt",
	ReportStatusRejected: "Từ chối",
	ReportStatusResolved: "Đã giải quyết",
}