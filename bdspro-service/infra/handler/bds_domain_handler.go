package handler

import (
	"bdspro/infra/mapper"
	"bdspro/internal/usecases"
	"context"
	bdspropb "pb/types/bdspro"
)

type BdsDomainHandler struct {
	bdspropb.UnimplementedBdsDomainServiceServer
	BdsDomainUC     *usecases.BdsDomainUsecase
	BdsDomainMapper *mapper.BdsDomainMapper
}

func NewBdsDomainHandler(
	bdsDomainUC *usecases.BdsDomainUsecase,
	bdsDomainMapper *mapper.BdsDomainMapper,
) *BdsDomainHandler {
	return &BdsDomainHandler{
		BdsDomainUC:     bdsDomainUC,
		BdsDomainMapper: bdsDomainMapper,
	}
}

// @Summary Lấy danh sách bất động sản
// @Description Lấy danh sách bất động sản với phân trang và bộ lọc
// @Tags User: Bất động sản
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang hiện tại" default(1)
// @Param size query int false "Số lượng item" default(20)
// @Param sort query string false "Sắp xếp"
// @Param text query string false "Tìm kiếm theo text"
// @Param productId query int false "Lọc theo Product ID"
// @Param assetId query int false "Lọc theo Asset ID"
// @Param provinceId query int false "Lọc theo Tỉnh/TP"
// @Param wardId query int false "Lọc theo Phường/Xã"
// @Param districtId query int false "Lọc theo Quận/Huyện"
// @Success 200 {object} bdspropb.BdsDomainListResponse
// @Router /v2/bdspro/v2/bds/list [get]
func (h *BdsDomainHandler) GetBdsDomainList(ctx context.Context, req *bdspropb.BdsDomainListRequest) (*bdspropb.BdsDomainListResponse, error) {
	// Convert protobuf request to DTO
	dto := h.BdsDomainMapper.PbToBdsDomainListRequest(req)

	// Get list from usecase
	results, total, err := h.BdsDomainUC.GetList(ctx, dto)
	if err != nil {
		return nil, err
	}

	// Convert domain to protobuf
	bdsList := h.BdsDomainMapper.BdsDomainListToProto(results)

	return &bdspropb.BdsDomainListResponse{
		Data:  bdsList,
		Total: total,
	}, nil
}
