package handler_grpc

import (
	"context"
	"log"

	_dto "common/domain/dto"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// PoiCategoryGrpcHandler handles gRPC requests for POI category
type PoiCategoryGrpcHandler struct {
	tqdpb.UnimplementedPoiCategoryServiceServer
	usecase usecase.PoiCategoryUsecase
	mapper  *mapper.PoiCategoryMapper
}

// NewPoiCategoryGrpcHandler creates new handler
func NewPoiCategoryGrpcHandler(
	usecase usecase.PoiCategoryUsecase,
	mapper *mapper.PoiCategoryMapper,
) *PoiCategoryGrpcHandler {
	return &PoiCategoryGrpcHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

// CreatePoiCategory handles create request
func (h *PoiCategoryGrpcHandler) CreatePoiCategory(ctx context.Context, req *tqdpb.PoiCategory) (*tqdpb.PoiCategory, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "poi category name is required")
	}
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "poi category code is required")
	}

	// Convert proto to create request
	createReq, err := h.mapper.FromProto(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Call usecase
	result, err := h.usecase.Create(ctx, createReq)
	if err != nil {
		log.Printf("error creating poi category: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create poi category: %v", err)
	}

	return h.mapper.ToProtoFromResponse(result), nil
}

// GetPoiCategory handles get by ID request
func (h *PoiCategoryGrpcHandler) GetPoiCategory(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.PoiCategory, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi category id is required")
	}

	result, err := h.usecase.GetByID(ctx, req.Id)
	if err != nil {
		log.Printf("error getting poi category: %v", err)
		return nil, status.Errorf(codes.NotFound, "poi category not found: %v", err)
	}

	return h.mapper.ToProtoFromResponse(result), nil
}

// UpdatePoiCategory handles update request
func (h *PoiCategoryGrpcHandler) UpdatePoiCategory(ctx context.Context, req *tqdpb.PoiCategory) (*tqdpb.PoiCategory, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi category id is required")
	}

	// Build update request
	updateReq := &dto.UpdatePoiCategoryRequest{}

	if req.Code != "" {
		updateReq.Code = &req.Code
	}
	if req.Name != "" {
		updateReq.Name = &req.Name
	}
	if req.Description != "" {
		updateReq.Description = &req.Description
	}
	if req.Icon != "" {
		updateReq.Icon = &req.Icon
	}
	if req.Color != "" {
		updateReq.Color = &req.Color
	}
	if req.ParentId != 0 {
		updateReq.ParentID = &req.ParentId
	}
	if req.SortOrder != 0 {
		sortOrder := int32(req.SortOrder)
		updateReq.SortOrder = &sortOrder
	}
	updateReq.IsActive = &req.IsActive

	// Call usecase
	result, err := h.usecase.Update(ctx, req.Id, updateReq)
	if err != nil {
		log.Printf("error updating poi category: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update poi category: %v", err)
	}

	return h.mapper.ToProtoFromResponse(result), nil
}

// DeletePoiCategory handles delete request
func (h *PoiCategoryGrpcHandler) DeletePoiCategory(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi category id is required")
	}

	if err := h.usecase.Delete(ctx, req.Id); err != nil {
		log.Printf("error deleting poi category: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete poi category: %v", err)
	}

	return &sharepb.SubmitResponse{
		Message: "POI category deleted successfully",
	}, nil
}

// ListPoiCategories handles list request
func (h *PoiCategoryGrpcHandler) ListPoiCategories(ctx context.Context, req *tqdpb.ListPoiCategoriesRequest) (*tqdpb.ListPoiCategoriesResponse, error) {
	// Build filter
	filter := &dto.PoiCategoryFilter{
		Pagable: _dto.Pagable{
			Page: normalizePage(req.Page),
			Size: normalizeSize(req.Size),
		},
		Search:   req.Search,
		IsActive: &req.IsActive,
	}

	// Set optional filters
	if req.ParentId != 0 {
		filter.ParentID = &req.ParentId
	}
	if req.Level != 0 {
		level := int(req.Level)
		filter.Level = &level
	}

	// Call usecase
	result, err := h.usecase.List(ctx, filter)
	if err != nil {
		log.Printf("error listing poi categories: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to list poi categories: %v", err)
	}

	// Convert to proto
	categories := h.mapper.ToProtoListFromResponses(result.Data)

	return &tqdpb.ListPoiCategoriesResponse{
		Data:  categories,
		Total: result.Total,
	}, nil
}

// GetPoiCategoryTree handles get tree request
func (h *PoiCategoryGrpcHandler) GetPoiCategoryTree(ctx context.Context, req *emptypb.Empty) (*tqdpb.GetPoiCategoryTreeResponse, error) {
	results, err := h.usecase.GetTree(ctx)
	if err != nil {
		log.Printf("error getting poi category tree: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get poi category tree: %v", err)
	}

	// Build tree nodes
	treeNodes := h.buildProtoTree(results)

	return &tqdpb.GetPoiCategoryTreeResponse{
		Data: treeNodes,
	}, nil
}

// buildProtoTree builds proto tree from tree responses
func (h *PoiCategoryGrpcHandler) buildProtoTree(treeResponses []dto.PoiCategoryTreeResponse) []*tqdpb.PoiCategoryTreeNode {
	nodes := make([]*tqdpb.PoiCategoryTreeNode, len(treeResponses))
	for i, resp := range treeResponses {
		node := &tqdpb.PoiCategoryTreeNode{
			Data: &tqdpb.PoiCategory{
				Id:          resp.ID,
				Name:        resp.Name,
				Description: resp.Description,
				Code:        resp.Code,
				Icon:        resp.Icon,
				Color:       resp.Color,
				IsActive:    resp.IsActive,
				SortOrder:   uint32(resp.SortOrder),
				ParentId:    h.uint64PtrToValue(resp.ParentID),
				Level:       uint32(resp.Level),
				Path:        resp.Path,
				PoiCount:    uint32(resp.POICount),
			},
			Children: h.buildProtoTree(resp.Children),
		}
		nodes[i] = node
	}
	return nodes
}

func (h *PoiCategoryGrpcHandler) uint64PtrToValue(ptr *uint64) uint64 {
	if ptr == nil {
		return 0
	}
	return *ptr
}
