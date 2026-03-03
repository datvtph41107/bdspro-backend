package handler

import (
	"context"
	"fmt"

	_provider "common/provider"
	"notification/infra/mapper"
	"notification/internal/usecase"
	notificationpb "pb/types/notification"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PropertyHistoryHandler struct {
	notificationpb.UnimplementedPropertyHistoryServiceServer
	usecase      *usecase.PropertyHistoryUseCase
	mapper       *mapper.PropertyHistoryMapper
	SyncProvider *_provider.SyncProvider
}

func NewPropertyHistoryHandler(
	uc *usecase.PropertyHistoryUseCase,
	SyncProvider *_provider.SyncProvider,
) *PropertyHistoryHandler {
	return &PropertyHistoryHandler{
		usecase:      uc,
		mapper:       mapper.NewPropertyHistoryMapper(),
		SyncProvider: SyncProvider,
	}
}

func (h *PropertyHistoryHandler) CreatePropertyHistory(
	ctx context.Context,
	req *notificationpb.CreatePropertyHistoryRequest,
) (*notificationpb.CreatePropertyHistoryResponse, error) {

	input := h.mapper.CreateProtoToDTO(req)

	result, err := h.usecase.Create(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &notificationpb.CreatePropertyHistoryResponse{
		Success: true,
		Id:      result.ID,
	}, nil
}

func (h *PropertyHistoryHandler) SoftDeletePropertyHistoryId(
	ctx context.Context,
	req *notificationpb.SoftDeletePropertyHistoryIdRequest,
) (*notificationpb.SoftDeletePropertyHistoryIdResponse, error) {
	if err := h.usecase.SoftDelete(ctx, req.Id, req.ActorId); err != nil {
		return nil, err
	}

	return &notificationpb.SoftDeletePropertyHistoryIdResponse{
		Success: true,
	}, nil
}

func (h *PropertyHistoryHandler) GetPropertyHistory(
	ctx context.Context,
	req *notificationpb.GetPropertyHistoryRequest,
) (*notificationpb.GetPropertyHistoryResponse, error) {
	updated := h.SyncProvider.HasUpdated(ctx, fmt.Sprintf("time:product_property_histories:%d", req.SubjectId), req.TimeStamp)
	// client mặc định sẽ truyền page = 0 GET /v2/deals?page=0
	if !updated && req.Page == 0 {
		return &notificationpb.GetPropertyHistoryResponse{}, nil
	}
	query := h.mapper.SearchProtoToDTO(req)
	histories, userMap, total, err := h.usecase.Search(ctx, query)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if req.Page == 0 && len(histories) > 0 {
		h.SyncProvider.PutTimestamp(ctx, fmt.Sprintf("time:product_property_histories:%d", req.SubjectId), histories[0].UpdatedAt.UnixMilli())
	}

	return h.mapper.EntitiesToListResponse(
		histories,
		userMap,
		total,
		int32(query.Page),
		int32(query.Size),
	), nil
}

func (h *PropertyHistoryHandler) GetPropertyHistoryByID(
	ctx context.Context,
	req *notificationpb.GetPropertyHistoryByIDRequest,
) (*notificationpb.GetPropertyHistoryByIDResponse, error) {

	history, actor, err := h.usecase.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return h.mapper.EntityToDetailResponse(history, actor), nil
}
