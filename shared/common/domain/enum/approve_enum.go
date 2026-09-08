package _enum

type EApproveStatus int

const (
	EApproveStatusPending  EApproveStatus = 10
	EApproveStatusApproved EApproveStatus = 20
	EApproveStatusRejected EApproveStatus = 30
)

var ApproveStatusMap = map[EApproveStatus]string{
	EApproveStatusPending:  "Chờ duyệt",
	EApproveStatusApproved: "Đã duyệt",
	EApproveStatusRejected: "Từ chối",
}
