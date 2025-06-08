package enums

type Visibility int

const (
	VisibilityPublic  Visibility = 10
	VisibilityFriend  Visibility = 30
	VisibilityPrivate Visibility = 20
	VisibilityGroup   Visibility = 40
)

var VisibilityMap = map[Visibility]string{
	VisibilityPublic:  "Công khai",
	VisibilityFriend:  "Bạn bè",
	VisibilityPrivate: "Riêng tư",
	VisibilityGroup:   "Nhóm",
}
