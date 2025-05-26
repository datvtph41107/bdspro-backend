package enums

type VisibilityStatus uint

const (
	VisiShow        VisibilityStatus = 10 //"Đang hiển thị"
	VisiblePublic   VisibilityStatus = 20 // công khai
	VisibleInternal VisibilityStatus = 30 // nội bộ
)

// type A struct {
// 	B *B
// }
// type B struct {
// 	A *A
// }
