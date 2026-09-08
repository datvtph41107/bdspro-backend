package enums

type Visibility int

const (
	VisibilityPublic  Visibility = 10
	VisibilityFriend  Visibility = 20
	VisibilityPrivate Visibility = 30
	VisibilityGroup   Visibility = 40
)

var VisibilityMap = map[Visibility]string{
	VisibilityPublic:  "Công khai",
	VisibilityFriend:  "Bạn bè",
	VisibilityPrivate: "Riêng tư",
	VisibilityGroup:   "Nhóm",
}
