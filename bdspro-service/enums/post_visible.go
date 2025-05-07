package enums

type PostVisibility uint

const (
	PostVisibleShowing  PostVisibility = 1 // Đang hiển thị
	PostVisiblePublic   PostVisibility = 2 // Công khai
	PostVisibleInternal PostVisibility = 3 // Nội bộ
)
