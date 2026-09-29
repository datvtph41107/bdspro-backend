package handler_grpc

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	_dto "common/domain/dto"
	_utils "common/utils"
	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/dto"
	"tqd/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PoiGrpcHandler handles gRPC requests for POI
type PoiGrpcHandler struct {
	tqdpb.UnimplementedPoiServiceServer
	usecase      usecase.PoiUsecase
	mapper       *mapper.PoiMapper
	SyncProvider *_utils.SyncUtil
}

func NewPoiGrpcHandler(
	usecase usecase.PoiUsecase,
	mapper *mapper.PoiMapper,
	syncProvider *_utils.SyncUtil,
) *PoiGrpcHandler {
	return &PoiGrpcHandler{
		usecase:      usecase,
		mapper:       mapper,
		SyncProvider: syncProvider,
	}
}

// CreatePoi handles create request
func (h *PoiGrpcHandler) CreatePoi(ctx context.Context, req *tqdpb.Poi) (*tqdpb.Poi, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "poi name is required")
	}
	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address is required")
	}
	if req.CategoryId == 0 {
		return nil, status.Error(codes.InvalidArgument, "category ID is required")
	}

	// Convert proto to create request
	createReq := &dto.CreatePoiRequest{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Address:     req.Address,
		Phone:       req.Phone,
		Email:       req.Email,
		Website:     req.Website,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		CategoryID:  req.CategoryId,
		IsActive:    req.IsActive,
		IsFeatured:  req.IsFeatured,
		CoverImage:  req.CoverImage,
		Images:      req.Images,
		Tags:        req.Tags,
	}

	// Call usecase
	result, err := h.usecase.Create(ctx, createReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error creating poi: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to create poi: %v", err)
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDPOIList, 0), t.UnixMilli())

	return h.mapper.ToProtoFromResponse(result), nil
}

// GetPoi handles get by ID request
func (h *PoiGrpcHandler) GetPoi(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.Poi, error) {
	if req == nil || req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	key := h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDPOIDetail, req.Id)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.Poi{}, nil
	}

	result, err := h.usecase.GetByID(ctx, req.Id)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error getting poi: %v", err))
		return nil, status.Errorf(codes.NotFound, "poi not found: %v", err)
	}

	return h.mapper.ToProtoFromResponse(result), nil
}

// GetPoiByCode handles get by code request
func (h *PoiGrpcHandler) GetPoiByCode(ctx context.Context, req *tqdpb.GetPoiByCodeRequest) (*tqdpb.Poi, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "poi code is required")
	}

	result, err := h.usecase.GetByCode(ctx, req.Code)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error getting poi by code: %v", err))
		return nil, status.Errorf(codes.NotFound, "poi not found: %v", err)
	}

	return h.mapper.ToProtoFromResponse(result), nil
}

// GetNearbyPois handles nearby search request
func (h *PoiGrpcHandler) GetNearbyPois(ctx context.Context, req *tqdpb.GetNearbyPoisRequest) (*tqdpb.ListPoisResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDPOIList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListPoisResponse{}, nil
	}

	if req.Latitude < -90 || req.Latitude > 90 || req.Longitude < -180 || req.Longitude > 180 {
		return nil, status.Error(codes.InvalidArgument, "latitude and longitude are required")
	}

	// Create nearby request
	nearbyReq := &dto.NearbyPoiRequest{
		Pagable: _dto.Pagable{
			Page: normalizePage(req.Page),
			Size: normalizeSize(req.Size),
		},
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Radius:      req.Radius,
		CategoryIDs: req.CategoryIds,
	}

	// Call usecase
	result, err := h.usecase.ListNearby(ctx, nearbyReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error getting nearby pois: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get nearby pois: %v", err)
	}

	// Convert to proto
	pois := h.mapper.ToProtoListFromResponses(result.Data)

	return &tqdpb.ListPoisResponse{
		Data:  pois,
		Total: result.Total,
	}, nil
}

// GetPoisByCategory handles get by category request
func (h *PoiGrpcHandler) GetPoisByCategory(ctx context.Context, req *tqdpb.GetPoisByCategoryRequest) (*tqdpb.ListPoisResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDPOIList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListPoisResponse{}, nil
	}

	if req.CategoryId == 0 {
		return nil, status.Error(codes.InvalidArgument, "category id is required")
	}

	limit := int(req.Size)
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	results, err := h.usecase.ListByCategory(ctx, req.CategoryId, limit)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error getting pois by category: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to get pois by category: %v", err)
	}

	pois := h.mapper.ToProtoListFromResponses(results)

	return &tqdpb.ListPoisResponse{
		Data:  pois,
		Total: int64(len(pois)),
	}, nil
}

