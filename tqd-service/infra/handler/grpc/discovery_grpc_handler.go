package handler_grpc

import (
	"context"
	"strings"

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

	service *application.Service
	mapper  *mapper.DiscoveryMapper
}

// NewDiscoveryGrpcHandler khởi tạo Discovery gRPC handler.
func NewDiscoveryGrpcHandler(
	service *application.Service,
	discoveryMapper *mapper.DiscoveryMapper,
) *DiscoveryGrpcHandler {
	return &DiscoveryGrpcHandler{
		service: service,
		mapper:  discoveryMapper,
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
		return nil, discoveryError(err, "identify")
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

	response, err := h.service.Search(ctx, request)
	if err != nil {
		return nil, discoveryError(err, "search")
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
		return nil, discoveryError(err, "get entity")
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

// discoveryError chuyển application/domain error thành gRPC status error.
//
// Đây mới chỉ là error mapping theo thông điệp hiện có.
// Về lâu dài nên thay bằng các typed domain error, chẳng hạn:
//
//   - domain.ErrInvalidPoint
//   - domain.ErrEntityNotFound
//   - domain.ErrAllProvidersFailed
//
// Khi đó không cần kiểm tra chuỗi error.
func discoveryError(err error, operation string) error {
	if err == nil {
		return nil
	}

	message := strings.ToLower(err.Error())

	switch {
	case strings.Contains(message, "invalid latitude"),
		strings.Contains(message, "invalid longitude"),
		strings.Contains(message, "invalid canonical"),
		strings.Contains(message, "entity id must be numeric"),
		strings.Contains(message, "administrative entity id must be"),
		strings.Contains(message, "unsupported administrative unit type"),
		strings.Contains(message, "unsupported entity kind"),
		strings.Contains(message, "query is required"):
		return status.Error(
			codes.InvalidArgument,
			err.Error(),
		)

	case strings.Contains(message, "not found"):
		return status.Error(
			codes.NotFound,
			err.Error(),
		)

	case strings.Contains(message, "all discovery providers failed"),
		strings.Contains(message, "all search providers failed"):
		return status.Error(
			codes.Unavailable,
			err.Error(),
		)

	default:
		return status.Errorf(
			codes.Internal,
			"%s failed: %v",
			operation,
			err,
		)
	}
}
