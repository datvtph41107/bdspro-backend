package enums

type EPropertySourceType uint32

const (
	PropertySourceCanonical EPropertySourceType = 10 // System
	PropertySourceBase      EPropertySourceType = 20 // User
)

var EPropertySourceTypeNames = map[EPropertySourceType]string{
	PropertySourceCanonical: "canonical", // system
	PropertySourceBase:      "base",      // user
}

func (e EPropertySourceType) String() string {
	return EPropertySourceTypeNames[e]
}

func (e EPropertySourceType) IsValid() bool {
	switch e {
	case PropertySourceCanonical, PropertySourceBase:
		return true
	}
	return false
}
