package enums

type DealType uint

const (
	DealTypeBrokerage DealType = 10
	DealTypeInvest    DealType = 20
	DealTypeJoint     DealType = 30
)

var DealTypeMap = map[DealType]string{
	DealTypeBrokerage: "Môi giới",
	DealTypeInvest:    "Góp vốn",
	DealTypeJoint:     "Kết hợp",
}
