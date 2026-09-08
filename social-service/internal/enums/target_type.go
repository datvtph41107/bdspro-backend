package enums

type TargetType int

const (
	TargetTypeComment  TargetType = 10
	TargetTypeNewsFeed TargetType = 20

	// user service
	TargetTypeRate    TargetType = 30
	TargetTypeProfile TargetType = 40
)

var TargetTypeMap = map[TargetType]string{
	TargetTypeComment:  "Bình luận",
	TargetTypeNewsFeed: "Bài viết",
	TargetTypeRate:     "Đánh giá",
	TargetTypeProfile:  "Người dùng",
}
