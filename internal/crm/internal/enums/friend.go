package enums

type FriendStatus int

const (
	FriendStatusPending  FriendStatus = 10 //"PENDING"
	FriendStatusAccepted FriendStatus = 20 //"ACCEPT"
	FriendStatusReject   FriendStatus = 30 //"REJECT"
)

var FriendStatusMap = map[FriendStatus]string{
	FriendStatusPending:  "Gửi lời mời",
	FriendStatusAccepted: "Đã chấp nhận",
	FriendStatusReject:   "Từ chối",
}