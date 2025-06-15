package enums

type EVisibility uint32

const (
	EVisiblePublic   EVisibility = 10
	EVisibleInternal EVisibility = 20
	EVisiblePrivate  EVisibility = 30
)

var EVisibilityNames = map[EVisibility]string{
	EVisiblePublic:   "Công khai",
	EVisibleInternal: "Nội bộ",
	EVisiblePrivate:  "Riêng tư",
}
