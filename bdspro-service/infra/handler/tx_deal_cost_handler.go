package handler

import (
	_dto "common/domain/dto"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"

	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
)

type TxDealCostHandler struct {
	bdspropb.UnimplementedDealCostServiceServer
	usecase usecases.DealCostUsecase
	mapper  *mapper.DealCostMapper
}

func NewDealCostHandler(uc usecases.DealCostUsecase, mapper *mapper.DealCostMapper) *TxDealCostHandler {
	return &TxDealCostHandler{usecase: uc, mapper: mapper}
}

// @Summary Tạo Chi phí/Doanh thu thương vụ
// @Description Tạo Chi phí/Doanh thu thương vụ
// @Tags DealCost
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dealCost body bdspropb.DealCostDTO true "Deal cost data"
// @Success 200 {object} bdspropb.DealCostDTO
// @Router /deal-cost [post]
func (h *TxDealCostHandler) Create(ctx context.Context, req *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error) {
	model := h.mapper.DealCostProtoToModel(req)
	if err := h.usecase.Create(ctx, model); err != nil {
		return nil, err
	}
	return h.mapper.DealCostModelToProto(model), nil
}

// @Summary Cập nhật Chi phí/Doanh thu thương vụ
// @Description Cập nhật Chi phí/Doanh thu thương vụ
// @Tags DealCost
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dealCost body bdspropb.DealCostDTO true "Deal cost data"
// @Success 200 {object} bdspropb.DealCostDTO
// @Router /deal-cost/{id} [put]
func (h *TxDealCostHandler) Update(ctx context.Context, req *bdspropb.DealCostDTO) (*bdspropb.DealCostDTO, error) {
	model := h.mapper.DealCostProtoToModel(req)
	if err := h.usecase.Update(ctx, model); err != nil {
		return nil, err
	}
	return h.mapper.DealCostModelToProto(model), nil
}

// @Summary Xóa Chi phí/Doanh thu thương vụ
// @Description Xóa Chi phí/Doanh thu thương vụ
// @Tags DealCost
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "Deal cost ID"
// @Success 200 {object} sharepb.Empty
// @Router /deal-cost/{id} [delete]
func (h *TxDealCostHandler) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	err := h.usecase.Delete(ctx, req.Id)
	if err != nil {
		return &sharepb.Empty{}, nil
	}
	return &sharepb.Empty{}, nil
}

// @Summary Lấy Chi phí/Doanh thu thương vụ theo ID
// @Description Lấy Chi phí/Doanh thu thương vụ theo ID
// @Tags DealCost
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "Deal cost ID"
// @Success 200 {object} bdspropb.DealCostDTO
// @Router /deal-cost/detail/{id} [get]
func (h *TxDealCostHandler) Detail(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.DealCostDTO, error) {
	model, err := h.usecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.mapper.DealCostModelToProto(model), nil
}

// @Summary Tìm kiếm Chi phí/Doanh thu thương vụ
// @Description Tìm kiếm Chi phí/Doanh thu thương vụ
// @Tags DealCost
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dealId path uint64 true "Deal ID"
// @Param dto query bdspropb.DealCostSearchDTO true "Thông tin tìm kiếm"
// @Router /deal/deal-cost/{dealId} [get]
func (h *TxDealCostHandler) Search(ctx context.Context, req *bdspropb.DealCostSearchDTO) (*bdspropb.DealCostListDTO, error) {
	search := &dto.TxDealCostSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		DealId: req.DealId,
	}
	list, err := h.usecase.ListByDealID(ctx, search)
	if err != nil {
		return nil, err
	}

	var data []*bdspropb.DealCostDTO
	for _, item := range list {
		data = append(data, h.mapper.DealCostModelToProto(item))
	}

	return &bdspropb.DealCostListDTO{
		Data:          data,
		TotalElements: uint32(len(data)),
	}, nil
}
