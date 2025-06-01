package enums

type VerifyStatus int

const (
	VerifyPending  VerifyStatus = 10
	VerifyVerified VerifyStatus = 20
	VerifyRejected VerifyStatus = 30
)

var VerifyStatusMap = map[VerifyStatus]string{
	VerifyPending:  "Chờ xác minh",
	VerifyVerified: "Đã xác minh",
	VerifyRejected: "Từ chối",
}
