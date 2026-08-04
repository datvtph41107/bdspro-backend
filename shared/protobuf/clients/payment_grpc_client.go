package clients

import (
	"common"
	_middleware "common/middleware"
	"context"
	"errors"
	"fmt"
	"log"
	paymentpb "pb/types/payment"
	sharepb "pb/types/shared"
	"time"

	"google.golang.org/grpc"
)

type PaymentClient struct {
	client paymentpb.InternalServiceClient
}

func MakePaymentClient(env_runtime string) *PaymentClient {
	paymentClient := &PaymentClient{}
	go func() {
		name := "payment"
		url := common.GetPrefixProtobufUrl(env_runtime, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8205", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(_middleware.StreamClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Khởi chạy goroutine để theo dõi trạng thái kết nối
		go _middleware.MonitorConnection(ctx, name, conn)

		client := paymentpb.NewInternalServiceClient(conn)
		paymentClient.client = client
	}()
	return paymentClient
}

func (c *PaymentClient) GetBankById(ctx context.Context, id uint64) (*paymentpb.Bank, error) {
	if c.client == nil {
		return nil, errors.New("payment client not avaiable")
	}

	bank, err := c.client.GetBankById(ctx, &sharepb.IdRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return bank, nil
}
