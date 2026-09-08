package enums

import "fmt"

type LegalLevel uint32

const (
	LegalLevelCentral  LegalLevel = 10
	LegalLevelProvince LegalLevel = 20
	LegalLevelDistrict LegalLevel = 30
	LegalLevelWard     LegalLevel = 40
	LegalLevelOther    LegalLevel = 0
)

func (e LegalLevel) String() string {
	switch e {
	case LegalLevelCentral:
		return "Trung ương"
	case LegalLevelProvince:
		return "Tỉnh"
	case LegalLevelDistrict:
		return "Huyện"
	case LegalLevelWard:
		return "Xã"
	default:
		return "Khác"
	}
}

var LegalLevelMap = map[LegalLevel]string{
	LegalLevelCentral:  "Trung ương",
	LegalLevelProvince: "Tỉnh",
	LegalLevelDistrict: "Huyện",
	LegalLevelWard:     "Xã",
	LegalLevelOther:    "Khác",
}

func (e LegalLevel) IsValid() bool {
	return e >= LegalLevelCentral && e <= LegalLevelWard
}

func (e LegalLevel) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, e.String())), nil
}
