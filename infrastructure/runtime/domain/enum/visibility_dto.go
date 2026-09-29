package _enum

type EVisibility uint32

const (
	EVisibilityDTOPublic   EVisibility = 10
	EVisibilityDTOPrivate  EVisibility = 20
	EVisibilityDTOInternal EVisibility = 30
)

var EVisibilityDTONames = map[EVisibility]string{
	EVisibilityDTOPublic:   "Công khai",
	EVisibilityDTOPrivate:  "Riêng tư",
	EVisibilityDTOInternal: "Nội bộ",
}

func (e EVisibility) IsValid() bool {
	switch e {
	case EVisibilityDTOPublic, EVisibilityDTOPrivate, EVisibilityDTOInternal:
		return true
	default:
		return false
	}
}
