package dto

type TxDealCostStatisticsDTO struct {
	DealID        uint64
	Payment       *TxStats
	AmountCost    *TxStats
	AmountRevenue *TxStats
	NumPayment    *TxStats
	NumCost       *TxStats
	NumRevenue    *TxStats
	TotalAmount   float64
	TotalCost     float64
	TotalProfit   float64
}

type TxStatisticDealsDTO struct {
	DealIDs       []uint64
	AmountTxt     float64
	AmountRevenue float64
	AmountCost    float64
	TotalProfit   float64
}
