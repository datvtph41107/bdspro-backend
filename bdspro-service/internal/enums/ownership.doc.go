package enums

type EDocType uint32

const (
	ERedBook         EDocType = 10 // sổ đỏ
	EPinkBook        EDocType = 20 // sổ hồng
	ENotarizedRecord EDocType = 30 // vi bằng
	EDocUnknown      EDocType = 40 // không xác định
)

var EDocTypeNames = map[EDocType]string{
	ERedBook:         "Sổ đỏ",
	EPinkBook:        "Sổ hồng",
	ENotarizedRecord: "Vi bằng",
	EDocUnknown:      "Không xác định",
}

func (e EDocType) String() string {
	return EDocTypeNames[e]
}

func (e EDocType) IsValid() bool {
	return e == ERedBook || e == EPinkBook || e == ENotarizedRecord || e == EDocUnknown
}

func (e EDocType) Parse(s *uint32) EDocType {
	if s == nil {
		return EDocUnknown
	}
	return EDocType(*s)
}

func (e EDocType) Value() uint32 {
	return uint32(e)
}
