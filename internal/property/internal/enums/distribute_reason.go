package enums

type DistributionReason uint32

const (
	DISTRIBUTION_REASON_UNSPECIFIED DistributionReason = 0
	DISTRIBUTION_REASON_FAKE        DistributionReason = 10
	DISTRIBUTION_REASON_CHANGE      DistributionReason = 20
	DISTRIBUTION_REASON_HAVE_OTHER  DistributionReason = 30
	DISTRIBUTION_REASON_OTHER       DistributionReason = 40
)

var DistributionReasonNameMap = map[DistributionReason]string{
	// DISTRIBUTION_REASON_UNSPECIFIED: "Không xác định",
	DISTRIBUTION_REASON_FAKE:       "Sai thông tin / chọn nhầm đối tác",
	DISTRIBUTION_REASON_CHANGE:     "Chủ nhà đổi phương án / đổi giá",
	DISTRIBUTION_REASON_HAVE_OTHER: "Đã có đối tác khác phụ trách",
	DISTRIBUTION_REASON_OTHER:      "Khác",
}

func (r DistributionReason) String() string {
	if name, ok := DistributionReasonNameMap[r]; ok {
		return name
	}
	return "Không xác định"
}
