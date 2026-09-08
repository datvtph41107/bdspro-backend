package enums

type CommissionType uint

const (
	CommissionTypePercent CommissionType = 10
	CommissionTypeVND     CommissionType = 20
)

var CommissionTypeMap = map[CommissionType]string{
	CommissionTypePercent: "Phần trăm (%)",
	CommissionTypeVND:     "Tiền VND",
}
