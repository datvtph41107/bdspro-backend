package clients

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"time"

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

func (c *TransactionGrpcClient) CreateSignTransaction(ctx context.Context, req *transactionpb.SignRequest) (*transactionpb.ActionResponse, error) {
	response, err := c.client.CreateSignTransaction(ctx, req)
	if err != nil {
		return nil, err
	}
	return response, nil
}
