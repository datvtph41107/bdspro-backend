package enums

type EMemberStatus int32

const (
	MemberStatusPending  EMemberStatus = 10
	MemberStatusActive   EMemberStatus = 20
	MemberStatusInActive EMemberStatus = 30
)

var MemberStatusMap = map[EMemberStatus]string{
	MemberStatusPending:  "Đang chờ xác thực",
	MemberStatusActive:   "Đang hoạt động",
	MemberStatusInActive: "Đã khóa",
}
