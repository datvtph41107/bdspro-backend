package handler

import (
	"context"

	hubpb "pb/types/hub"
)

// CreateRegion creates a new region
// @Summary Tạo khu vực mới
// @Description Tạo một region (internal admin)
// @Tags Region
// @Accept json
// @Produce json
// @Param request body hubpb.CreateRegionRequest true "Thông tin region"
// @Success 200 {object} hubpb.RegionResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v2/hub/regions [post]
func (h *InternalHandler) CreateRegion(ctx context.Context, req *hubpb.CreateRegionRequest) (*hubpb.RegionResponse, error) {
	region := h.regionMapper.FromCreateRequest(req)
	result, err := h.regionUsecase.Create(ctx, region)
	if err != nil {
		return nil, err
	}
	return &hubpb.RegionResponse{Region: h.regionMapper.ToProto(result)}, nil
}

// UpdateRegion updates an existing region
// @Summary Cập nhật region
// @Description Cập nhật thông tin region (internal admin)
// @Tags Region
// @Accept json
// @Produce json
// @Param id path int true "Region ID"
// @Param request body hubpb.UpdateRegionRequest true "Thông tin cập nhật"
// @Success 200 {object} hubpb.RegionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v2/hub/regions/{id} [put]
func (h *InternalHandler) UpdateRegion(ctx context.Context, req *hubpb.UpdateRegionRequest) (*hubpb.RegionResponse, error) {
	region := h.regionMapper.FromUpdateRequest(req)
	result, err := h.regionUsecase.Update(ctx, region)
	if err != nil {
		return nil, err
	}
	return &hubpb.RegionResponse{Region: h.regionMapper.ToProto(result)}, nil
}

// DeleteRegion deletes a region by ID
// @Summary Xóa region
// @Description Xóa mềm một region (internal admin)
// @Tags Region
// @Accept json
// @Produce json
// @Param id path int true "Region ID"
// @Success 200 {object} hubpb.DeleteRegionResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /v2/hub/regions/{id} [delete]
func (h *InternalHandler) DeleteRegion(ctx context.Context, req *hubpb.DeleteRegionRequest) (*hubpb.DeleteRegionResponse, error) {
	err := h.regionUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &hubpb.DeleteRegionResponse{Success: true}, nil
}

// GetRegionByID retrieves a single region
// @Summary Lấy region theo ID
// @Description Lấy thông tin một region
// @Tags Region
// @Accept json
// @Produce json
// @Param id path int true "Region ID"
// @Success 200 {object} hubpb.RegionResponse
// @Failure 404 {object} map[string]string
// @Router /v2/hub/regions/{id} [get]
func (h *InternalHandler) GetRegionByID(ctx context.Context, req *hubpb.GetRegionByIDRequest) (*hubpb.RegionResponse, error) {
	region, err := h.regionUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &hubpb.RegionResponse{Region: h.regionMapper.ToProto(region)}, nil
}

// ListRegions returns all active regions
// @Summary Danh sách region
// @Description Lấy tất cả region (active)
// @Tags Region
// @Accept json
// @Produce json
// @Success 200 {object} hubpb.ListRegionsResponse
// @Router /v2/hub/regions [get]
func (h *InternalHandler) ListRegions(ctx context.Context, req *hubpb.ListRegionsRequest) (*hubpb.ListRegionsResponse, error) {
	regions, err := h.regionUsecase.List(ctx)
	if err != nil {
		return nil, err
	}
	return &hubpb.ListRegionsResponse{Regions: h.regionMapper.ToProtoList(regions)}, nil
}
