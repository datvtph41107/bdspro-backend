package handler

import (
	_dto "common/domain/dto"
	"context"
	bdspropb "pb/types/bdspro"
	sharedpb "pb/types/shared"

	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
)

type TxCostTypeHandler struct {
	bdspropb.UnimplementedCostTypeServiceServer
	usecase usecases.TxCostTypeUsecase
	mapper  *mapper.TxCostTypeMapper
}

func NewTxCostTypeHandler(uc usecases.TxCostTypeUsecase, mapper *mapper.TxCostTypeMapper) *TxCostTypeHandler {
	return &TxCostTypeHandler{usecase: uc, mapper: mapper}
}

// @Summary Tạo loại chi phí
// @Description Tạo loại chi phí
// @Tags CostType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param costType body bdspropb.CostTypeDTO true "Loại chi phí"
// @Success 200 {object} bdspropb.CostTypeDTO
// @Router /cost-type [post]
func (h *TxCostTypeHandler) Create(ctx context.Context, req *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	model := h.mapper.TxCostTypeProtoToModel(req)
	if err := h.usecase.Create(ctx, model); err != nil {
		return nil, err
	}
	return h.mapper.TxCostTypeModelToProto(model), nil
}

// @Summary Cập nhật loại chi phí
// @Description Cập nhật loại chi phí
// @Tags CostType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID loại chi phí"
// @Param costType body bdspropb.CostTypeDTO true "Loại chi phí"
// @Success 200 {object} bdspropb.CostTypeDTO
// @Router /cost-type/{id} [put]
func (h *TxCostTypeHandler) Update(ctx context.Context, req *bdspropb.CostTypeDTO) (*bdspropb.CostTypeDTO, error) {
	model := h.mapper.TxCostTypeProtoToModel(req)
	if err := h.usecase.Update(ctx, model); err != nil {
		return nil, err
	}
	return h.mapper.TxCostTypeModelToProto(model), nil
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
func (h *TxCostTypeHandler) Delete(ctx context.Context, req *sharedpb.IdRequest) (*sharedpb.Empty, error) {
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
// @Param dto query bdspropb.CostTypeSearchDTO true "Thông tin tìm kiếm"
// @Success 200 {object} bdspropb.CostTypeListDTO
// @Router /cost-type [get]
func (h *TxCostTypeHandler) List(ctx context.Context, req *bdspropb.CostTypeSearchDTO) (*bdspropb.CostTypeListDTO, error) {
	search := &dto.TxCostTypeSearchDTO{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		TypeName: req.TypeName,
		CostType: req.CostType,
	}
	list, err := h.usecase.List(ctx, search)
	if err != nil {
		return nil, err
	}

	var data []*bdspropb.CostTypeDTO
	for _, ct := range list {
		data = append(data, h.mapper.TxCostTypeModelToProto(ct))
	}

	return &bdspropb.CostTypeListDTO{
		Data:  data,
		Total: uint32(len(data)),
	}, nil
}
