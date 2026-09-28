package enums

type TxApprovedStatus int

const (
	TxApprovedStatusPending  TxApprovedStatus = 10
	TxApprovedStatusApproved TxApprovedStatus = 20
	TxApprovedStatusRejected TxApprovedStatus = 30
)

var TxApprovedStatusMap = map[TxApprovedStatus]string{
	TxApprovedStatusPending:  "Chờ xác nhận",
	TxApprovedStatusApproved: "Đã xác nhận",
	TxApprovedStatusRejected: "Từ chối",
}

func (t TxApprovedStatus) IsValid() bool {
	return t == TxApprovedStatusPending || t == TxApprovedStatusApproved || t == TxApprovedStatusRejected
}
