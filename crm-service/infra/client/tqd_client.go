package client

import (
	"common"
	_middleware "common/middleware"
	"context"
	"crm/env"
	"crm/internal/dto"
	"errors"
	"fmt"
	"log"
	tqdpb "pb/types/tqd"
	"time"

	"google.golang.org/grpc"
)

// @bind: crm/internal/interface/provider.TqdProvider
type TqdClient struct {
	Client tqdpb.ParcelServiceClient
}

func NewTqdClient() *TqdClient {
	tqdClient := &TqdClient{Client: nil}
	go func() {
		name := "tqd"
		url := common.GetPrefixProtobufUrl(env.ENV_RUNTIME, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8219", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial tqd service: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go _middleware.MonitorConnection(ctx, name, conn)
		tqdClient.Client = tqdpb.NewParcelServiceClient(conn)
	}()
	return tqdClient
}

func (c *TqdClient) GetParcelSeoSource(ctx context.Context, parcelID uint64) (*dto.ParcelSeoSource, error) {
	if c.Client == nil {
		return nil, errors.New("tqd client not available")
	}

	resp, err := c.Client.GetParcelSeoSource(ctx, &tqdpb.GetParcelDetailRequest{
		ParcelId: parcelID,
	})
	if err != nil {
		return nil, err
	}
	return parcelSeoSourceFromProto(resp), nil
}

func (c *TqdClient) GetParcelSeoSourcesForGenerate(ctx context.Context, limit uint32) ([]dto.ParcelSeoSource, error) {
	if c.Client == nil {
		return nil, errors.New("tqd client not available")
	}

	resp, err := c.Client.ListParcelSeoSourcesForGenerate(ctx, &tqdpb.ListParcelSeoSourcesRequest{
		Limit: limit,
	})
	if err != nil {
		return nil, err
	}

	items := make([]dto.ParcelSeoSource, 0, len(resp.GetData()))
	for _, item := range resp.GetData() {
		if source := parcelSeoSourceFromProto(item); source != nil {
			items = append(items, *source)
		}
	}
	return items, nil
}

func (c *TqdClient) UpdateParcelSeoID(ctx context.Context, parcelID uint64, seoID uint64) error {
	if c.Client == nil {
		return errors.New("tqd client not available")
	}

	_, err := c.Client.UpdateParcelSeoID(ctx, &tqdpb.UpdateParcelSeoIDRequest{
		ParcelId: parcelID,
		SeoId:    seoID,
	})
	return err
}

func parcelSeoSourceFromProto(resp *tqdpb.ParcelSeoSourceResponse) *dto.ParcelSeoSource {
	if resp == nil || resp.ParcelId == 0 {
		return nil
	}
	source := &dto.ParcelSeoSource{
		ParcelID:  resp.ParcelId,
		AdrSearch: resp.AdrSearch,
	}
	if resp.SeoId > 0 {
		seoID := resp.SeoId
		source.SeoID = &seoID
	}
	return source
}