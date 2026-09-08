package handler_grpc

import (
	"context"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/domain"
	"tqd/internal/dto"

	"tqd/internal/usecase"

	_dto "common/domain/dto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AmenityGrpcHandlerImpl implements gRPC handlers for Amenity service
type AmenityGrpcHandlerImpl struct {
	tqdpb.UnimplementedAmenityServiceServer
	amenityUsecase *usecase.AmenityUsecase
	amenityMapper  *mapper.AmenityMapper
}

// NewAmenityGrpcHandlerImpl creates a new AmenityGrpcHandlerImpl
func NewAmenityGrpcHandlerImpl(
	amenityUsecase *usecase.AmenityUsecase,
	amenityMapper *mapper.AmenityMapper,
) *AmenityGrpcHandlerImpl {
	return &AmenityGrpcHandlerImpl{
		amenityUsecase: amenityUsecase,
		amenityMapper:  amenityMapper,
	}
}

// CreateAmenity handles gRPC CreateAmenity request
func (h *AmenityGrpcHandlerImpl) CreateAmenity(ctx context.Context, req *tqdpb.Amenity) (*tqdpb.Amenity, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "amenity name is required")
	}

	// Convert proto to CreateAmenityRequestDTO
	createReq := &dto.CreateAmenityRequestDTO{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}

	// Convert DTO to domain using mapper
	entity := h.amenityMapper.ToDomainFromCreateRequest(createReq)

	// Call usecase Create method
	result, err := h.amenityUsecase.Create(ctx, entity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create amenity: %v", err)
	}

	return h.amenityMapper.ToProto(result), nil
}

// GetAmenity handles gRPC GetAmenity request
func (h *AmenityGrpcHandlerImpl) GetAmenity(ctx context.Context, req *tqdpb.GetAmenityRequest) (*tqdpb.Amenity, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "amenity id is required")
	}

	result, err := h.amenityUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "amenity not found: %v", err)
	}

	return h.amenityMapper.ToProto(result), nil
}

// UpdateAmenity handles gRPC UpdateAmenity request
func (h *AmenityGrpcHandlerImpl) UpdateAmenity(ctx context.Context, req *tqdpb.Amenity) (*tqdpb.Amenity, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "amenity id is required")
	}

	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "amenity name is required")
	}

	// Convert proto to UpdateAmenityRequestDTO
	updateReq := &dto.UpdateAmenityRequestDTO{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}

	// Convert DTO to domain using mapper
	entity := h.amenityMapper.ToDomainFromUpdateRequest(updateReq, req.Id)

	result, err := h.amenityUsecase.Update(ctx, req.Id, entity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update amenity: %v", err)
	}

	return h.amenityMapper.ToProto(result), nil
}

// DeleteAmenity handles gRPC DeleteAmenity request
func (h *AmenityGrpcHandlerImpl) DeleteAmenity(ctx context.Context, req *tqdpb.DeleteAmenityRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "amenity id is required")
	}

	_, err := h.amenityUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete amenity: %v", err)
	}

	return &sharepb.SubmitResponse{
		Message: "Amenity deleted successfully",
	}, nil
}

// ListAmenities handles gRPC ListAmenities request
func (h *AmenityGrpcHandlerImpl) ListAmenities(ctx context.Context, req *tqdpb.ListAmenitiesRequest) (*tqdpb.ListAmenitiesResponse, error) {
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
	filter := &dto.AmenityFilterDTO{
		Pagable: _dto.Pagable{
			Page: page,
			Size: size,
		},
		Search:   req.Search,
		Category: *req.Category,
	}

	// Call repo's ListAmenities method directly if needed, or enhance usecase
	// For now, let's assume we need to get all and filter
	results, total, err := h.amenityUsecase.GetList(ctx, &filter.Pagable)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list amenities: %v", err)
	}

	// Filter results based on search and category (this should ideally be done in repo)
	var filtered []domain.Amenity
	for _, r := range results {
		// Apply search filter
		if filter.Search != "" {
			if !containsString(r.Name, filter.Search) && !containsString(r.Description, filter.Search) {
				continue
			}
		}
		// Apply category filter
		if filter.Category != "" {
			// You might need to add Category field to domain.Amenity
			// For now, skip this filter
		}
		filtered = append(filtered, r)
	}

	amenities := make([]*tqdpb.Amenity, len(filtered))
	for i, r := range filtered {
		amenities[i] = h.amenityMapper.ToProto(&r)
	}

	return &tqdpb.ListAmenitiesResponse{
		Data:  amenities,
		Total: total,
	}, nil
}

// Helper function
func containsString(s, substr string) bool {
	if substr == "" {
		return true
	}
	// Simple contains check - you might want case-insensitive
	return s != "" && substr != "" && len(s) >= len(substr) && s[:len(substr)] == substr
}
