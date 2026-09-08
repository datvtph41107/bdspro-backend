package domain

type CommissionStats struct {
	TotalCommission      float64 `json:"total_commission"`
	TotalMembers         uint32  `json:"total_members"`
	TotalAmountCommitted float64 `json:"total_amount_committed"`

	DealID               uint64  `json:"dealId"`
	TargetProfit         float64 `json:"targetProfit"`         // Lợi nhuận mục tiêu
	MemberCount          uint32  `json:"memberCount"`          // Số người đã chia
	RemainingProfit      float64 `json:"remainingProfit"`      // Lợi nhuận còn lại
	CommissionPercentage float64 `json:"commissionPercentage"` // Phần trăm đã chia
	// TotalCommission      float64 `json:"totalCommission"`      // Tổng tiền đã chia
}
