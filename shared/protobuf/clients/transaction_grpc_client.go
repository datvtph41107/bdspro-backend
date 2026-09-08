package clients

import (
	"context"

	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"

	"google.golang.org/grpc"
)

type TransactionGrpcClient struct {
	client transactionpb.InternalTransactionServiceClient
}

func BindTransactionClient(conn grpc.ClientConnInterface) *TransactionGrpcClient {
	return &TransactionGrpcClient{client: transactionpb.NewInternalTransactionServiceClient(conn)}
}

func (c *TransactionGrpcClient) GetTransactionDetail(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.TransactionResponse, error) {
	response, err := c.client.GetTransactionDetail(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *TransactionGrpcClient) GetTransaction(ctx context.Context, req *transactionpb.GetTransactionRequest) (*transactionpb.TransactionListResponse, error) {
	response, err := c.client.GetTransaction(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *TransactionGrpcClient) CreateAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	response, err := c.client.CreateAction(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *TransactionGrpcClient) CreateProductTransaction(ctx context.Context, req *transactionpb.CreateProductTransactionRequest) (*transactionpb.TransactionResponse, error) {
	response, err := c.client.CreateProductTransaction(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *TransactionGrpcClient) GetTransactionByIds(ctx context.Context, req []uint64) (*transactionpb.TransactionListResponse, error) {
	response, err := c.client.GetTransactionByIDs(ctx, &sharepb.IdRequest{Ids: req})
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *TransactionGrpcClient) GetDealStatisticByDealIds(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.StatisticDealsResponse, error) {
	response, err := c.client.GetDealStatisticByDealIds(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
