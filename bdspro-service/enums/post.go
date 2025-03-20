package enums

type PostStatus int

const (
	PostPublic   PostStatus = 1
	PostFinished PostStatus = 2
	PostHide     PostStatus = 3
	PostSold     PostStatus = 4
	PostRent     PostStatus = 5
	PostDelete   PostStatus = 6 // tùy nghiệp vụ về sau
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
	case PostDelete:
		return "Đã xóa"
	default:
		return ""
	}
}
