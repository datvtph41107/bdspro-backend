package enums

type OwnershipDocType int8

const (
	RedBook         OwnershipDocType = 1 // sổ đỏ
	PinkBook        OwnershipDocType = 2 // sổ hồng
	NotarizedRecord OwnershipDocType = 3 // vi bằng
)

func (s OwnershipDocType) String() string {
	switch s {
	case RedBook:
		return "Sổ đỏ"
	case PinkBook:
		return "Sổ hồng"
	case NotarizedRecord:
		return "Vi bằng"
	default:
		return ""
	}
}
