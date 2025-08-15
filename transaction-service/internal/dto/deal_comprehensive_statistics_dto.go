package dto

type DealComprehensiveStatisticsDTO struct {
	DealID  uint64
	Payment *DealAmountResponse
	Cost    *DealCostStatisticsDTO
}

type DealComprehensiveStatisticsRequest struct {
	DealID  uint64 `json:"deal_id" validate:"required"`
	Payment bool   `json:"payment"`
	Cost    bool   `json:"cost"`
}
