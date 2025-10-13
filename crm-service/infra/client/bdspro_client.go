package client

import (
	_middleware "common/middleware"
	"context"
	"crm/env"
	"fmt"
	"log"
	bdspropb "pb/types/bdspro"
	crmpb "pb/types/crm"
	"time"

	"google.golang.org/grpc"
)

type BdsproClient struct {
	Client bdspropb.ProductServiceClient
}

func NewBdsproClient() *BdsproClient {
	bdsproClient := &BdsproClient{Client: nil}
	go func() {
		url := "localhost"
		name := "bdspro"
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

		client := bdspropb.NewProductServiceClient(conn)
		bdsproClient.Client = client
	}()
	return bdsproClient
}

func (c *BdsproClient) GetProductByIds(ctx context.Context, ids []uint64) ([]*bdspropb.ProductResponse, error) {
	request := &bdspropb.IdRequest{
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
