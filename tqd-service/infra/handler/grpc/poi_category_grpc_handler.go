package handler_grpc

import (
	"context"
	"log"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PoiCategoryGrpcHandler implements gRPC handlers for PoiCategory service
type PoiCategoryGrpcHandler struct {
	tqdpb.UnimplementedPoiCategoryServiceServer
	poiCategoryUsecase *usecase.PoiCategoryUsecase
	poiCategoryMapper  *mapper.PoiCategoryMapper
}

// NewPoiCategoryGrpcHandler creates a new PoiCategoryGrpcHandler
func NewPoiCategoryGrpcHandler(
	poiCategoryUsecase *usecase.PoiCategoryUsecase,
	poiCategoryMapper *mapper.PoiCategoryMapper,
) *PoiCategoryGrpcHandler {
	return &PoiCategoryGrpcHandler{
		poiCategoryUsecase: poiCategoryUsecase,
		poiCategoryMapper:  poiCategoryMapper,
	}
}

// CreatePoiCategory handles gRPC CreatePoiCategory request
func (h *PoiCategoryGrpcHandler) CreatePoiCategory(ctx context.Context, req *tqdpb.PoiCategory) (*tqdpb.PoiCategory, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "poi category name is required")
	}

	createReq := &dto.CreatePoiCategoryRequestDTO{
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		Icon:        req.Icon,
		Color:       req.Color,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		ParentId:    req.ParentId,
	}

	result, err := h.poiCategoryUsecase.CreatePoiCategory(ctx, createReq)
	if err != nil {
		log.Printf("error creating poi category: %v", err)
		return nil, status.Error(codes.Internal, "failed to create poi category")
	}

	return h.poiCategoryMapper.ToProto(result), nil
}

// GetPoiCategory handles gRPC GetPoiCategory request
func (h *PoiCategoryGrpcHandler) GetPoiCategory(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.PoiCategory, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi category id is required")
	}

	result, err := h.poiCategoryUsecase.GetPoiCategoryByID(ctx, req.Id)
	if err != nil {
		log.Printf("error getting poi category: %v", err)
		return nil, status.Error(codes.NotFound, "poi category not found")
	}

	return h.poiCategoryMapper.ToProto(result), nil
}

// UpdatePoiCategory handles gRPC UpdatePoiCategory request
func (h *PoiCategoryGrpcHandler) UpdatePoiCategory(ctx context.Context, req *tqdpb.PoiCategory) (*tqdpb.PoiCategory, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi category id is required")
	}

	updateReq := &dto.PoiCategoryDTO{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		Icon:        req.Icon,
		Color:       req.Color,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		ParentId:    req.ParentId,
	}

	result, err := h.poiCategoryUsecase.UpdatePoiCategory(ctx, updateReq)
	if err != nil {
		log.Printf("error updating poi category: %v", err)
		return nil, status.Error(codes.Internal, "failed to update poi category")
	}

	return h.poiCategoryMapper.ToProto(result), nil
}

// DeletePoiCategory handles gRPC DeletePoiCategory request
func (h *PoiCategoryGrpcHandler) DeletePoiCategory(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi category id is required")
	}

	err := h.poiCategoryUsecase.DeletePoiCategory(ctx, req.Id)
	if err != nil {
		log.Printf("error deleting poi category: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete poi category")
	}

	return &sharepb.SubmitResponse{
		Code:    200,
		Message: "POI category deleted successfully",
		Data:    "",
	}, nil
}

// ListPoiCategories handles gRPC ListPoiCategories request
func (h *PoiCategoryGrpcHandler) ListPoiCategories(ctx context.Context, req *tqdpb.ListPoiCategoriesRequest) (*tqdpb.ListPoiCategoriesResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	filters := &dto.PoiCategoryFilterDTO{
		Page:     int(req.Page),
		Size:     int(req.Size),
		Search:   req.Search,
		IsActive: req.IsActive,
		ParentId: req.ParentId,
		Level:    req.Level,
	}

	results, total, err := h.poiCategoryUsecase.ListPoiCategories(ctx, filters)
	if err != nil {
		log.Printf("error listing poi categories: %v", err)
		return nil, status.Error(codes.Internal, "failed to list poi categories")
	}

	categories := make([]*tqdpb.PoiCategory, len(results))
	for i, result := range results {
		categories[i] = h.poiCategoryMapper.ToProto(result)
	}

	return &tqdpb.ListPoiCategoriesResponse{
		Data:  categories,
		Total: int64(total),
	}, nil
}

// GetPoiCategoryTree handles gRPC GetPoiCategoryTree request
func (h *PoiCategoryGrpcHandler) GetPoiCategoryTree(ctx context.Context, req *emptypb.Empty) (*tqdpb.GetPoiCategoryTreeResponse, error) {
	results, err := h.poiCategoryUsecase.GetPoiCategoryTree(ctx)
	if err != nil {
		log.Printf("error getting poi category tree: %v", err)
		return nil, status.Error(codes.Internal, "failed to get poi category tree")
	}

	treeNodes := h.poiCategoryMapper.BuildTreeNodes(results)

	return &tqdpb.GetPoiCategoryTreeResponse{
		Data: treeNodes,
	}, nil
}
