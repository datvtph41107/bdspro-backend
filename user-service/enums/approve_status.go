package enums

type ApproveStatus uint32

const (
	StatusPending  ApproveStatus = 10
	StatusApproved ApproveStatus = 20
	StatusRejected ApproveStatus = 30
)

var ApproveStatusMap = map[ApproveStatus]string{
	StatusPending:  "Đang xử lý",
	StatusApproved: "Đã giải quyết",
	StatusRejected: "Từ chối",
}
