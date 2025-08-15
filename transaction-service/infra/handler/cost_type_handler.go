package handler

import (
	_dto "common/dto"
	"context"
	sharedpb "pb/types/shared"
	transactionpb "pb/types/transaction"
	"tx/infra/mapper"
	"tx/internal/dto"
	transaction_usecase "tx/internal/usecases"
)

type CostTypeHandler struct {
	transactionpb.UnimplementedCostTypeServiceServer
	usecase transaction_usecase.CostTypeUsecase
	mapper  *mapper.CostTypeMapper
}

func NewCostTypeHandler(uc transaction_usecase.CostTypeUsecase, mapper *mapper.CostTypeMapper) *CostTypeHandler {
	return &CostTypeHandler{usecase: uc, mapper: mapper}
}

// @Summary Tạo loại chi phí
// @Description Tạo loại chi phí
// @Tags CostType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param costType body transactionpb.CostTypeDTO true "Loại chi phí"
// @Success 200 {object} transactionpb.CostTypeDTO
// @Router /cost-type [post]
func (h *CostTypeHandler) Create(ctx context.Context, req *transactionpb.CostTypeDTO) (*transactionpb.CostTypeDTO, error) {
	model := h.mapper.CostTypeProtoToModel(req)
	if err := h.usecase.Create(ctx, model); err != nil {
		return nil, err
	}
	return h.mapper.CostTypeModelToProto(model), nil
}

// @Summary Cập nhật loại chi phí
// @Description Cập nhật loại chi phí
// @Tags CostType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID loại chi phí"
// @Param costType body transactionpb.CostTypeDTO true "Loại chi phí"
// @Success 200 {object} transactionpb.CostTypeDTO
// @Router /cost-type/{id} [put]
func (h *CostTypeHandler) Update(ctx context.Context, req *transactionpb.CostTypeDTO) (*transactionpb.CostTypeDTO, error) {
	model := h.mapper.CostTypeProtoToModel(req)
	if err := h.usecase.Update(ctx, model); err != nil {
		return nil, err
	}
	return h.mapper.CostTypeModelToProto(model), nil
}

// @Summary Xóa loại chi phí
// @Description Xóa loại chi phí
// @Tags CostType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID loại chi phí"
// @Success 200 {object} sharedpb.Empty
// @Router /cost-type/{id} [delete]
func (h *CostTypeHandler) Delete(ctx context.Context, req *sharedpb.IdRequest) (*sharedpb.Empty, error) {
	err := h.usecase.Delete(ctx, req.Id)
	if err != nil {
		return &sharedpb.Empty{}, err
	}
	return &sharedpb.Empty{}, nil
}

// @Summary Lấy danh sách loại chi phí
// @Description Lấy danh sách loại chi phí
// @Tags CostType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dto query transactionpb.CostTypeSearchDTO true "Thông tin tìm kiếm"
// @Success 200 {object} transactionpb.CostTypeListDTO
// @Router /cost-type [get]
func (h *CostTypeHandler) List(ctx context.Context, req *transactionpb.CostTypeSearchDTO) (*transactionpb.CostTypeListDTO, error) {
	search := &dto.CostTypeSearchDTO{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
		TypeName: req.TypeName,
		CostType: req.CostType,
	}
	list, err := h.usecase.List(ctx, search)
	if err != nil {
		return nil, err
	}

	var data []*transactionpb.CostTypeDTO
	for _, ct := range list {
		data = append(data, h.mapper.CostTypeModelToProto(ct))
	}

	return &transactionpb.CostTypeListDTO{
		Data:  data,
		Total: uint32(len(data)),
	}, nil
}