// UpdatePoi handles update request
func (h *PoiGrpcHandler) UpdatePoi(ctx context.Context, req *tqdpb.Poi) (*tqdpb.Poi, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi id is required")
	}

	// Build update request
	updateReq := &dto.UpdatePoiRequest{}

	if req.Name != "" {
		updateReq.Name = &req.Name
	}
	if req.Code != "" {
		updateReq.Code = &req.Code
	}
	if req.Description != "" {
		updateReq.Description = &req.Description
	}
	if req.Address != "" {
		updateReq.Address = &req.Address
	}
	if req.Phone != "" {
		updateReq.Phone = &req.Phone
	}
	if req.Email != "" {
		updateReq.Email = &req.Email
	}
	if req.Website != "" {
		updateReq.Website = &req.Website
	}
	if req.Latitude != 0 {
		updateReq.Latitude = &req.Latitude
	}
	if req.Longitude != 0 {
		updateReq.Longitude = &req.Longitude
	}
	if req.CategoryId != 0 {
		updateReq.CategoryID = &req.CategoryId
	}
	updateReq.IsActive = &req.IsActive
	updateReq.IsFeatured = &req.IsFeatured

	if req.CoverImage != "" {
		updateReq.CoverImage = &req.CoverImage
	}
	if len(req.Images) > 0 {
		updateReq.Images = req.Images
	}
	if len(req.Tags) > 0 {
		updateReq.Tags = req.Tags
	}

	// Call usecase
	result, err := h.usecase.Update(ctx, req.Id, updateReq)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error updating poi: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to update poi: %v", err)
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDPOIList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDPOIDetail, req.Id), t.UnixMilli())

	return h.mapper.ToProtoFromResponse(result), nil
}

// UpdatePoiRating handles update rating request
func (h *PoiGrpcHandler) UpdatePoiRating(ctx context.Context, req *tqdpb.UpdatePoiRatingRequest) (*tqdpb.Poi, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi id is required")
	}

	result, err := h.usecase.UpdateRating(ctx, req.Id, req.Rating, req.ReviewCount)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error updating poi rating: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to update poi rating: %v", err)
	}

	return h.mapper.ToProtoFromResponse(result), nil
}

// DeletePoi handles delete request
func (h *PoiGrpcHandler) DeletePoi(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi id is required")
	}

	if err := h.usecase.Delete(ctx, req.Id); err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error deleting poi: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to delete poi: %v", err)
	}

	t := time.Now()
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDPOIList, 0), t.UnixMilli())
	h.SyncProvider.PutTimeRequest(ctx, h.SyncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyTQDPOIDetail, req.Id), t.UnixMilli())

	return &sharepb.SubmitResponse{
		Message: "POI deleted successfully",
	}, nil
}

// ListPois handles list request
func (h *PoiGrpcHandler) ListPois(ctx context.Context, req *tqdpb.ListPoisRequest) (*tqdpb.ListPoisResponse, error) {
	key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyTQDPOIList, 0)
	updated := h.SyncProvider.HasUpdated(ctx, key, req.Timestamp)
	h.SyncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &tqdpb.ListPoisResponse{}, nil
	}

	filter := &dto.PoiFilter{
		Pagable: _dto.Pagable{
			Page: normalizePage(req.Page),
			Size: normalizeSize(req.Size),
		},
		Search:      req.Search,
		CategoryIDs: req.CategoryIds,
		SortBy:      req.SortBy,
	}

	// Set optional filters
	if req.CategoryId != 0 {
		filter.CategoryID = &req.CategoryId
	}
	filter.IsActive = &req.IsActive
	filter.IsVerified = &req.IsVerified

	if req.MinRating > 0 {
		filter.MinRating = &req.MinRating
	}
	if req.MaxRating > 0 {
		filter.MaxRating = &req.MaxRating
	}

	// Call usecase
	result, err := h.usecase.List(ctx, filter)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("error listing pois: %v", err))
		return nil, status.Errorf(codes.Internal, "failed to list pois: %v", err)
	}

	// Convert to proto
	pois := h.mapper.ToProtoListFromResponses(result.Data)

	return &tqdpb.ListPoisResponse{
		Data:  pois,
		Total: result.Total,
	}, nil
}
