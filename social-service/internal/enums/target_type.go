package enums

type TargetType int

const (
	TargetTypeComment  TargetType = 10
	TargetTypeNewsFeed TargetType = 20
)

var TargetTypeMap = map[TargetType]string{
	TargetTypeComment:  "Bình luận",
	TargetTypeNewsFeed: "Bài viết",
}
