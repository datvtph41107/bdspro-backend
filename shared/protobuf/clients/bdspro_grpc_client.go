package clients

import (
	_middleware "common/middleware"
	"context"
	"errors"
	"fmt"
	"log"
	bdspropb "pb/types/bdspro"
	"time"

	"google.golang.org/grpc"
)

type BdsproGrpcClient struct {
	Client         bdspropb.ProductServiceClient
	InternalClient bdspropb.BdsproInternalServiceClient
}

func MakeBdsproGrpcClient(env_runtime string) *BdsproGrpcClient {
	bdsproClient := &BdsproGrpcClient{Client: nil, InternalClient: nil}
	go func() {
		url := "localhost"
		name := "bdspro"
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

		client := bdspropb.NewProductServiceClient(conn)
		internalClient := bdspropb.NewBdsproInternalServiceClient(conn)
		bdsproClient.Client = client
		bdsproClient.InternalClient = internalClient
	}()
	return bdsproClient
}

func (c *BdsproGrpcClient) GetProductByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductInternalResponse, error) {
	if c.Client == nil {
		return nil, errors.New("bdspro client not avaiable")
	}

	response, err := c.Client.GetDetailByIds(ctx, &bdspropb.IdRequest{Ids: ids})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *BdsproGrpcClient) GetProductAttachmentByIds(ctx context.Context, ids []uint64) (*bdspropb.ProductAttachmentResponse, error) {
	if c.Client == nil {
		return nil, errors.New("bdspro client not avaiable")
	}

	response, err := c.Client.GetProductAttachmentByIds(ctx, &bdspropb.IdRequest{Ids: ids})
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *BdsproGrpcClient) GetCountByOwner(ctx context.Context, req *bdspropb.GetCountByOwnerRequest) (*bdspropb.GetCountByOwnerResponse, error) {
	if c.InternalClient == nil {
		return nil, errors.New("bdspro internal client not available")
	}

	response, err := c.InternalClient.GetCountByOwner(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}
