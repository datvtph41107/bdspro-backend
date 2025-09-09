package dto

// DealOverview represents the overview statistics for deals
type DealOverview struct {
	TotalDealProcessing uint32  `gorm:"column:total_deal_processing" json:"totalDealProcessing"` // Tổng số thương vụ
	TotalDealCompleted  uint32  `gorm:"column:total_deal_completed" json:"totalDealCompleted"`   // Hoàn tất (theo status)
	TotalDealCanceled   uint32  `gorm:"column:total_deal_canceled" json:"totalDealCanceled"`     // Hủy (theo status)
	TotalCapital        float64 `gorm:"column:total_capital" json:"totalCapital"`                // Tổng vốn (sum của targetProfit)
	EstimatedProfit     float64 `gorm:"column:estimated_profit" json:"estimatedProfit"`          // Lợi nhuận ước tính (từ transaction service)
	ProcessingDeals     uint32  `json:"processingDeals"`                                         // Đang xử lý (theo status)
}
