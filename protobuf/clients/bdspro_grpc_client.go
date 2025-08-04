package clients

import (
	_middleware "common/middleware"
	"context"
	"errors"
	"fmt"
	"log"
	pb_bdspro "pb/types/bdspro"
	"time"

	"google.golang.org/grpc"
)

type BdsproGrpcClient struct {
	Client pb_bdspro.ProductServiceClient
}

func MakeBdsproGrpcClient(env_runtime string) *BdsproGrpcClient {
	bdsproClient := &BdsproGrpcClient{Client: nil}
	go func() {
		url := "localhost"
		name := "bdspro_grpc"
		if env_runtime == "release" {
			url = name
		}

		conn, err := grpc.Dial(fmt.Sprintf("%s:8202", url),
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

		client := pb_bdspro.NewProductServiceClient(conn)
		bdsproClient.Client = client
	}()
	return bdsproClient
}

func (c *BdsproGrpcClient) GetProductByIds(ctx context.Context, ids []uint64) (*pb_bdspro.ProductInternalResponse, error) {
	if c.Client == nil {
		return nil, errors.New("bdspro client not avaiable")
	}

	response, err := c.Client.GetDetailByIds(ctx, &pb_bdspro.IdRequest{Ids: ids})
	if err != nil {
		return nil, err
	}

	return response, nil
}