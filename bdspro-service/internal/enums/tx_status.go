package enums

type TxStatus uint32

const (
	TxStatusPending TxStatus = 10
	TxStatusDone    TxStatus = 20
	TxStatusCancel  TxStatus = 30
)

var TxStatusMap = map[TxStatus]string{
	TxStatusPending: "Đang diễn ra",
	TxStatusDone:    "Hoàn thành",
	TxStatusCancel:  "Hủy",
}
