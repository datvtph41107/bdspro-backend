package client

import (
	"common"
	_middleware "common/middleware"
	"context"
	"crm/env"
	"crm/internal/dto"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	tqdpb "pb/types/tqd"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// @bind: crm/internal/interface/provider.TqdProvider
type TqdClient struct {
	Client              tqdpb.ParcelServiceClient
	SeoProjectionClient tqdpb.TqdSeoProjectionServiceClient
}

func NewTqdClient() *TqdClient {
	client := &TqdClient{}
	name := "tqd"
	serviceURL := common.GetPrefixProtobufUrl(env.ENV_RUNTIME, name)

	// grpc.Dial is non-blocking unless WithBlock is supplied, so constructing
	// both stubs synchronously removes the old nil-client race without delaying
	// CRM startup while the connection is still becoming READY.
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:8219", serviceURL),
		grpc.WithInsecure(),
		grpc.WithBackoffMaxDelay(5*time.Second),
		grpc.WithUnaryInterceptor(_middleware.UnaryClientInterceptor),
	)
	if err != nil {
		log.Printf("failed to dial tqd service: %v", err)
		return client
	}

	monitorCtx, cancel := context.WithCancel(context.Background())
	_ = cancel // connection lifetime follows the CRM process.
	go _middleware.MonitorConnection(monitorCtx, name, conn)

	client.Client = tqdpb.NewParcelServiceClient(conn)
	client.SeoProjectionClient = tqdpb.NewTqdSeoProjectionServiceClient(conn)
	return client
}

func (c *TqdClient) GetAdministrativeUnitProjection(
	ctx context.Context,
	identity string,
) (*dto.TqdAdministrativeUnitProjection, error) {
	if c == nil || c.SeoProjectionClient == nil {
		return nil, errors.New("tqd SEO projection client not available")
	}
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return nil, errors.New("administrative unit identity is required")
	}

	resp, err := c.SeoProjectionClient.GetAdministrativeUnitProjection(
		ctx,
		wrapperspb.String(identity),
	)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Value) == 0 {
		return nil, errors.New("tqd returned an empty administrative unit projection")
	}

	var projection dto.TqdAdministrativeUnitProjection
	if err := json.Unmarshal(resp.Value, &projection); err != nil {
		return nil, fmt.Errorf("decode administrative unit projection: %w", err)
	}
	if strings.TrimSpace(projection.EntityID) == "" ||
		strings.TrimSpace(projection.Name) == "" ||
		strings.TrimSpace(projection.AdministrativeCode) == "" ||
		!strings.HasPrefix(strings.TrimSpace(projection.CanonicalPath), "/dia-ban/") {
		return nil, errors.New("tqd administrative unit projection is missing stable identity fields")
	}
	if projection.SchemaVersion != dto.TqdAdministrativeUnitProjectionSchemaVersion {
		return nil, fmt.Errorf("unsupported administrative unit projection schema %q", projection.SchemaVersion)
	}
	return &projection, nil
}

func (c *TqdClient) GetPlanningProjectProjection(
	ctx context.Context,
	identity string,
) (*dto.TqdPlanningProjectProjection, error) {
	if c == nil || c.SeoProjectionClient == nil {
		return nil, errors.New("tqd SEO projection client not available")
	}
	identity = strings.TrimSpace(identity)
	if identity == "" {
		return nil, errors.New("planning project identity is required")
	}

	resp, err := c.SeoProjectionClient.GetPlanningProjectProjection(
		ctx,
		wrapperspb.String(identity),
	)
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.Value) == 0 {
		return nil, errors.New("tqd returned an empty planning project projection")
	}

	var projection dto.TqdPlanningProjectProjection
	if err := json.Unmarshal(resp.Value, &projection); err != nil {
		return nil, fmt.Errorf("decode planning project projection: %w", err)
	}
	if projection.ID == 0 || strings.TrimSpace(projection.Name) == "" {
		return nil, errors.New("tqd planning project projection is missing id or name")
	}
	return &projection, nil
}

func (c *TqdClient) GetParcelSeoSource(ctx context.Context, parcelID uint64) (*dto.ParcelSeoSource, error) {
	if c == nil || c.Client == nil {
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
	if c == nil || c.Client == nil {
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
	if c == nil || c.Client == nil {
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
