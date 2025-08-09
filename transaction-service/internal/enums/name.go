package enums

type EOwnerType uint32

const (
	EOwnerUser EOwnerType = 10
	EOwnerDeal EOwnerType = 20
)

var DomainNames = map[EOwnerType]string{
	EOwnerUser: "User",
	EOwnerDeal: "Deal",
}
