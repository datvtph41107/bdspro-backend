package enums

type PostVisibility uint

const (
	PostVisibleShowing  PostVisibility = 10 // Đang hiển thị
	PostVisiblePublic   PostVisibility = 20 // Công khai
	PostVisibleInternal PostVisibility = 30 // Nội bộ
)
