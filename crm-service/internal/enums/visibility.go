package enums

type EVisibility uint

const (
	VisibilityPublic   EVisibility = 10
	VisibilityInternal EVisibility = 20
	VisibilityPrivate  EVisibility = 30
)

var VisibilityMap = map[EVisibility]string{
	VisibilityPublic:   "Công khai",
	VisibilityInternal: "Nội bộ",
	VisibilityPrivate:  "Chỉ mình tôi",
}