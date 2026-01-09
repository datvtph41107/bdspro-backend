package enums

type DistributionStatus uint

const (
	DISTRIBUTION_STATUS_ACTIVE  DistributionStatus = 10 // Đang có hiệu lực
	DISTRIBUTION_STATUS_REVOKED DistributionStatus = 20 // Đã thu quyền
	DISTRIBUTION_STATUS_EXPRIED DistributionStatus = 30 // Hết hiệu lực (tự nhiên)
)

var DistributionStatusMap = map[DistributionStatus]string{
	DISTRIBUTION_STATUS_ACTIVE:  "active",
	DISTRIBUTION_STATUS_REVOKED: "revoked",
	DISTRIBUTION_STATUS_EXPRIED: "expired",
}
