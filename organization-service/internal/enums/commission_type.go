package enums

type CommissionType uint

const (
	CommissionTypePercent CommissionType = 10 // Hoa hồng theo phần trăm
	CommissionTypeVND    CommissionType = 20 // Hoa hồng theo tiền VND
)

var CommissionTypeMap = map[CommissionType]string{
	CommissionTypePercent: "Phần trăm (%)",
	CommissionTypeVND:     "Tiền VND",
} 