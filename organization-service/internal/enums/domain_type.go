package enums

type DomainType int

const (
	DomainTypeDeal DomainType = 50
	DomainTypeOrg  DomainType = 100
)

func (d DomainType) String() string {
	switch d {
	case DomainTypeDeal:
		return "deal"
	case DomainTypeOrg:
		return "org"
	default:
		return "unknown"
	}
}

func (d DomainType) Value() int {
	return int(d)
}
