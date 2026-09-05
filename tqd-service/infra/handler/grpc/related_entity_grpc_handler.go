package handler_grpc

import (
	"context"
	"errors"

	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	relateddomain "tqd/internal/domain/relatedentity/model"
	relatedapplication "tqd/internal/usecase/relatedentity/application"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RelatedEntityGrpcHandler is the transport adapter for bounded entity
// relationships. It performs transport validation, domain/protobuf mapping and
// error mapping only. Spatial querying, relationship classification, ranking
// and deduplication are owned by tqd-service application/domain layers.
type RelatedEntityGrpcHandler struct {
	tqdpb.UnimplementedRelatedEntityServiceServer

	service         *relatedapplication.Service
	discoveryMapper *mapper.DiscoveryMapper
}

func NewRelatedEntityGrpcHandler(
	service *relatedapplication.Service,
	discoveryMapper *mapper.DiscoveryMapper,
) *RelatedEntityGrpcHandler {
	return &RelatedEntityGrpcHandler{
		service:         service,
		discoveryMapper: discoveryMapper,
	}
}

func (h *RelatedEntityGrpcHandler) GetRelatedEntities(
	ctx context.Context,
	input *tqdpb.GetRelatedEntitiesRequest,
) (*tqdpb.GetRelatedEntitiesResponse, error) {
	if input == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	response, err := h.service.GetRelatedEntities(ctx, relateddomain.Request{
		Key:          input.GetKey(),
		RadiusMeters: input.GetRadiusMeters(),
		Limit:        int(input.GetLimit()),
	})
	if err != nil {
		return nil, relatedEntityError(err)
	}

	result := &tqdpb.GetRelatedEntitiesResponse{
		RequestId:           response.RequestID,
		Partial:             response.Partial,
		Warnings:            append([]string(nil), response.Warnings...),
		AppliedRadiusMeters: response.AppliedRadiusMeters,
		AppliedLimit:        int32(response.AppliedLimit),
		Candidates:          make([]*tqdpb.EntityCandidate, 0, len(response.Candidates)),
	}

	if response.Primary != nil {
		result.Primary, err = h.discoveryMapper.ToProtoEntityCandidate(response.Primary)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "map related primary: %v", err)
		}
	}

	for index := range response.Candidates {
		candidate, mapErr := h.discoveryMapper.ToProtoEntityCandidate(&response.Candidates[index].Entity)
		if mapErr != nil {
			return nil, status.Errorf(codes.Internal, "map related candidate: %v", mapErr)
		}
		result.Candidates = append(result.Candidates, candidate)
	}

	return result, nil
}

func relatedEntityError(err error) error {
	if err == nil {
		return nil
	}

	var validationError *relateddomain.ValidationError
	switch {
	case errors.As(err, &validationError):
		return status.Error(codes.InvalidArgument, validationError.Error())
	case errors.Is(err, relateddomain.ErrEntityNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, relateddomain.ErrProvidersFailed):
		return status.Error(codes.Unavailable, err.Error())
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, err.Error())
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, err.Error())
	default:
		return status.Error(codes.Internal, "get related entities failed")
	}
}
