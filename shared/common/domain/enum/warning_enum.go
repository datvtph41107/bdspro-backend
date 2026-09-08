package _enum

type EWarningLevel int

const (
	EWarningLevelNone     EWarningLevel = 0  // Không có cảnh báo
	EWarningLevelInfo     EWarningLevel = 10 // Thông tin
	EWarningLevelWarning  EWarningLevel = 20 // Cảnh báo
	EWarningLevelDanger   EWarningLevel = 30 // Nguy hiểm
	EWarningLevelCritical EWarningLevel = 40 // Nghiêm trọng
)

var WarningLevelNames = map[EWarningLevel]string{
	EWarningLevelNone:     "Không có cảnh báo",
	EWarningLevelInfo:     "Thông tin",
	EWarningLevelWarning:  "Cảnh báo",
	EWarningLevelDanger:   "Nguy hiểm",
	EWarningLevelCritical: "Nghiêm trọng",
}

func (e EWarningLevel) IsValid() bool {
	switch e {
	case EWarningLevelNone, EWarningLevelInfo, EWarningLevelWarning, EWarningLevelDanger, EWarningLevelCritical:
		return true
	default:
		return false
	}
}

func (e EWarningLevel) String() string {
	if name, ok := WarningLevelNames[e]; ok {
		return name
	}
	return "Unknown"
}
