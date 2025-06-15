package enums

type FriendStatus int

const (
	Pending  FriendStatus = 10 //"PENDING"
	Accepted FriendStatus = 20 //"ACCEPT"
	Reject   FriendStatus = 30 //"REJECT"
)

var FriendStatusMap = map[FriendStatus]string{
	Pending:  "Gửi lời mời",
	Accepted: "Đã chấp nhận",
	Reject:   "Từ chối",
}
