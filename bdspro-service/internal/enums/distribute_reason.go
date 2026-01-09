package enums

type DistributionReason uint32

const (
	DISTRIBUTION_REASON_UNSPECIFIED DistributionReason = 0
	DISTRIBUTION_REASON_PROMOTION   DistributionReason = 10
	DISTRIBUTION_REASON_PARTNER     DistributionReason = 20
	DISTRIBUTION_REASON_OTHER       DistributionReason = 30
)

var DistributionReasonNameMap = map[DistributionReason]string{
	DISTRIBUTION_REASON_UNSPECIFIED: "Unspecified",
	DISTRIBUTION_REASON_PROMOTION:   "Promotion",
	DISTRIBUTION_REASON_PARTNER:     "Partner",
	DISTRIBUTION_REASON_OTHER:       "Other",
}

func (r DistributionReason) String() string {
	if name, ok := DistributionReasonNameMap[r]; ok {
		return name
	}
	return "Unknown"
}
