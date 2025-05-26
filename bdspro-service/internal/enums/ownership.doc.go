package enums

type OwnershipDocType int8

const (
	RedBook         OwnershipDocType = 10 // sổ đỏ
	PinkBook        OwnershipDocType = 20 // sổ hồng
	NotarizedRecord OwnershipDocType = 30 // vi bằng
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
