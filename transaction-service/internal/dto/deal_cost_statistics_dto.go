package dto

type DealCostStatisticsDTO struct {
	DealID        uint64
	Payment       *Stats
	AmountCost    *Stats
	AmountRevenue *Stats
	NumPayment    *Stats
	NumCost       *Stats
	NumRevenue    *Stats
	TotalAmount   float64
	TotalCost     float64
	TotalProfit   float64
}

type Stats struct {
	Approved float64
	Reject   float64
	Pending  float64
	Total    float64
}
