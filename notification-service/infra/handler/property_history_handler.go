package handler

import (
	"context"

	"notification/infra/mapper"
	"notification/internal/usecase"
	notificationpb "pb/types/notification"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PropertyHistoryHandler struct {
	notificationpb.UnimplementedPropertyHistoryServiceServer
	usecase *usecase.PropertyHistoryUseCase
	mapper  *mapper.PropertyHistoryMapper
}

func NewPropertyHistoryHandler(
	uc *usecase.PropertyHistoryUseCase,
) *PropertyHistoryHandler {

	return &PropertyHistoryHandler{
		usecase: uc,
		mapper:  mapper.NewPropertyHistoryMapper(),
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

	query := h.mapper.SearchProtoToDTO(req)

	entities, total, err := h.usecase.Search(ctx, query)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.mapper.EntitiesToListResponse(
		entities,
		total,
		int32(query.Page),
		int32(query.Size),
	), nil
}

func (h *PropertyHistoryHandler) GetPropertyHistoryByID(
	ctx context.Context,
	req *notificationpb.GetPropertyHistoryByIDRequest,
) (*notificationpb.GetPropertyHistoryByIDResponse, error) {

	entity, err := h.usecase.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return h.mapper.EntityToDetailResponse(entity), nil
}
