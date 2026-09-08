package enums

type DealType uint

const (
	DealTypeBrokerage DealType = 10
	DealTypeInvest    DealType = 20
	DealTypeJoint     DealType = 30

	DealTypeSale       DealType = 40
	DealTypeRent       DealType = 50
	DealTypeInvestment DealType = 60
	DealTypeOther      DealType = 100
)

var DealTypeNames = map[DealType]string{
	// DealTypeBrokerage: "Môi giới",
	// DealTypeInvest:    "Góp vốn",
	// DealTypeJoint:     "Kết hợp",
	DealTypeSale:       "Bán",
	DealTypeRent:       "Thuê",
	DealTypeInvestment: "Đầu tư",
	DealTypeOther:      "Khác",
}
