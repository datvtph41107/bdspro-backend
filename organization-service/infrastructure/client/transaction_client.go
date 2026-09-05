package client

import (
	"context"
	"pb/clients"
	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"
)

// @bind: organization/internal/interface.TransactionClient
type TransactionClient struct {
	TransactionGrpcClient *clients.TransactionGrpcClient
}

func NewTransactionClient(rpcClient *clients.TransactionGrpcClient) *TransactionClient {
	return &TransactionClient{TransactionGrpcClient: rpcClient}
}

func (c *TransactionClient) MapTransactionToDealContract(ctx context.Context, contracts []*transactionpb.DealContractResponse) {
	ids := make([]uint64, len(contracts))
	for i, contract := range contracts {
		ids[i] = contract.TransactionId
	}
	response, err := c.TransactionGrpcClient.GetTransactionByIds(ctx, ids)
	if err != nil {
		return
	}
	transactionMap := make(map[uint64]*transactionpb.TransactionResponse)
	for _, transaction := range response.Data {
		transactionMap[transaction.TransactionId] = transaction
	}
	for _, contract := range contracts {
		if tx, ok := transactionMap[contract.TransactionId]; ok {
			contract.Transaction = tx
			contract.TransactionName = tx.TransactionName
			contract.Status = uint64(tx.Status)
			// contract.StatusName = tx.StatusName
		}
	}
}

// func (c *TransactionClient) GetTotalAmountByDealID(ctx context.Context, dealID uint64) (*iusecase.GetTotalAmountByDealIDResponse, error) {
// 	response, err := c.TransactionGrpcClient.GetTotalAmountByDealID(ctx, &sharepb.IdRequest{Id: dealID})
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &iusecase.GetTotalAmountByDealIDResponse{
// 		DealID:         response.DealId,
// 		TotalPayment:   response.TotalPayment,
// 		NumTransaction: response.NumTransaction,
// 	}, nil
// }

// GetEstimatedProfitByDealIDs lấy lợi nhuận ước tính theo danh sách deal IDs
func (c *TransactionClient) GetEstimatedProfitByDealIDs(ctx context.Context, dealIDs []uint64) (float64, error) {
	if len(dealIDs) == 0 {
		return 0, nil
	}
	response, err := c.TransactionGrpcClient.GetDealStatisticByDealIds(ctx, &sharepb.IdRequest{Ids: dealIDs})
	if err != nil {
		return 0, err
	}

	return response.TotalProfit, nil
}
