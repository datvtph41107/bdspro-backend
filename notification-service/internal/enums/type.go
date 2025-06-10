package enums

// TypeEnum định nghĩa các loại thông báo
type TypeEnum int8

const (
	FriendRequest    TypeEnum = 10 //"FRIEND_REQUEST"
	FriendResponse   TypeEnum = 20 //"FRIEND_RESPONSE"
	Follow           TypeEnum = 30 //"FOLLOW"
	CommentNewsFeed  TypeEnum = 40 //"COMMENT_NEWS_FEED"
	CommentChildren  TypeEnum = 41 //"COMMENT_REPLY"
	LikeNewsFeed     TypeEnum = 50 //"LIKE_NEWS_FEED"
	LikeComment      TypeEnum = 60 //"LIKE_COMMENT"
	ShareNewsFeed    TypeEnum = 70 //"SHARE_NEWS_FEED"
	JoinGroup        TypeEnum = 80
	JoinOrganization TypeEnum = 90
)

var TypeEnumMap = map[TypeEnum]string{
	FriendRequest:    "Gửi yêu cầu kết bạn",
	FriendResponse:   "Kết bạn",
	Follow:           "Theo dõi",
	CommentNewsFeed:  "Bình luận bài viết",
	CommentChildren:  "Bình luận bài viết",
	LikeNewsFeed:     "Thích bài viết",
	LikeComment:      "Thích bình luận",
	ShareNewsFeed:    "Chia sẻ bài viết",
	JoinGroup:        "Tham gia nhóm",
	JoinOrganization: "Tham gia tổ chức",
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
