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
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "directory supplier name is required")
	}

	createReq := &dto.CreateDirectorySupplierRequestDTO{
		Name:             req.Name,
		Code:             req.Code,
		Description:      req.Description,
		ContactPerson:    req.ContactPerson,
		Phone:            req.Phone,
		Email:            req.Email,
		Address:          req.Address,
		Website:          req.Website,
		TaxCode:          req.TaxCode,
		BusinessLicense:  req.BusinessLicense,
		Rating:           req.Rating,
		IsActive:         req.IsActive,
		Notes:            req.Notes,
		Tags:             req.Tags,
		Categories:       req.Categories,
	}

	result, err := h.directorySupplierUsecase.CreateDirectorySupplier(ctx, createReq)
	if err != nil {
		log.Printf("error creating directory supplier: %v", err)
		return nil, status.Error(codes.Internal, "failed to create directory supplier")
	}

	return h.directorySupplierMapper.ToProto(result), nil
}

// GetDirectorySupplier handles gRPC GetDirectorySupplier request
func (h *DirectorySupplierGrpcHandler) GetDirectorySupplier(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.DirectorySupplier, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory supplier id is required")
	}

	result, err := h.directorySupplierUsecase.GetDirectorySupplierByID(ctx, req.Id)
	if err != nil {
		log.Printf("error getting directory supplier: %v", err)
		return nil, status.Error(codes.NotFound, "directory supplier not found")
	}

	return h.directorySupplierMapper.ToProto(result), nil
}

// UpdateDirectorySupplier handles gRPC UpdateDirectorySupplier request
func (h *DirectorySupplierGrpcHandler) UpdateDirectorySupplier(ctx context.Context, req *tqdpb.DirectorySupplier) (*tqdpb.DirectorySupplier, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory supplier id is required")
	}

	updateReq := &dto.DirectorySupplierDTO{
		ID:               req.Id,
		Name:             req.Name,
		Code:             req.Code,
		Description:      req.Description,
		ContactPerson:    req.ContactPerson,
		Phone:            req.Phone,
		Email:            req.Email,
		Address:          req.Address,
		Website:          req.Website,
		TaxCode:          req.TaxCode,
		BusinessLicense:  req.BusinessLicense,
		Rating:           req.Rating,
		IsActive:         req.IsActive,
		Notes:            req.Notes,
		Tags:             req.Tags,
		Categories:       req.Categories,
	}

	result, err := h.directorySupplierUsecase.UpdateDirectorySupplier(ctx, updateReq)
	if err != nil {
		log.Printf("error updating directory supplier: %v", err)
		return nil, status.Error(codes.Internal, "failed to update directory supplier")
	}

	return h.directorySupplierMapper.ToProto(result), nil
}

// DeleteDirectorySupplier handles gRPC DeleteDirectorySupplier request
func (h *DirectorySupplierGrpcHandler) DeleteDirectorySupplier(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "directory supplier id is required")
	}

	err := h.directorySupplierUsecase.DeleteDirectorySupplier(ctx, req.Id)
	if err != nil {
		log.Printf("error deleting directory supplier: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete directory supplier")
	}

	return &sharepb.SubmitResponse{
		Code:    200,
		Message: "Directory supplier deleted successfully",
		Data:    "",
	}, nil
}

// ListDirectorySuppliers handles gRPC ListDirectorySuppliers request
func (h *DirectorySupplierGrpcHandler) ListDirectorySuppliers(ctx context.Context, req *tqdpb.ListDirectorySuppliersRequest) (*tqdpb.ListDirectorySuppliersResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	filters := &dto.DirectorySupplierFilterDTO{
		Page:       int(req.Page),
		Size:       int(req.Size),
		Search:     req.Search,
		Categories: req.Categories,
		IsActive:   req.IsActive,
		MinRating:  req.MinRating,
		MaxRating:  req.MaxRating,
	}

	results, total, err := h.directorySupplierUsecase.ListDirectorySuppliers(ctx, filters)
	if err != nil {
		log.Printf("error listing directory suppliers: %v", err)
		return nil, status.Error(codes.Internal, "failed to list directory suppliers")
	}

	suppliers := make([]*tqdpb.DirectorySupplier, len(results))
	for i, result := range results {
		suppliers[i] = h.directorySupplierMapper.ToProto(result)
	}

	return &tqdpb.ListDirectorySuppliersResponse{
		Data:  suppliers,
		Total: int64(total),
	}, nil
}
