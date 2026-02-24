package enums

type EPropertyScope uint32

const (
	PropertyScopePrivate EPropertyScope = 10 // Riêng tư
	PropertyScopeShared  EPropertyScope = 20 // Được chia sẻ
	PropertyScopePublic  EPropertyScope = 30 // Dùng chung
)

var EPropertyScopeNames = map[EPropertyScope]string{
	PropertyScopePrivate: "Riêng tư",
	PropertyScopeShared:  "Được chia sẻ",
	PropertyScopePublic:  "Dùng chung",
}

func (e EPropertyScope) String() string {
	return EPropertyScopeNames[e]
}

func (e EPropertyScope) IsValid() bool {
	switch e {
	case PropertyScopePrivate, PropertyScopeShared, PropertyScopePublic:
		return true
	}
	return false
}
