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

// PoiGrpcHandler implements gRPC handlers for Poi service
type PoiGrpcHandler struct {
	tqdpb.UnimplementedPoiServiceServer
	poiUsecase *usecase.PoiUsecase
	poiMapper  *mapper.PoiMapper
}

// NewPoiGrpcHandler creates a new PoiGrpcHandler
func NewPoiGrpcHandler(
	poiUsecase *usecase.PoiUsecase,
	poiMapper *mapper.PoiMapper,
) *PoiGrpcHandler {
	return &PoiGrpcHandler{
		poiUsecase: poiUsecase,
		poiMapper:  poiMapper,
	}
}

// CreatePoi handles gRPC CreatePoi request
func (h *PoiGrpcHandler) CreatePoi(ctx context.Context, req *tqdpb.Poi) (*tqdpb.Poi, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "poi name is required")
	}

	createReq := &dto.CreatePoiRequestDTO{
		Name:       req.Name,
		Code:       req.Code,
		Description: req.Description,
		Address:    req.Address,
		Phone:      req.Phone,
		Email:      req.Email,
		Website:    req.Website,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		CategoryId: req.CategoryId,
		IsActive:   req.IsActive,
		IsFeatured: req.IsFeatured,
		CoverImage: req.CoverImage,
		Images:     req.Images,
		Tags:       req.Tags,
	}

	result, err := h.poiUsecase.CreatePoi(ctx, createReq)
	if err != nil {
		log.Printf("error creating poi: %v", err)
		return nil, status.Error(codes.Internal, "failed to create poi")
	}

	return h.poiMapper.ToProto(result), nil
}

// GetPoi handles gRPC GetPoi request
func (h *PoiGrpcHandler) GetPoi(ctx context.Context, req *sharepb.IdRequest) (*tqdpb.Poi, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi id is required")
	}

	result, err := h.poiUsecase.GetPoiByID(ctx, req.Id)
	if err != nil {
		log.Printf("error getting poi: %v", err)
		return nil, status.Error(codes.NotFound, "poi not found")
	}

	return h.poiMapper.ToProto(result), nil
}

// GetPoiByCode handles gRPC GetPoiByCode request
func (h *PoiGrpcHandler) GetPoiByCode(ctx context.Context, req *tqdpb.GetPoiByCodeRequest) (*tqdpb.Poi, error) {
	if req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "poi code is required")
	}

	result, err := h.poiUsecase.GetPoiByCode(ctx, req.Code)
	if err != nil {
		log.Printf("error getting poi by code: %v", err)
		return nil, status.Error(codes.NotFound, "poi not found")
	}

	return h.poiMapper.ToProto(result), nil
}

