package client

import (
	"context"
	"organization/env"
	"pb/clients"
	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"

	iusecase "organization/internal/interface"
)

// @bind: organization/internal/interface.TransactionClient
type TransactionClient struct {
	TransactionGrpcClient *clients.TransactionGrpcClient
}

func NewTransactionClient() *TransactionClient {
	return &TransactionClient{
		TransactionGrpcClient: clients.NewTransactionGrpcClient(env.ENV_RUNTIME),
	}
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

func (c *TransactionClient) GetTotalAmountByDealID(ctx context.Context, dealID uint64) (*iusecase.GetTotalAmountByDealIDResponse, error) {
	response, err := c.TransactionGrpcClient.GetTotalAmountByDealID(ctx, &sharepb.IdRequest{Id: dealID})
	if err != nil {
		return nil, err
	}
	return &iusecase.GetTotalAmountByDealIDResponse{
		DealID:         response.DealId,
		TotalPayment:   response.TotalPayment,
		NumTransaction: response.NumTransaction,
	}, nil
}
