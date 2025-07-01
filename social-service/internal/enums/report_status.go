package enums

type ReportStatus uint8

const (
	ReportStatusPending    ReportStatus = 10
	ReportStatusProcessing ReportStatus = 20
	ReportStatusDone       ReportStatus = 30
	ReportStatusCancel     ReportStatus = 40
)

var ReportStatusMap = map[ReportStatus]string{
	ReportStatusPending:    "Chưa xử lý",
	ReportStatusProcessing: "Đang xử lý",
	ReportStatusDone:       "Đã xử lý",
	ReportStatusCancel:     "Đã hủy",
}
