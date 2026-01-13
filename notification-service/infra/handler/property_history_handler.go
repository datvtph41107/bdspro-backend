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
}

func NewPropertyHistoryHandler(
	uc *usecase.PropertyHistoryUseCase,
) *PropertyHistoryHandler {
	return &PropertyHistoryHandler{
		usecase: uc,
	}
}

func (h *PropertyHistoryHandler) CreatePropertyHistory(
	ctx context.Context,
	req *notificationpb.CreatePropertyHistoryRequest,
) (*notificationpb.CreatePropertyHistoryResponse, error) {

	entity, err := mapper.RequestToPrototypeHistoryEntity(req)

	result, err := h.usecase.Create(ctx, entity)
	if err != nil {
		return &notificationpb.CreatePropertyHistoryResponse{
			Success: false,
			// Message: err.Error(),
		}, err
	}

	return &notificationpb.CreatePropertyHistoryResponse{
		Success: true,
		// Message: "success",
		Id: result.ID,
	}, nil
}

func (h *PropertyHistoryHandler) GetPropertyHistory(ctx context.Context, req *notificationpb.GetPropertyHistoryRequest) (*notificationpb.GetPropertyHistoryResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetActivityHistory not implemented")
}
func (h *PropertyHistoryHandler) GetPropertyHistoryByID(ctx context.Context, req *notificationpb.GetPropertyHistoryByIDRequest) (*notificationpb.GetPropertyHistoryByIDResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetActivityHistoryByID not implemented")
}
