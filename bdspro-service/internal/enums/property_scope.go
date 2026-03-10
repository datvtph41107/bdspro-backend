package enums

type EPropertyScope uint32

const (
	PropertyScopePrivate  EPropertyScope = 10 // Riêng tư
	PropertyScopeShared   EPropertyScope = 20 // Được chia sẻ
	PropertyScopePublic   EPropertyScope = 30 // Dùng chung
	PropertyScopeInternal EPropertyScope = 40 // nội bộ
)

var EPropertyScopeNames = map[EPropertyScope]string{
	PropertyScopePrivate:  "Riêng tư",
	PropertyScopeShared:   "Được chia sẻ",
	PropertyScopePublic:   "Công khai",
	PropertyScopeInternal: "Nội bộ",
}

func (e EPropertyScope) String() string {
	return EPropertyScopeNames[e]
}

func (e EPropertyScope) IsValid() bool {
	switch e {
	case PropertyScopePrivate, PropertyScopeShared, PropertyScopePublic, PropertyScopeInternal:
		return true
	}
	return false
}
