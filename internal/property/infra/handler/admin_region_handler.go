package handler

import (
	"bdspro/infra/mapper"
	admin_usecases "bdspro/internal/usecases/admin"
	"context"

	bdspropb "pb/types/bdspro"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminRegionHandler struct {
	bdspropb.UnimplementedAdminRegionServiceServer
	RegionUsecase *admin_usecases.RegionUsecase
	RegionMapper  *mapper.RegionMapper
}

func NewAdminRegionHandler(
	regionUsecase *admin_usecases.RegionUsecase,
	regionMapper *mapper.RegionMapper) *AdminRegionHandler {
	return &AdminRegionHandler{
		RegionUsecase: regionUsecase,
		RegionMapper:  regionMapper,
	}
}

// @Summary Lấy danh sách địa chỉ (tỉnh/huyện/xã)
// @Description Lấy danh sách địa chỉ theo level: 1=Tỉnh/TP, 2=Quận/Huyện, 3=Xã/Phường
// @Tags Admin - Region
// @Accept json
// @Produce json
// @Param page query int false "Trang hiện tại" default(1)
// @Param size query int false "Số lượng item" default(20)
// @Param text query string false "Tìm kiếm theo tên"
// @Param level query int false "Cấp địa chỉ: 1=Tỉnh/TP, 2=Quận/Huyện, 3=Xã/Phường"
// @Param parentId query int false "ID của cấp cha"
// @Success 200 {object} bdspropb.RegionListResponse
// @Router /v2/bdspro/admin/region/list [get]
func (h *AdminRegionHandler) GetRegionList(ctx context.Context, req *bdspropb.RegionListRequest) (*bdspropb.RegionListResponse, error) {
	// Convert proto request to DTO
	dtoReq := h.RegionMapper.PbToRegionListRequest(req)

	// Call usecase
	regions, total, err := h.RegionUsecase.GetRegionList(ctx, dtoReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get region list: %v", err)
	}

	// Convert to proto response
	response := &bdspropb.RegionListResponse{
		Data:  h.RegionMapper.RegionListToProto(regions),
		Total: total,
	}

	return response, nil
}

// @Summary Lấy chi tiết địa chỉ
// @Description Lấy thông tin chi tiết của một địa chỉ
// @Tags Admin - Region
// @Accept json
// @Produce json
// @Param id path int true "ID của địa chỉ"
// @Success 200 {object} bdspropb.Region
// @Router /v2/bdspro/admin/region/{id} [get]
func (h *AdminRegionHandler) GetRegionDetail(ctx context.Context, req *bdspropb.RegionDetailRequest) (*bdspropb.Region, error) {
	// Call usecase
	region, err := h.RegionUsecase.GetRegionDetail(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "region not found: %v", err)
	}

	// Convert to proto response
	return h.RegionMapper.RegionToProto(region), nil
}
