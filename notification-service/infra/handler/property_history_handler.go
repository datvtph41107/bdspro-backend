package handler

import (
	"context"

	_utils "common/utils"
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
	SyncProvider *_utils.SyncUtil
}

func NewPropertyHistoryHandler(
	uc *usecase.PropertyHistoryUseCase,
	SyncProvider *_utils.SyncUtil,
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
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyPropertyHistories, req.SubjectId)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	// client mặc định sẽ truyền page = 0 GET /v2/deals?page=0
	if !updated && req.Page == 0 {
		return &notificationpb.GetPropertyHistoryResponse{}, nil
	}
	query := h.mapper.SearchProtoToDTO(req)
	histories, userMap, total, err := h.usecase.Search(ctx, query)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	// if req.Page == 0 && len(histories) > 0 {
	// 	time := time.Now().UnixMilli() - 24*60*60*1000
	// 	if len(histories) > 0 {
	// 		time = histories[0].UpdatedAt.UnixMilli()
	// 	}
	// 	h.SyncProvider.PutTimestamp(ctx, key, time)
	// }

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
