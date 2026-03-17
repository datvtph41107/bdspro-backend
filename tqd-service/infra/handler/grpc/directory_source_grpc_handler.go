package handler_grpc

import (
	"context"
	"log"

	_dto "common/domain/dto"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/infra/validator"
	"tqd/internal/domain"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DirectorySourceGrpcHandler implements gRPC handlers for DirectorySource service
type DirectorySourceGrpcHandler struct {
	tqdpb.UnimplementedDirectorySourceServiceServer
	directorySourceUsecase *usecase.DirectorySourceUsecase
	directorySourceMapper  *mapper.DirectorySourceMapper
	validator              *validator.DirectorySourceValidator
}

// NewDirectorySourceGrpcHandler creates a new DirectorySourceGrpcHandler
func NewDirectorySourceGrpcHandler(
	directorySourceUsecase *usecase.DirectorySourceUsecase,
	directorySourceMapper *mapper.DirectorySourceMapper,
	validator *validator.DirectorySourceValidator,
) *DirectorySourceGrpcHandler {
	return &DirectorySourceGrpcHandler{
		directorySourceUsecase: directorySourceUsecase,
		directorySourceMapper:  directorySourceMapper,
		validator:              validator,
	}
}

// CreateDirectorySource handles gRPC CreateDirectorySource request
func (h *DirectorySourceGrpcHandler) CreateDirectorySource(ctx context.Context, req *tqdpb.DirectorySource) (*tqdpb.DirectorySource, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "directory source name is required")
	}

	source := &dto.DirectorySourceDTO{
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		Type:           req.Type,
		Icon:           req.Icon,
		Color:          req.Color,
		IsActive:       req.IsActive,
		SortOrder:      req.SortOrder,
		ExpectedAmount: req.ExpectedAmount,
		ActualAmount:   req.ActualAmount,
		Frequency:      req.Frequency,
		IsRecurring:    req.IsRecurring,
		PaymentMethod:  req.PaymentMethod,
		Notes:          req.Notes,
	}

	// Validate request
	if err := h.validator.ValidateCreateRequest(source); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := h.directorySourceUsecase.Create(ctx, source)
	if err != nil {
		log.Printf("error creating directory source: %v", err)
		return nil, status.Error(codes.Internal, "failed to create directory source")
	}

	return h.directorySourceMapper.ToProto(result), nil
}

// GetDirectorySource handles gRPC GetDirectorySource request
func (h *DirectorySourceGrpcHandler) GetDirectorySource(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.DirectorySource, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory source id is required")
	}

	result, err := h.directorySourceUsecase.GetByID(ctx, req.Id)
	if err != nil {
		log.Printf("error getting directory source: %v", err)
		return nil, status.Error(codes.NotFound, "directory source not found")
	}

	return h.directorySourceMapper.ToProto(result), nil
}

// UpdateDirectorySource handles gRPC UpdateDirectorySource request
func (h *DirectorySourceGrpcHandler) UpdateDirectorySource(ctx context.Context, req *tqdpb.DirectorySource) (*tqdpb.DirectorySource, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory source id is required")
	}

	source := &domain.DirectorySource{
		Name:           req.Name,
		Code:           req.Code,
		Description:    req.Description,
		Type:           req.Type,
		Icon:           req.Icon,
		Color:          req.Color,
		IsActive:       req.IsActive,
		SortOrder:      req.SortOrder,
		ExpectedAmount: req.ExpectedAmount,
		ActualAmount:   req.ActualAmount,
		Frequency:      req.Frequency,
		IsRecurring:    req.IsRecurring,
		PaymentMethod:  req.PaymentMethod,
		Notes:          req.Notes,
	}

	// Validate request
	if err := h.validator.ValidateUpdateRequest(source); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := h.directorySourceUsecase.Update(ctx, source)
	if err != nil {
		log.Printf("error updating directory source: %v", err)
		return nil, status.Error(codes.Internal, "failed to update directory source")
	}

	return h.directorySourceMapper.ToProto(result), nil
}

// DeleteDirectorySource handles gRPC DeleteDirectorySource request
func (h *DirectorySourceGrpcHandler) DeleteDirectorySource(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory source id is required")
	}

	err := h.directorySourceUsecase.Delete(ctx, req.Id)
	if err != nil {
		log.Printf("error deleting directory source: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete directory source")
	}

	return &sharepb.SubmitResponse{
		Message: "Directory source deleted successfully",
	}, nil
}

// ListDirectorySources handles gRPC ListDirectorySources request
func (h *DirectorySourceGrpcHandler) ListDirectorySources(ctx context.Context, req *tqdpb.ListDirectorySourcesRequest) (*tqdpb.ListDirectorySourcesResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	filter := &dto.ListDirectorySourcesRequestDTO{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Search:      req.Search,
		Category:    func(s string) *string { return &s }(req.Category),
		Type:        func(s string) *string { return &s }(req.Type),
		IsActive:    func(b bool) *bool { return &b }(req.IsActive),
		IsRecurring: func(b bool) *bool { return &b }(req.IsRecurring),
	}

	results, total, err := h.directorySourceUsecase.List(ctx, filter)
	if err != nil {
		log.Printf("error listing directory sources: %v", err)
		return nil, status.Error(codes.Internal, "failed to list directory sources")
	}

	sources := make([]*tqdpb.DirectorySource, len(results))
	for i, result := range results {
		sources[i] = h.directorySourceMapper.ToProto(&result)
	}

	return &tqdpb.ListDirectorySourcesResponse{
		Data:  sources,
		Total: total,
	}, nil
}
