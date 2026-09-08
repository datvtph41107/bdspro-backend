package enums

type KYCStatus int

const (
	KYCStatusPending  KYCStatus = 10 // Chờ duyệt
	KYCStatusApproved KYCStatus = 20 // Đã duyệt
	KYCStatusRejected KYCStatus = 30 // Từ chối
)

var KYCStatusMap = map[KYCStatus]string{
	KYCStatusPending:  "Chờ duyệt",
	KYCStatusApproved: "Đã duyệt",
	KYCStatusRejected: "Từ chối",
}
