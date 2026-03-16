package handler_grpc

import (
	_dto "common/domain/dto"
	"context"
	"log"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/infra/validator"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// DirectoryCategoryGrpcHandler implements gRPC handlers for DirectoryCategory service
type DirectoryCategoryGrpcHandler struct {
	tqdpb.UnimplementedDirectoryCategoryServiceServer
	directoryCategoryUsecase *usecase.DirectoryCategoryUsecase
	directoryCategoryMapper  *mapper.DirectoryCategoryMapper
	validator                *validator.DirectoryCategoryValidator
}

// NewDirectoryCategoryGrpcHandler creates a new DirectoryCategoryGrpcHandler
func NewDirectoryCategoryGrpcHandler(
	directoryCategoryUsecase *usecase.DirectoryCategoryUsecase,
	directoryCategoryMapper *mapper.DirectoryCategoryMapper,
	validator *validator.DirectoryCategoryValidator,
) *DirectoryCategoryGrpcHandler {
	return &DirectoryCategoryGrpcHandler{
		directoryCategoryUsecase: directoryCategoryUsecase,
		directoryCategoryMapper:  directoryCategoryMapper,
		validator:                validator,
	}
}

// CreateDirectoryCategory handles gRPC CreateDirectoryCategory request
func (h *DirectoryCategoryGrpcHandler) CreateDirectoryCategory(ctx context.Context, req *tqdpb.DirectoryCategory) (*tqdpb.DirectoryCategory, error) {
	category := &domain.DirectoryCategory{
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		Icon:        req.Icon,
		Color:       req.Color,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		Level:       req.Level,
	}

	// Validate request
	if err := h.validator.ValidateCreateRequest(category); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.directoryCategoryUsecase.CreateDirectoryCategory(ctx, category); err != nil {
		log.Printf("error creating directory category: %v", err)
		return nil, status.Error(codes.Internal, "failed to create directory category")
	}

	return h.directoryCategoryMapper.ToProto(category), nil
}

// GetDirectoryCategory handles gRPC GetDirectoryCategory request
func (h *DirectoryCategoryGrpcHandler) GetDirectoryCategory(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.DirectoryCategory, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory category id is required")
	}

	result, err := h.directoryCategoryUsecase.GetDirectoryCategoryByID(ctx, req.Id)
	if err != nil {
		log.Printf("error getting directory category: %v", err)
		return nil, status.Error(codes.NotFound, "directory category not found")
	}

	return h.directoryCategoryMapper.ToProto(result), nil
}

// UpdateDirectoryCategory handles gRPC UpdateDirectoryCategory request
func (h *DirectoryCategoryGrpcHandler) UpdateDirectoryCategory(ctx context.Context, req *tqdpb.DirectoryCategory) (*tqdpb.DirectoryCategory, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory category id is required")
	}

	category := &domain.DirectoryCategory{
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		Icon:        req.Icon,
		Color:       req.Color,
		IsActive:    req.IsActive,
		SortOrder:   req.SortOrder,
		Level:       req.Level,
	}

	// Validate request
	if err := h.validator.ValidateUpdateRequest(category); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := h.directoryCategoryUsecase.UpdateDirectoryCategory(ctx, req.Id, category)
	if err != nil {
		log.Printf("error updating directory category: %v", err)
		return nil, status.Error(codes.Internal, "failed to update directory category")
	}

	return h.directoryCategoryMapper.ToProto(result), nil
}

// DeleteDirectoryCategory handles gRPC DeleteDirectoryCategory request
func (h *DirectoryCategoryGrpcHandler) DeleteDirectoryCategory(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory category id is required")
	}

	err := h.directoryCategoryUsecase.DeleteDirectoryCategory(ctx, req.Id)
	if err != nil {
		log.Printf("error deleting directory category: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete directory category")
	}

	return &sharepb.SubmitResponse{
		Message: "Directory category deleted successfully",
	}, nil
}

// ListServiceCategories handles gRPC ListServiceCategories request
func (h *DirectoryCategoryGrpcHandler) ListServiceCategories(ctx context.Context, req *tqdpb.ListServiceCategoriesRequest) (*tqdpb.ListServiceCategoriesResponse, error) {
	page := req.Page
	size := req.Size
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}

	listReq := &dto.ListDirectoryCategoriesRequest{
		Search:   req.Search,
		IsActive: nil, // for all
	}
	listReq.Pagable.Page = page
	listReq.Pagable.Size = size

	results, total, err := h.directoryCategoryUsecase.ListDirectoryCategories(ctx, listReq)
	if err != nil {
		log.Printf("error listing directory categories: %v", err)
		return nil, status.Error(codes.Internal, "failed to list directory categories")
	}

	categories := make([]*tqdpb.DirectoryCategory, len(results))
	for i, result := range results {
		categories[i] = h.directoryCategoryMapper.ToProto(&result)
	}

	return &tqdpb.ListServiceCategoriesResponse{
		Data:  categories,
		Total: total,
	}, nil
}

// GetActiveServiceCategories handles gRPC GetActiveServiceCategories request
func (h *DirectoryCategoryGrpcHandler) GetActiveServiceCategories(ctx context.Context, req *emptypb.Empty) (*tqdpb.ListServiceCategoriesResponse, error) {
	results, total, err := h.directoryCategoryUsecase.ListDirectoryCategories(ctx, &dto.ListDirectoryCategoriesRequest{
		Pagable:  _dto.Pagable{Page: 1, Size: 1000},
		IsActive: func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		log.Printf("error listing active directory categories: %v", err)
		return nil, status.Error(codes.Internal, "failed to list active directory categories")
	}

	categories := make([]*tqdpb.DirectoryCategory, len(results))
	for i, result := range results {
		categories[i] = h.directoryCategoryMapper.ToProto(&result)
	}

	return &tqdpb.ListServiceCategoriesResponse{
		Data:  categories,
		Total: total,
	}, nil
}

// GetServiceCategoriesByLevel handles gRPC GetServiceCategoriesByLevel request
func (h *DirectoryCategoryGrpcHandler) GetServiceCategoriesByLevel(ctx context.Context, req *tqdpb.GetServiceCategoriesByLevelRequest) (*tqdpb.ListServiceCategoriesResponse, error) {
	if req.Level <= 0 {
		return nil, status.Error(codes.InvalidArgument, "level must be greater than 0")
	}

	results, _, err := h.directoryCategoryUsecase.ListDirectoryCategories(ctx, &dto.ListDirectoryCategoriesRequest{
		Pagable: _dto.Pagable{Page: 1, Size: 1000},
	})
	if err != nil {
		log.Printf("error listing directory categories by level: %v", err)
		return nil, status.Error(codes.Internal, "failed to list directory categories by level")
	}

	// Filter by level manually
	filtered := make([]*tqdpb.DirectoryCategory, 0, len(results))
	for _, result := range results {
		if result.Level == req.Level {
			filtered = append(filtered, h.directoryCategoryMapper.ToProto(&result))
		}
	}

	return &tqdpb.ListServiceCategoriesResponse{
		Data:  filtered,
		Total: int64(len(filtered)),
	}, nil
}
