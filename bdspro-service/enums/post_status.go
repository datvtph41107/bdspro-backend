package enums

type PostStatus int

const (
	PostPublic          PostStatus = 1
	PostFinished        PostStatus = 2
	PostHide            PostStatus = 3
	PostSold            PostStatus = 4
	PostRent            PostStatus = 5
	PostStop            PostStatus = 6 // tùy nghiệp vụ về sau
	PostStatusCompleted PostStatus = 7
)

func (s PostStatus) String() string {
	switch s {
	case PostPublic:
		return "Công khai"
	case PostFinished:
		return "Hết hạn"
	case PostHide:
		return "Tạm ẩn"
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
