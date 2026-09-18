package handler_grpc

import (
	"context"
	"strings"
	"time"
	"tqd/infra/handler/errorresponse"

	"pb/clients"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/usecase/discovery/application"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DiscoveryGrpcHandler là transport adapter của DiscoveryService.
//
// Handler chỉ chịu trách nhiệm:
//   - kiểm tra request ở mức transport;
//   - gọi DiscoveryMapper để chuyển protobuf sang domain;
//   - gọi application service;
//   - gọi DiscoveryMapper để chuyển domain sang protobuf;
//   - chuyển application error thành gRPC status.
//
// Handler không trực tiếp chứa logic mapping protobuf/domain.
type DiscoveryGrpcHandler struct {
	tqdpb.UnimplementedDiscoveryServiceServer

	service   *application.Service
	mapper    *mapper.DiscoveryMapper
	hubClient *clients.HubGrpcClient
}

// NewDiscoveryGrpcHandler khởi tạo Discovery gRPC handler.
func NewDiscoveryGrpcHandler(
	service *application.Service,
	discoveryMapper *mapper.DiscoveryMapper,
	hubClient *clients.HubGrpcClient,
) *DiscoveryGrpcHandler {
	return &DiscoveryGrpcHandler{
		service:   service,
		mapper:    discoveryMapper,
		hubClient: hubClient,
	}
}

// IdentifyAtPoint nhận tọa độ từ protobuf request và thực hiện nhận diện
// các đối tượng không gian tại vị trí đó.
func (h *DiscoveryGrpcHandler) IdentifyAtPoint(
	ctx context.Context,
	input *tqdpb.IdentifyAtPointRequest,
) (*tqdpb.IdentifyAtPointResponse, error) {
	if input == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"request is required",
		)
	}

	if input.GetPoint() == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"point is required",
		)
	}

	request := h.mapper.FromProtoIdentifyRequest(input)

	response, err := h.service.Identify(ctx, request)
	if err != nil {
		return nil, discoveryError(ctx, err, "identify")
	}

	result, err := h.mapper.ToProtoIdentifyResponse(response)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"map identify response: %v",
			err,
		)
	}

	return result, nil
}

// SearchDiscovery tìm kiếm các đối tượng Discovery thuộc phạm vi TQD.
func (h *DiscoveryGrpcHandler) SearchDiscovery(
	ctx context.Context,
	input *tqdpb.SearchDiscoveryRequest,
) (*tqdpb.SearchDiscoveryResponse, error) {
	if input == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"request is required",
		)
	}

	request := h.mapper.FromProtoSearchRequest(input)
	q := strings.TrimSpace(request.Query)

	start := time.Now()
	response, err := h.service.Search(ctx, request)
	durationMs := time.Since(start).Milliseconds()

	if q != "" && h.hubClient != nil {
		h.hubClient.TrackKeywordSearch(ctx, q, "tqd")
		h.hubClient.TrackSearchLatency(ctx, "GET", "/v2/tqd/discovery/search", "tqd", durationMs)
	}

	if err != nil {
		return nil, discoveryError(ctx, err, "search")
	}

	result, err := h.mapper.ToProtoSearchResponse(response)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"map search response: %v",
			err,
		)
	}

	return result, nil
}

// GetDiscoveryEntity phân giải một canonical entity key thành một
// EntityCandidate đầy đủ.
func (h *DiscoveryGrpcHandler) GetDiscoveryEntity(
	ctx context.Context,
	input *tqdpb.GetDiscoveryEntityRequest,
) (*tqdpb.EntityCandidate, error) {
	if input == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"request is required",
		)
	}

	key := strings.TrimSpace(input.GetKey())
	if key == "" {
		return nil, status.Error(
			codes.InvalidArgument,
			"canonical entity key is required",
		)
	}

	entity, err := h.service.GetEntity(ctx, key)
	if err != nil {
		return nil, discoveryError(ctx, err, "get entity")
	}

	if entity == nil {
		return nil, status.Error(
			codes.NotFound,
			"discovery entity not found",
		)
	}

	result, err := h.mapper.ToProtoEntityCandidate(entity)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			"map discovery entity: %v",
			err,
		)
	}

	return result, nil
}

// discoveryError giữ operation trong log; phân loại do adapter TQD sở hữu.
func discoveryError(ctx context.Context, err error, operation string) error {
	return errorresponse.GRPC(ctx, err)
}
