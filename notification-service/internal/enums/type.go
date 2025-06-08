package enums

// TypeEnum định nghĩa các loại thông báo
type TypeEnum int8

const (
	FriendRequest  TypeEnum = 10 //"FRIEND_REQUEST"
	FriendResponse TypeEnum = 20 //"FRIEND_RESPONSE"
	Follow         TypeEnum = 30 //"FOLLOW"
	Comment        TypeEnum = 40 //"COMMENT"
	Like           TypeEnum = 50 //"LIKE"
	Share          TypeEnum = 60 //"SHARE"
	FollowRequest  TypeEnum = 70 //"FOLLOW_REQUEST"
	FollowResponse TypeEnum = 80 //"FOLLOW_RESPONSE"
)

var TypeEnumMap = map[TypeEnum]string{
	FriendRequest:  "Gửi yêu cầu kết bạn",
	FriendResponse: "Kết bạn",
	Follow:         "Theo dõi",
	Comment:        "Bình luận",
	Like:           "Thích",
	Share:          "Chia sẻ",
	FollowRequest:  "Gửi yêu cầu theo dõi",
	FollowResponse: "Theo dõi",
}

// IsValid kiểm tra xem giá trị có hợp lệ không
func (t TypeEnum) IsValid() bool {
	switch t {
	case FriendRequest, FriendResponse, Follow:
		return true
	default:
		return false
	}
}
