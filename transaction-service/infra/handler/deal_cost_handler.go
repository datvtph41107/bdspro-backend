package handler

import (
	_dto "common/dto"
	"context"
	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"

	"tx/infra/mapper"
	"tx/internal/dto"
	transaction_usecase "tx/internal/usecases"
)

type DealCostHandler struct {
	transactionpb.UnimplementedDealCostServiceServer
	usecase transaction_usecase.DealCostUsecase
	mapper  *mapper.DealCostMapper
}

func NewDealCostHandler(uc transaction_usecase.DealCostUsecase, mapper *mapper.DealCostMapper) *DealCostHandler {
	return &DealCostHandler{usecase: uc, mapper: mapper}
}

// @Summary Tạo Chi phí/Doanh thu thương vụ
// @Description Tạo Chi phí/Doanh thu thương vụ
// @Tags DealCost
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dealCost body transactionpb.DealCostDTO true "Deal cost data"
// @Success 200 {object} transactionpb.DealCostDTO
// @Router /deal-cost [post]
func (h *DealCostHandler) Create(ctx context.Context, req *transactionpb.DealCostDTO) (*transactionpb.DealCostDTO, error) {
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
// @Param dealCost body transactionpb.DealCostDTO true "Deal cost data"
// @Success 200 {object} transactionpb.DealCostDTO
// @Router /deal-cost/{id} [put]
func (h *DealCostHandler) Update(ctx context.Context, req *transactionpb.DealCostDTO) (*transactionpb.DealCostDTO, error) {
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
func (h *DealCostHandler) Delete(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
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
// @Success 200 {object} transactionpb.DealCostDTO
// @Router /deal-cost/detail/{id} [get]
func (h *DealCostHandler) Detail(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.DealCostDTO, error) {
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
// @Param dto query transactionpb.DealCostSearchDTO true "Thông tin tìm kiếm"
// @Router /deal/deal-cost/{dealId} [get]
func (h *DealCostHandler) Search(ctx context.Context, req *transactionpb.DealCostSearchDTO) (*transactionpb.DealCostListDTO, error) {
	search := &dto.DealCostSearchDTO{
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

	var data []*transactionpb.DealCostDTO
	for _, item := range list {
		data = append(data, h.mapper.DealCostModelToProto(item))
	}

	return &transactionpb.DealCostListDTO{
		Data:          data,
		TotalElements: uint32(len(data)),
	}, nil
}
