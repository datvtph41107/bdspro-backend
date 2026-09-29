package handler

import (
	"context"

	"bdspro/infra/mapper"
	admin_usecases "bdspro/internal/usecases/admin"
	bdspropb "pb/types/bdspro"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminAreaRegionHandler struct {
	bdspropb.UnimplementedAdminAreaRegionServiceServer
	usecase *admin_usecases.AreaRegionUsecase
	mapper  *mapper.AreaRegionMapper
}

func NewAdminAreaRegionHandler(
	usecase *admin_usecases.AreaRegionUsecase,
	mapper *mapper.AreaRegionMapper,
) *AdminAreaRegionHandler {
	return &AdminAreaRegionHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

func (h *AdminAreaRegionHandler) CreateAreaRegion(ctx context.Context, req *bdspropb.CreateAreaRegionRequest) (*bdspropb.AreaRegionResponse, error) {
	region := h.mapper.FromCreateRequest(req)
	result, err := h.usecase.Create(ctx, region)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create: %v", err)
	}
	return &bdspropb.AreaRegionResponse{Region: h.mapper.ToProto(result)}, nil
}

func (h *AdminAreaRegionHandler) UpdateAreaRegion(ctx context.Context, req *bdspropb.UpdateAreaRegionRequest) (*bdspropb.AreaRegionResponse, error) {
	region := h.mapper.FromUpdateRequest(req)
	result, err := h.usecase.Update(ctx, region)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update: %v", err)
	}
	return &bdspropb.AreaRegionResponse{Region: h.mapper.ToProto(result)}, nil
}

func (h *AdminAreaRegionHandler) DeleteAreaRegion(ctx context.Context, req *bdspropb.DeleteAreaRegionRequest) (*bdspropb.DeleteAreaRegionResponse, error) {
	if err := h.usecase.Delete(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete: %v", err)
	}
	return &bdspropb.DeleteAreaRegionResponse{Success: true}, nil
}

func (h *AdminAreaRegionHandler) GetAreaRegionByID(ctx context.Context, req *bdspropb.GetAreaRegionByIDRequest) (*bdspropb.AreaRegionResponse, error) {
	region, err := h.usecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "region not found: %v", err)
	}
	return &bdspropb.AreaRegionResponse{Region: h.mapper.ToProto(region)}, nil
}

func (h *AdminAreaRegionHandler) ListAreaRegions(ctx context.Context, req *bdspropb.ListAreaRegionsRequest) (*bdspropb.ListAreaRegionsResponse, error) {
	regions, err := h.usecase.List(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list: %v", err)
	}
	return &bdspropb.ListAreaRegionsResponse{Data: h.mapper.ToProtoList(regions)}, nil
}

func (h *AdminAreaRegionHandler) GetAreaRegionsByIDs(ctx context.Context, req *bdspropb.GetAreaRegionsByIDsRequest) (*bdspropb.GetAreaRegionsByIDsResponse, error) {
	if req == nil || len(req.Ids) == 0 {
		return &bdspropb.GetAreaRegionsByIDsResponse{Regions: []*bdspropb.AreaRegion{}}, nil
	}
	regions, err := h.usecase.GetByIDs(ctx, req.Ids)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get by ids: %v", err)
	}
	return &bdspropb.GetAreaRegionsByIDsResponse{Regions: h.mapper.ToProtoList(regions)}, nil
}
