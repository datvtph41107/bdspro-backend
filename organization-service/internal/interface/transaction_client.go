package iusecase

import (
	"context"
)

// TransactionClient defines the interface for transaction service client
type TransactionClient interface {
	// GetTotalAmountByDealID(ctx context.Context, dealId uint64) (*GetTotalAmountByDealIDResponse, error)
	GetEstimatedProfitByDealIDs(ctx context.Context, dealIDs []uint64) (float64, error)
}

// type GetTotalAmountByDealIDResponse struct {
// 	DealID         uint64  `json:"dealId"`
// 	TotalPayment   float64 `json:"totalPayment"`
// 	NumTransaction uint32  `json:"numTransaction"`
// }
