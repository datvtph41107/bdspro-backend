package dto

type TxDealComprehensiveStatisticsDTO struct {
	DealID  uint64
	Payment *TxDealAmountResponse
	Cost    *TxDealCostStatisticsDTO
}

type TxDealComprehensiveStatisticsRequest struct {
	DealID  uint64 `json:"deal_id" validate:"required"`
	Payment bool   `json:"payment"`
	Cost    bool   `json:"cost"`
}
