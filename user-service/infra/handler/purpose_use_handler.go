package handler

import (
	"context"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"
)

type PurposeUseHandler struct {
	userpb.UnimplementedPurposeUseServiceServer
	purposeUsecase usecases.IAdminPurposeUseUsecase
	mapper         *mapper.PurposeUseMapper
}

func NewPurposeUseHandler(
	purposeUsecase usecases.IAdminPurposeUseUsecase,
	mapper *mapper.PurposeUseMapper,
) *PurposeUseHandler {
	return &PurposeUseHandler{
		purposeUsecase: purposeUsecase,
		mapper:         mapper,
	}
}

// GetPurposeUses lấy danh sách mục đích sử dụng
// @Summary Lấy danh sách mục đích sử dụng
// @Description Lấy danh sách mục đích sử dụng (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sharepb.IdRequest true "Thông tin mục đích sử dụng"
// @Success 200 {object} userpb.PurposeUsePage
// @Router /v2/user/admin/purpose-use [get]
func (h *PurposeUseHandler) GetPurposeUses(ctx context.Context, req *sharepb.IdRequest) (*userpb.PurposeUsePage, error) {
	entities, err := h.purposeUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	data := h.mapper.EntitiesToProtos(entities)

	return &userpb.PurposeUsePage{
		Data:  data,
		Total: uint32(len(data)),
	}, nil
}

// CreatePurposeUse tạo mục đích sử dụng mới
// @Summary Tạo mục đích sử dụng mới
// @Description Tạo mục đích sử dụng mới (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpb.PurposeUse true "Thông tin mục đích sử dụng"
// @Success 200 {object} userpb.PurposeUse
// @Router /v2/user/admin/purpose-use [post]
func (h *PurposeUseHandler) CreatePurposeUse(ctx context.Context, req *userpb.PurposeUse) (*userpb.PurposeUse, error) {
	entity := h.mapper.ProtoToEntity(req)

	result, err := h.purposeUsecase.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return h.mapper.EntityToProto(result), nil
}

// UpdatePurposeUse cập nhật mục đích sử dụng
// @Summary Cập nhật mục đích sử dụng
// @Description Cập nhật mục đích sử dụng (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID mục đích sử dụng"
// @Param request body userpb.PurposeUse true "Thông tin mục đích sử dụng"
// @Success 200 {object} userpb.PurposeUse
// @Router /v2/user/admin/purpose-use/{id} [put]
func (h *PurposeUseHandler) UpdatePurposeUse(ctx context.Context, req *userpb.PurposeUse) (*userpb.PurposeUse, error) {
	entity, err := h.purposeUsecase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	entity.Name = req.GetName()
	if req.Description != nil {
		entity.Description = req.GetDescription()
	}

	updated, err := h.purposeUsecase.Update(ctx, req.GetId(), entity)
	if err != nil {
		return nil, err
	}

	return h.mapper.EntityToProto(updated), nil
}

// DeletePurposeUse xóa mục đích sử dụng
// @Summary Xóa mục đích sử dụng
// @Description Xóa mềm mục đích sử dụng (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID mục đích sử dụng"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /v2/user/admin/purpose-use/{id} [delete]
func (h *PurposeUseHandler) DeletePurposeUse(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if _, err := h.purposeUsecase.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.GetId(),
		Message: "Xóa mục đích sử dụng thành công",
	}, nil
}

// GetPurposeUseItems lấy danh sách mục đích sử dụng
// @Summary Lấy danh sách mục đích sử dụng
// @Description Lấy danh sách mục đích sử dụng (User)
// @Tags UserDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body sharepb.IdRequest true "Thông tin mục đích sử dụng"
// @Success 200 {object} sharepb.ItemV3PageProto
// @Router /v2/user/purpose-use [get]
func (h *PurposeUseHandler) GetPurposeUseItems(ctx context.Context, req *sharepb.IdRequest) (*sharepb.ItemV3PageProto, error) {
	entities, err := h.purposeUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	data := h.mapper.EntitiesToItemProtos(entities)

	return &sharepb.ItemV3PageProto{
		Data:  data,
		Total: uint32(len(data)),
	}, nil
}
