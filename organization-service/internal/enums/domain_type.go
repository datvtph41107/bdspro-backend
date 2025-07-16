package enums

const (
	DomainTypeDeal = 50
)

type DomainType int

const (
	Deal DomainType = DomainTypeDeal
)

func (d DomainType) String() string {
	switch d {
	case Deal:
		return "deal"
	default:
		return "unknown"
	}
}

func (d DomainType) Value() int {
	return int(d)
} 