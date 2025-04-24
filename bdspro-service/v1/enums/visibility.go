package enums

type VisibilityStatus uint

const (
	VisiShow        VisibilityStatus = 1 //"Đang hiển thị"
	VisiblePublic   VisibilityStatus = 2 // công khai
	VisibleInternal VisibilityStatus = 3 // nội bộ
)
