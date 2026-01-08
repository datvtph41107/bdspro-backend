package handler

import (
	"context"
	"notification/infra/mapper"
	"notification/internal/usecase"
	notificationpb "pb/types/notification"
)

type DealHistoryHandler struct {
	notificationpb.UnimplementedDealHistoryServiceServer
	dealHistoryUsecase *usecase.DealHistoryUsecase
}

func NewDealHistoryHandler(dealHistoryUsecase *usecase.DealHistoryUsecase) *DealHistoryHandler {
	return &DealHistoryHandler{
		dealHistoryUsecase: dealHistoryUsecase,
	}
}

// CreateDealHistory tạo deal history
func (s *DealHistoryHandler) CreateDealHistory(ctx context.Context, req *notificationpb.CreateDealHistoryRequest) (*notificationpb.CreateDealHistoryResponse, error) {
	// Convert request to domain request
	domainReq := mapper.RequestToDealHistoryRequest(req)

	// Create deal history
	response, err := s.dealHistoryUsecase.CreateDealHistory(ctx, domainReq)
	if err != nil {
		return &notificationpb.CreateDealHistoryResponse{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &notificationpb.CreateDealHistoryResponse{
		Success: response.Success,
		Message: response.Message,
		Id:      response.ID,
	}, nil
}

// GetDealHistory lấy danh sách deal history theo deal ID
func (s *DealHistoryHandler) GetDealHistory(ctx context.Context, req *notificationpb.GetDealHistoryRequest) (*notificationpb.GetDealHistoryResponse, error) {
	// Validate pagination parameters
	page := int(req.Page)
	size := int(req.Size)

	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20 // Default page size
	}

	// Get deal history
	entities, total, err := s.dealHistoryUsecase.GetDealHistory(ctx, req.DealId, page, size)
	if err != nil {
		return nil, err
	}

	// Convert to response
	response := mapper.EntitiesToDealHistoryResponse(entities, total, int32(page), int32(size))
	return response, nil
}

// GetDealHistoryByID lấy deal history theo ID
func (s *DealHistoryHandler) GetDealHistoryByID(ctx context.Context, req *notificationpb.GetDealHistoryByIDRequest) (*notificationpb.GetDealHistoryByIDResponse, error) {
	// Get deal history by ID
	entity, err := s.dealHistoryUsecase.GetDealHistoryByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if entity == nil {
		return &notificationpb.GetDealHistoryByIDResponse{}, nil
	}

	// Convert to response
	response := mapper.EntityToDealHistoryByIDResponse(entity)
	return response, nil
}
