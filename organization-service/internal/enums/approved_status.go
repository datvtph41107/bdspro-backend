package enums

type ApprovedStatus int

const (
	ApprovedStatusPending  ApprovedStatus = 10
	ApprovedStatusApproved ApprovedStatus = 20
	ApprovedStatusRejected ApprovedStatus = 30
)

var ApprovedStatusMap = map[ApprovedStatus]string{
	ApprovedStatusPending:  "Chờ xác nhận",
	ApprovedStatusApproved: "Đã xác nhận",
	ApprovedStatusRejected: "Từ chối",
}

func (t ApprovedStatus) IsValid() bool {
	return t == ApprovedStatusPending || t == ApprovedStatusApproved || t == ApprovedStatusRejected
}
