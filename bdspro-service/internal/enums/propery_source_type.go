package enums

type EPropertySourceType uint32

const (
	PropertySourceSystem EPropertySourceType = 10 // System
	PropertySourceUser   EPropertySourceType = 20 // User
)

var EPropertySourceTypeNames = map[EPropertySourceType]string{
	PropertySourceSystem: "system",
	PropertySourceUser:   "user",
}

func (e EPropertySourceType) String() string {
	return EPropertySourceTypeNames[e]
}

func (e EPropertySourceType) IsValid() bool {
	switch e {
	case PropertySourceSystem, PropertySourceUser:
		return true
	}
	return false
}
