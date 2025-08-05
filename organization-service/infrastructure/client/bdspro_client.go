package client

import (
	_middleware "common/middleware"
	"context"
	"fmt"
	"log"
	"organization/env"
	bdspropb "pb/types/bdspro"
	organizationpb "pb/types/organization"
	"time"

	"google.golang.org/grpc"
)

// type BdsproClient interface {
// 	GetProductByIds(ctx context.Context, in *bdspropb.IdRequest) (*bdspropb.ProductInternalResponse, error)
// }

type BdsproClient struct {
	client bdspropb.ProductServiceClient
}

func NewBdsproClient() *BdsproClient {
	bdsproClient := &BdsproClient{client: nil}
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

		client := bdspropb.NewProductServiceClient(conn)
		bdsproClient.client = client
	}()
	return bdsproClient
}

func (c *BdsproClient) MapProductToDealPb(ctx context.Context, deals []*organizationpb.GroupDeal) {
	ids := make([]uint64, len(deals))
	for _, deal := range deals {
		if deal.Products != nil {
			for _, product := range deal.Products {
				ids = append(ids, product.Id)
			}
		}
	}

	products, err := c.client.GetProductAttachmentByIds(ctx, &bdspropb.IdRequest{Ids: ids})
	if err != nil {
		return
	}

	mapProduct := make(map[uint64]*bdspropb.ProductAttachment)
	for _, product := range products.Data {
		mapProduct[product.Id] = product
	}

	for _, deal := range deals {
		if deal.Products != nil {
			for i, product := range deal.Products {
				deal.Products[i] = mapProduct[product.Id]
			}
		}
	}
}

func (c *BdsproClient) GetProductAttachmentByIds(ctx context.Context, ids []uint64) ([]*bdspropb.ProductAttachment, error) {
	products, err := c.client.GetProductAttachmentByIds(ctx, &bdspropb.IdRequest{Ids: ids})
	if err != nil {
		return nil, err
	}
	return products.Data, nil
}
