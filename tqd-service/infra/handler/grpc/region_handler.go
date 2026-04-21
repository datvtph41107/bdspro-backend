package handler_grpc

import (
	_errors "common/errors"
	"context"

	sharepb "pb/types/shared"
	tqdpb "pb/types/tqd"
	"tqd/infra/mapper"
	"tqd/internal/enums"
	"tqd/internal/usecase"
)

type RegionGrpcHandler struct {
	tqdpb.UnimplementedRegionServiceServer
	regionUsecase usecase.RegionUsecase
	regionMapper  *mapper.RegionMapper
}

func NewRegionGrpcHandler(regionUsecase usecase.RegionUsecase) *RegionGrpcHandler {
	return &RegionGrpcHandler{
		regionUsecase: regionUsecase,
		regionMapper:  mapper.NewRegionMapper(),
	}
}

func (h *RegionGrpcHandler) CreateRegion(ctx context.Context, req *tqdpb.CreateRegionRequest) (*tqdpb.RegionResponse, error) {
	// Validate
	if req.Name == "" {
		return nil, _errors.ReturnError(400, "name is required")
	}
	if req.DisplayName == "" {
		return nil, _errors.ReturnError(400, "displayName is required")
	}
	if req.Geometry == "" {
		return nil, _errors.ReturnError(400, "geometry is required")
	}
	if req.Type < 1 || req.Type > 4 {
		return nil, _errors.ReturnError(400, "type must be 1-4")
	}
	if req.Status < 10 || req.Status > 30 {
		return nil, _errors.ReturnError(400, "status must be 10,20,30")
	}

	createReq := h.regionMapper.ToCreateRegionDTO(req)

	region, err := h.regionUsecase.Create(ctx, createReq)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return h.regionMapper.ToProtoFromDomain(region), nil
}

func (h *RegionGrpcHandler) UpdateRegion(ctx context.Context, req *tqdpb.UpdateRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	// Convert to DTO using mapper
	updateReq := h.regionMapper.ToUpdateRegionDTO(req)

	region, err := h.regionUsecase.Update(ctx, uint(req.Id), updateReq)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return h.regionMapper.ToProtoFromDomain(region), nil
}

func (h *RegionGrpcHandler) DeleteRegion(ctx context.Context, req *tqdpb.DeleteRegionRequest) (*sharepb.Empty, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	if err := h.regionUsecase.Delete(ctx, uint(req.Id)); err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}
	return &sharepb.Empty{}, nil
}

func (h *RegionGrpcHandler) GetRegion(ctx context.Context, req *tqdpb.GetRegionRequest) (*tqdpb.RegionResponse, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	region, err := h.regionUsecase.GetByID(ctx, uint(req.Id))
	if err != nil {
		return nil, _errors.ReturnError(404, err.Error())
	}
	if region == nil {
		return nil, _errors.ReturnError(404, "region not found")
	}

	return h.regionMapper.ToProtoFromDTO(region), nil
}

func (h *RegionGrpcHandler) ListRegions(ctx context.Context, req *tqdpb.ListRegionsRequest) (*tqdpb.ListRegionsResponse, error) {
	// Convert to DTO using mapper
	listReq := h.regionMapper.ToListRegionDTO(req)

	resp, err := h.regionUsecase.List(ctx, listReq)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return &tqdpb.ListRegionsResponse{
		Data:  h.regionMapper.ToProtoListFromDTO(resp.Data),
		Total: resp.Total,
		Page:  int32(resp.Page),
		Limit: int32(resp.Limit),
	}, nil
}

func (h *RegionGrpcHandler) GetRegionByLocation(ctx context.Context, req *tqdpb.GetRegionByLocationRequest) (*tqdpb.GetRegionByLocationResponse, error) {
	var regionType *enums.RegionType
	if req.Type != nil && *req.Type > 0 {
		t := enums.RegionType(*req.Type)
		regionType = &t
	}

	var layerID *uint
	if req.LayerId != nil && *req.LayerId > 0 {
		id := uint(*req.LayerId)
		layerID = &id
	}

	regions, err := h.regionUsecase.GetByLocation(ctx, req.Lat, req.Lng, regionType, layerID)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return &tqdpb.GetRegionByLocationResponse{
		Data: h.regionMapper.ToProtoListFromDomain(regions),
	}, nil
}

func (h *RegionGrpcHandler) ImportGeoJSON(ctx context.Context, req *tqdpb.ImportGeoJSONRequest) (*tqdpb.ImportGeoJSONResponse, error) {
	if req.Name == "" {
		return nil, _errors.ReturnError(400, "name is required")
	}
	if req.GeoJSON == "" {
		return nil, _errors.ReturnError(400, "geoJSON is required")
	}

	// Convert to DTO using mapper
	importReq := h.regionMapper.ToImportGeoJSONDTO(req)

	resp, err := h.regionUsecase.ImportGeoJSON(ctx, importReq)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return &tqdpb.ImportGeoJSONResponse{
		Region:       h.regionMapper.ToProtoFromDTO(&resp.Region),
		FeatureCount: int32(resp.FeatureCount),
	}, nil
}

func (h *RegionGrpcHandler) GetRegionIntersections(ctx context.Context, req *tqdpb.GetRegionIntersectionsRequest) (*tqdpb.GetRegionIntersectionsResponse, error) {
	if req.Id == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}

	// Convert to DTO using mapper
	intersectReq := h.regionMapper.ToIntersectionQueryDTO(req)

	regions, err := h.regionUsecase.GetIntersections(ctx, intersectReq)
	if err != nil {
		return nil, _errors.ReturnError(500, err.Error())
	}

	return &tqdpb.GetRegionIntersectionsResponse{
		Data: h.regionMapper.ToProtoListFromDTO(regions),
	}, nil
}
