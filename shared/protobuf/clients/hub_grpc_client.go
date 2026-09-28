package clients

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	hubpb "pb/types/hub"
	sharepb "pb/types/shared"

	_dto "common/domain/dto"
	_errors "common/errors"

	"google.golang.org/grpc"
)

type HubGrpcClient struct {
	Client           hubpb.LocationServiceClient
	InternalClient   hubpb.HubInternalServiceClient
	LocationV2Client hubpb.LocationV2ServiceClient
	ApiKeyClient     hubpb.ApiKeyServiceClient
	UserGuideClient  hubpb.UserGuideServiceClient
	FAQClient        hubpb.FAQServiceClient
	VersionClient    hubpb.VersionServiceClient
}

func BindHubGrpcClient(conn grpc.ClientConnInterface) *HubGrpcClient {
	return &HubGrpcClient{
		Client:           hubpb.NewLocationServiceClient(conn),
		InternalClient:   hubpb.NewHubInternalServiceClient(conn),
		LocationV2Client: hubpb.NewLocationV2ServiceClient(conn),
		ApiKeyClient:     hubpb.NewApiKeyServiceClient(conn),
		UserGuideClient:  hubpb.NewUserGuideServiceClient(conn),
		FAQClient:        hubpb.NewFAQServiceClient(conn),
		VersionClient:    hubpb.NewVersionServiceClient(conn),
	}
}

func (c *HubGrpcClient) GetByIDs(ctx context.Context, req *hubpb.GetRegionsRequest) (*hubpb.GetRegionsResponse, error) {
	if c == nil || c.InternalClient == nil {
		return nil, _errors.ReturnError(500, "hub client not available")
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.InternalClient.GetByIDs(ctx, req)
	if err != nil {
		return nil, err
	}

	if resp == nil {
		return nil, _errors.ReturnError(500, "empty response from hub service")
	}

	return resp, nil
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

// TrackKeywordSearch fire-and-forget đếm từ khóa tìm kiếm sang hub (CRM/TQD/...).
func (c *HubGrpcClient) TrackKeywordSearch(ctx context.Context, keyword, source string) {
	if c == nil || c.InternalClient == nil {
		return
	}
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return
	}
	go func() {
		trackCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err := c.InternalClient.TrackKeywordSearch(trackCtx, &hubpb.TrackKeywordSearchRequest{
			Keyword: keyword,
			Source:  source,
			Count:   1,
		})
		if err != nil {
			log.Printf("hub TrackKeywordSearch error: %v", err)
		}
	}()
}

// TrackSearchLatency fire-and-forget ghi latency API search có text sang hub.
func (c *HubGrpcClient) TrackSearchLatency(ctx context.Context, method, path, service string, durationMs int64) {
	if c == nil || c.InternalClient == nil {
		return
	}
	method = strings.TrimSpace(method)
	path = strings.TrimSpace(path)
	if method == "" || path == "" || durationMs < 0 {
		return
	}
	go func() {
		trackCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err := c.InternalClient.TrackSearchLatency(trackCtx, &hubpb.TrackSearchLatencyRequest{
			Method:     method,
			Path:       path,
			Service:    service,
			DurationMs: durationMs,
		})
		if err != nil {
			log.Printf("hub TrackSearchLatency error: %v", err)
		}
	}()
}