// GetNearbyPois handles gRPC GetNearbyPois request
func (h *PoiGrpcHandler) GetNearbyPois(ctx context.Context, req *tqdpb.GetNearbyPoisRequest) (*tqdpb.ListPoisResponse, error) {
	if req.Latitude == 0 || req.Longitude == 0 {
		return nil, status.Error(codes.InvalidArgument, "latitude and longitude are required")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	filters := &dto.GetNearbyPoisFilterDTO{
		Page:        int(req.Page),
		Size:        int(req.Size),
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Radius:      req.Radius,
		CategoryIds: req.CategoryIds,
	}

	results, total, err := h.poiUsecase.GetNearbyPois(ctx, filters)
	if err != nil {
		log.Printf("error getting nearby pois: %v", err)
		return nil, status.Error(codes.Internal, "failed to get nearby pois")
	}

	pois := make([]*tqdpb.Poi, len(results))
	for i, result := range results {
		pois[i] = h.poiMapper.ToProto(result)
	}

	return &tqdpb.ListPoisResponse{
		Data:  pois,
		Total: int64(total),
	}, nil
}

// GetPoisByCategory handles gRPC GetPoisByCategory request
func (h *PoiGrpcHandler) GetPoisByCategory(ctx context.Context, req *tqdpb.GetPoisByCategoryRequest) (*tqdpb.ListPoisResponse, error) {
	if req.CategoryId == 0 {
		return nil, status.Error(codes.InvalidArgument, "category id is required")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	filters := &dto.PoiFilterDTO{
		Page:       int(req.Page),
		Size:       int(req.Size),
		Search:     req.Search,
		CategoryId: req.CategoryId,
	}

	results, total, err := h.poiUsecase.ListPois(ctx, filters)
	if err != nil {
		log.Printf("error getting pois by category: %v", err)
		return nil, status.Error(codes.Internal, "failed to get pois by category")
	}

	pois := make([]*tqdpb.Poi, len(results))
	for i, result := range results {
		pois[i] = h.poiMapper.ToProto(result)
	}

	return &tqdpb.ListPoisResponse{
		Data:  pois,
		Total: int64(total),
	}, nil
}

// UpdatePoi handles gRPC UpdatePoi request
func (h *PoiGrpcHandler) UpdatePoi(ctx context.Context, req *tqdpb.Poi) (*tqdpb.Poi, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi id is required")
	}

	updateReq := &dto.PoiDTO{
		ID:         req.Id,
		Name:       req.Name,
		Code:       req.Code,
		Description: req.Description,
		Address:    req.Address,
		Phone:      req.Phone,
		Email:      req.Email,
		Website:    req.Website,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		CategoryId: req.CategoryId,
		IsActive:   req.IsActive,
		IsFeatured: req.IsFeatured,
		CoverImage: req.CoverImage,
		Images:     req.Images,
		Tags:       req.Tags,
	}

	result, err := h.poiUsecase.UpdatePoi(ctx, updateReq)
	if err != nil {
		log.Printf("error updating poi: %v", err)
		return nil, status.Error(codes.Internal, "failed to update poi")
	}

	return h.poiMapper.ToProto(result), nil
}

// UpdatePoiRating handles gRPC UpdatePoiRating request
func (h *PoiGrpcHandler) UpdatePoiRating(ctx context.Context, req *tqdpb.UpdatePoiRatingRequest) (*tqdpb.Poi, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi id is required")
	}

	result, err := h.poiUsecase.UpdatePoiRating(ctx, req.Id, req.Rating, uint32(req.ReviewCount))
	if err != nil {
		log.Printf("error updating poi rating: %v", err)
		return nil, status.Error(codes.Internal, "failed to update poi rating")
	}

	return h.poiMapper.ToProto(result), nil
}

// DeletePoi handles gRPC DeletePoi request
func (h *PoiGrpcHandler) DeletePoi(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "poi id is required")
	}

	err := h.poiUsecase.DeletePoi(ctx, req.Id)
	if err != nil {
		log.Printf("error deleting poi: %v", err)
		return nil, status.Error(codes.Internal, "failed to delete poi")
	}

	return &sharepb.SubmitResponse{
		Code:    200,
		Message: "POI deleted successfully",
		Data:    "",
	}, nil
}

// ListPois handles gRPC ListPois request
func (h *PoiGrpcHandler) ListPois(ctx context.Context, req *tqdpb.ListPoisRequest) (*tqdpb.ListPoisResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	filters := &dto.PoiFilterDTO{
		Page:        int(req.Page),
		Size:        int(req.Size),
		Search:      req.Search,
		CategoryId:  req.CategoryId,
		CategoryIds: req.CategoryIds,
		IsActive:    req.IsActive,
		IsVerified:  req.IsVerified,
		MinRating:   req.MinRating,
		MaxRating:   req.MaxRating,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Radius:      req.Radius,
		SortBy:      req.SortBy,
	}

	results, total, err := h.poiUsecase.ListPois(ctx, filters)
	if err != nil {
		log.Printf("error listing pois: %v", err)
		return nil, status.Error(codes.Internal, "failed to list pois")
	}

	pois := make([]*tqdpb.Poi, len(results))
	for i, result := range results {
		pois[i] = h.poiMapper.ToProto(result)
	}

	return &tqdpb.ListPoisResponse{
		Data:  pois,
		Total: int64(total),
	}, nil
}
