package enums

type PostStatus int

const (
	PostPublic          PostStatus = 10
	PostFinished        PostStatus = 20
	PostSold            PostStatus = 40
	PostRent            PostStatus = 50
	PostStop            PostStatus = 60 // tùy nghiệp vụ về sau
	PostStatusCompleted PostStatus = 70
	// PostHide            PostStatus = 3
)

func (s PostStatus) String() string {
	switch s {
	case PostPublic:
		return "Công khai"
	case PostFinished:
		return "Hết hạn"
	// case PostHide:
	// 	return "Tạm ẩn"
	case PostSold:
		return "Đã bán"
	case PostRent:
		return "Đã cho thuê"
	case PostStop:
		return "Dừng"
		// return "Đã xóa"
	default:
		return ""
	}
}
