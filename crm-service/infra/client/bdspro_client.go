package client

import (
	_middleware "common/middleware"
	"context"
	"crm/env"
	"fmt"
	"log"
	pb_bdspro "pb/types/bdspro"
	crmpb "pb/types/crm"
	"time"

	"google.golang.org/grpc"
)

type BdsproClient struct {
	Client pb_bdspro.ProductServiceClient
}

func NewBdsproClient() *BdsproClient {
	bdsproClient := &BdsproClient{Client: nil}
	go func() {
		url := "localhost"
		name := "bdspro_grpc"
		if env.ENV_RUNTIME == "release" {
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

func (c *BdsproClient) GetProductByIds(ctx context.Context, ids []uint64) ([]*pb_bdspro.ProductResponse, error) {
	request := &pb_bdspro.IdRequest{
		Ids: ids,
	}

	response, err := c.Client.GetDetailByIds(ctx, request)
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

func (c *BdsproClient) MapProductPbByIds(ctx context.Context, lead *crmpb.LeadDetailDTO) {
	products, err := c.GetProductByIds(ctx, lead.ProductIds)
	if err != nil {
		return
	}

	lead.Products = products
}
