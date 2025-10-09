package enums

type EFriendStatus int

const (
	EFriendStatusPending  EFriendStatus = 10 //"PENDING"
	EFriendStatusAccepted EFriendStatus = 20 //"ACCEPT"
	EFriendStatusReject   EFriendStatus = 30 //"REJECT"
)
