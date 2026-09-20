package handler_grpc

import (
	_dto "common/domain/dto"
	"context"
	"fmt"
	"log/slog"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// DirectorySupplierGrpcHandler implements gRPC handlers for DirectorySupplier service
type DirectorySupplierGrpcHandler struct {
	tqdpb.UnimplementedDirectorySupplierServiceServer
	directorySupplierUsecase *usecase.DirectorySupplierUsecase
	directorySupplierMapper  *mapper.DirectorySupplierMapper
}

// NewDirectorySupplierGrpcHandler creates a new DirectorySupplierGrpcHandler
func NewDirectorySupplierGrpcHandler(
	directorySupplierUsecase *usecase.DirectorySupplierUsecase,
	directorySupplierMapper *mapper.DirectorySupplierMapper,
) *DirectorySupplierGrpcHandler {
	return &DirectorySupplierGrpcHandler{
		directorySupplierUsecase: directorySupplierUsecase,
		directorySupplierMapper:  directorySupplierMapper,
	}
}

// CreateDirectorySupplier handles gRPC CreateDirectorySupplier request
func (h *DirectorySupplierGrpcHandler) CreateDirectorySupplier(ctx context.Context, req *tqdpb.DirectorySupplier) (*tqdpb.DirectorySupplier, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "directory supplier name is required")
	}
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "directory supplier code is required")
	}

	// Convert proto to CreateDirectorySupplierRequestDTO
	createReq := &dto.CreateDirectorySupplierRequestDTO{
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		ContactPerson:   req.ContactPerson,
		Phone:           req.Phone,
		Email:           req.Email,
		Address:         req.Address,
		Website:         req.Website,
		TaxCode:         req.TaxCode,
		BusinessLicense: req.BusinessLicense,
		Rating:          req.Rating,
		IsActive:        req.IsActive,
		Notes:           req.Notes,
		Tags:            req.Tags,
		Categories:      req.Categories,
	}

	// Set dates if provided
	if req.JoinDate != "" {
		createReq.JoinDate = &req.JoinDate
	}
	if req.LastContactDate != "" {
		createReq.LastContactDate = &req.LastContactDate
	}

	// Call usecase
	result, err := h.directorySupplierUsecase.Create(ctx, createReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error creating directory supplier: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to create directory supplier: %v", err)
	}

	return h.directorySupplierMapper.ToProto(result), nil
}

// GetDirectorySupplier handles gRPC GetDirectorySupplier request
func (h *DirectorySupplierGrpcHandler) GetDirectorySupplier(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.DirectorySupplier, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory supplier id is required")
	}

	result, err := h.directorySupplierUsecase.GetByID(ctx, req.Id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error getting directory supplier: %v", err))
		return nil, status.Errorf(codes.NotFound, "directory supplier not found: %v", err)
	}

	return h.directorySupplierMapper.ToProto(result), nil
}

// UpdateDirectorySupplier handles gRPC UpdateDirectorySupplier request
func (h *DirectorySupplierGrpcHandler) UpdateDirectorySupplier(ctx context.Context, req *tqdpb.DirectorySupplier) (*tqdpb.DirectorySupplier, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory supplier id is required")
	}

	// Validate required fields
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "directory supplier name is required")
	}
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "directory supplier code is required")
	}

	// Convert proto to UpdateDirectorySupplierRequestDTO
	updateReq := &dto.UpdateDirectorySupplierRequestDTO{
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		ContactPerson:   req.ContactPerson,
		Phone:           req.Phone,
		Email:           req.Email,
		Address:         req.Address,
		Website:         req.Website,
		TaxCode:         req.TaxCode,
		BusinessLicense: req.BusinessLicense,
		Rating:          req.Rating,
		IsActive:        req.IsActive,
		Notes:           req.Notes,
		Tags:            req.Tags,
		Categories:      req.Categories,
	}

	// Set dates if provided
	if req.JoinDate != "" {
		updateReq.JoinDate = &req.JoinDate
	}
	if req.LastContactDate != "" {
		updateReq.LastContactDate = &req.LastContactDate
	}

	// Call usecase
	result, err := h.directorySupplierUsecase.Update(ctx, req.Id, updateReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error updating directory supplier: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to update directory supplier: %v", err)
	}

	return h.directorySupplierMapper.ToProto(result), nil
}

// DeleteDirectorySupplier handles gRPC DeleteDirectorySupplier request
func (h *DirectorySupplierGrpcHandler) DeleteDirectorySupplier(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory supplier id is required")
	}

	err := h.directorySupplierUsecase.Delete(ctx, req.Id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error deleting directory supplier: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to delete directory supplier: %v", err)
	}

	return &sharepb.SubmitResponse{
		Message: "Directory supplier deleted successfully",
	}, nil
}

// ListDirectorySuppliers handles gRPC ListDirectorySuppliers request
func (h *DirectorySupplierGrpcHandler) ListDirectorySuppliers(ctx context.Context, req *tqdpb.ListDirectorySuppliersRequest) (*tqdpb.ListDirectorySuppliersResponse, error) {
	// Set default pagination
	page := uint32(req.Page)
	if page < 1 {
		page = 1
	}
	size := uint32(req.Size)
	if size < 1 {
		size = 10
	}
	if size > 100 {
		size = 100
	}

	// Create filter DTO
	filters := &dto.DirectorySupplierFilterDTO{
		Pagable: _dto.Pagable{
			Page: page,
			Size: size,
		},
		Search: req.Search,
	}

	// Set optional filters
	if len(req.Categories) > 0 {
		filters.Categories = req.Categories
	}

	// Handle boolean filters
	if req.IsActive {
		filters.IsActive = &req.IsActive
	}

	// Handle rating filters
	if req.MinRating > 0 {
		filters.MinRating = &req.MinRating
	}
	if req.MaxRating > 0 {
		filters.MaxRating = &req.MaxRating
	}

	// Call usecase
	result, err := h.directorySupplierUsecase.List(ctx, filters)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error listing directory suppliers: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to list directory suppliers: %v", err)
	}

	// Convert DTOs to proto
	suppliers := make([]*tqdpb.DirectorySupplier, len(result.Data))
	for i, item := range result.Data {
		suppliers[i] = h.directorySupplierMapper.ToProto(&item)
	}

	return &tqdpb.ListDirectorySuppliersResponse{
		Data:  suppliers,
		Total: result.Total,
	}, nil
}
