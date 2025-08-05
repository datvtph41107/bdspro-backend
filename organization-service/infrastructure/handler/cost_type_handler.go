package handler

import (
	_dto "common/dto"
	"context"
	"organization/infrastructure/transformer"
	"organization/internal/dto"
	"organization/internal/usecase"

	pb "pb/types/organization"
	sharedpb "pb/types/shared"
)

type CostTypeHandler struct {
	pb.UnimplementedCostTypeServiceServer
	usecase usecase.CostTypeUsecase
	mapper  *transformer.CostTypeMapper
}

func NewCostTypeHandler(uc usecase.CostTypeUsecase, mapper *transformer.CostTypeMapper) *CostTypeHandler {
	return &CostTypeHandler{usecase: uc, mapper: mapper}
}

// @Summary Tạo loại chi phí
// @Description Tạo loại chi phí
// @Tags CostType
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param costType body pb.CostTypeDTO true "Loại chi phí"
// @Success 200 {object} pb.CostTypeDTO
// @Router /cost-type [post]
func (h *CostTypeHandler) Create(ctx context.Context, req *pb.CostTypeDTO) (*pb.CostTypeDTO, error) {
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
// @Param costType body pb.CostTypeDTO true "Loại chi phí"
// @Success 200 {object} pb.CostTypeDTO
// @Router /cost-type/{id} [put]
func (h *CostTypeHandler) Update(ctx context.Context, req *pb.CostTypeDTO) (*pb.CostTypeDTO, error) {
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
// @Param dto query pb.CostTypeSearchDTO true "Thông tin tìm kiếm"
// @Success 200 {object} pb.CostTypeListDTO
// @Router /cost-type [get]
func (h *CostTypeHandler) List(ctx context.Context, req *pb.CostTypeSearchDTO) (*pb.CostTypeListDTO, error) {
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

	var data []*pb.CostTypeDTO
	for _, ct := range list {
		data = append(data, h.mapper.CostTypeModelToProto(ct))
	}

	return &pb.CostTypeListDTO{
		Data:  data,
		Total: uint32(len(data)),
	}, nil
}
