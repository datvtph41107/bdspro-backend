package entity

// CommissionStats chứa thống kê về hoa hồng của deal
type CommissionStats struct {
	DealID              uint64  `json:"dealId"`
	TargetProfit        float64 `json:"targetProfit"`        // Lợi nhuận mục tiêu
	TotalCommission     float64 `json:"totalCommission"`     // Tổng tiền đã chia
	MemberCount         uint32  `json:"memberCount"`         // Số người đã chia
	RemainingProfit     float64 `json:"remainingProfit"`     // Lợi nhuận còn lại
	CommissionPercentage float64 `json:"commissionPercentage"` // Phần trăm đã chia
} 