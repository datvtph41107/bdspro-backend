package clients

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"time"

	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"

	"google.golang.org/grpc"
)

type TransactionGrpcClient struct {
	client transactionpb.InternalTransactionServiceClient
}

func NewTransactionGrpcClient(env_runtime string) *TransactionGrpcClient {
	transactionClient := &TransactionGrpcClient{}
	go func() {
		url := "localhost"
		name := "transaction"
		if env_runtime == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8215", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client := transactionpb.NewInternalTransactionServiceClient(conn)
		transactionClient.client = client
	}()
	return transactionClient
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

func (c *TransactionGrpcClient) CreateDepositAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	response, err := c.client.CreateDepositAction(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *TransactionGrpcClient) CreateSignAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	response, err := c.client.CreateSignAction(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *TransactionGrpcClient) CreatePaymentAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	response, err := c.client.CreatePaymentAction(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *TransactionGrpcClient) CreateHandoverAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	response, err := c.client.CreateHandoverAction(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
