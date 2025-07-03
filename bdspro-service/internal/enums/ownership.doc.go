package enums

type EDocType uint32

const (
	ERedBook         EDocType = 10 // sổ đỏ
	EPinkBook        EDocType = 20 // sổ hồng
	ENotarizedRecord EDocType = 30 // vi bằng
)

var EDocTypeNames = map[EDocType]string{
	ERedBook:         "Sổ đỏ",
	EPinkBook:        "Sổ hồng",
	ENotarizedRecord: "Vi bằng",
}
