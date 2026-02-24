package clients

import (
	"common"
	"context"
	"fmt"
	"log"
	hubpb "pb/types/hub"
	sharepb "pb/types/shared"
	"time"

	_dto "common/domain/dto"
	_errors "common/errors"
	_middleware "common/middleware"

	"google.golang.org/grpc"
)

type HubGrpcClient struct {
	Client           hubpb.LocationServiceClient
	InternalClient   hubpb.HubInternalServiceClient
	LocationV2Client hubpb.LocationV2ServiceClient
	ApiKeyClient     hubpb.ApiKeyServiceClient
	UserGuideClient  hubpb.UserGuideServiceClient
}

func MakeHubGrpcClient(env_runtime string) *HubGrpcClient {
	hubClient := &HubGrpcClient{Client: nil, InternalClient: nil}
	go func() {
		name := "hub"
		url := common.GetPrefixProtobufUrl(env_runtime, name)

		conn, err := grpc.Dial(fmt.Sprintf("%s:8280", url),
			grpc.WithInsecure(),
			grpc.WithBackoffMaxDelay(5*time.Second),
			grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
		)
		if err != nil {
			log.Fatalf("failed to dial hub service: %v", err)
		}

		ctx, _ := context.WithCancel(context.Background())

		// Monitor connection status
		go _middleware.MonitorConnection(ctx, name, conn)

		hubClient.Client = hubpb.NewLocationServiceClient(conn)
		hubClient.InternalClient = hubpb.NewHubInternalServiceClient(conn)
		hubClient.LocationV2Client = hubpb.NewLocationV2ServiceClient(conn)
		hubClient.ApiKeyClient = hubpb.NewApiKeyServiceClient(conn)
		hubClient.UserGuideClient = hubpb.NewUserGuideServiceClient(conn)
	}()
	return hubClient
}

func (c *HubGrpcClient) InferAddressFromText(ctx context.Context, text string) (*_dto.AddressV3DTO, error) {
	if c == nil || c.InternalClient == nil {
		return nil, _errors.ReturnError(500, "hub client not avaiable")
	}

	addressV3, err := c.InternalClient.InferAddressFromText(ctx, &sharepb.RequestV3Proto{Text: text})
	if err != nil {
		return nil, err
	}
	if addressV3 == nil {
		return nil, _errors.ReturnError(500, "failed to infer address from text")
	}
	return &_dto.AddressV3DTO{
		Detail:       addressV3.Detail,
		ProvinceID:   addressV3.ProvinceId,
		ProvinceName: addressV3.ProvinceName,
		DistrictID:   addressV3.DistrictId,
		DistrictName: addressV3.DistrictName,
		WardID:       addressV3.WardId,
		WardName:     addressV3.WardName,
	}, nil
}

func (c *HubGrpcClient) VerifyApiKey(ctx context.Context, apiKey string) (*hubpb.VerifyApiKeyResponse, error) {
	if c == nil || c.ApiKeyClient == nil {
		return nil, fmt.Errorf("hub api key client not available")
	}

	return c.ApiKeyClient.VerifyApiKey(ctx, &hubpb.VerifyApiKeyRequest{ApiKey: apiKey})
}

func (c *HubGrpcClient) PutUpdate(ctx context.Context, resource string, resourceId uint64, ownerIds []uint64, updatedAt int64) error {
	_, err := c.InternalClient.PutUpdate(ctx, &hubpb.PutUpdateRequest{
		Resource:  resource,
		Id:        resourceId,
		OwnerIds:  ownerIds,
		UpdatedAt: updatedAt,
	})
	return err
}

func (c *HubGrpcClient) DelUpdate(ctx context.Context, resource string, id uint64, ownerId uint64) error {
	return nil
}
