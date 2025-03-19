package enums

// TypeEnum định nghĩa các loại thông báo
type TypeEnum int8

const (
	FriendRequest  TypeEnum = 1 //"FRIEND_REQUEST"
	FriendResponse TypeEnum = 2 //"FRIEND_RESPONSE"
	Follow         TypeEnum = 3 //"FOLLOW"
)

// IsValid kiểm tra xem giá trị có hợp lệ không
func (t TypeEnum) IsValid() bool {
	switch t {
	case FriendRequest, FriendResponse, Follow:
		return true
	default:
		return false
	}
}
