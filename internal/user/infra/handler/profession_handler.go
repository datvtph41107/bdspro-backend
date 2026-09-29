package handler

import (
	_dto "common/domain/dto"
	"context"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"user/infra/mapper"
	"user/internal/usecases"
)

type ProfessionHandler struct {
	userpb.UnimplementedProfessionServiceServer
	usecase usecases.IAdminProfessionUsecase
	mapper  *mapper.ProfessionMapper
}

func NewProfessionHandler(
	usecase usecases.IAdminProfessionUsecase,
	mapper *mapper.ProfessionMapper,
) *ProfessionHandler {
	return &ProfessionHandler{
		usecase: usecase,
		mapper:  mapper,
	}
}

// ListProfessions lấy danh sách nghề nghiệp
// @Summary Lấy danh sách nghề nghiệp
// @Description Lấy danh sách nghề nghiệp với phân trang (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang hiện tại (bắt đầu từ 0)"
// @Param size query int false "Kích thước trang (tối đa 100)"
// @Router /v2/user/admin/profession [get]
func (h *ProfessionHandler) ListProfessions(ctx context.Context, req *sharepb.RequestV3Proto) (*userpb.ProfessionPage, error) {
	pagable := &_dto.Pagable{
		Page: req.GetPage(),
		Size: req.GetSize(),
	}

	entities, total, err := h.usecase.GetList(ctx, pagable)
	if err != nil {
		return nil, err
	}

	return &userpb.ProfessionPage{
		Data:  h.mapper.EntitiesToProtos(entities),
		Total: uint32(total),
	}, nil
}

// CreateProfession tạo nghề nghiệp mới
// @Summary Tạo nghề nghiệp
// @Description Tạo nghề nghiệp mới (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpb.Profession true "Thông tin nghề nghiệp"
// @Router /v2/user/admin/profession [post]
func (h *ProfessionHandler) CreateProfession(ctx context.Context, req *userpb.Profession) (*userpb.Profession, error) {
	entity := h.mapper.ProtoToEntity(req)

	result, err := h.usecase.Create(ctx, entity)
	if err != nil {
		return nil, err
	}

	return h.mapper.EntityToProto(result), nil
}

// UpdateProfession cập nhật nghề nghiệp
// @Summary Cập nhật nghề nghiệp
// @Description Cập nhật thông tin nghề nghiệp (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID nghề nghiệp"
// @Param request body userpb.Profession true "Thông tin nghề nghiệp"
// @Router /v2/user/admin/profession/{id} [put]
func (h *ProfessionHandler) UpdateProfession(ctx context.Context, req *userpb.Profession) (*userpb.Profession, error) {
	entity := h.mapper.ProtoToEntity(req)

	result, err := h.usecase.Update(ctx, req.GetId(), entity)
	if err != nil {
		return nil, err
	}

	return h.mapper.EntityToProto(result), nil
}

// DeleteProfession xóa nghề nghiệp
// @Summary Xóa nghề nghiệp
// @Description Xóa mềm nghề nghiệp (Admin)
// @Tags AdminDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "ID nghề nghiệp"
// @Router /v2/user/admin/profession/{id} [delete]
func (h *ProfessionHandler) DeleteProfession(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if _, err := h.usecase.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.GetId(),
		Message: "Xóa nghề nghiệp thành công",
	}, nil
}

// GetProfessionItems lấy danh sách nghề nghiệp cho user
// @Summary Lấy danh sách nghề nghiệp
// @Description Lấy danh sách nghề nghiệp đã được duyệt (User)
// @Tags UserDictionary
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Trang hiện tại (bắt đầu từ 0)"
// @Param size query int false "Kích thước trang (tối đa 100)"
// @Router /v2/user/profession [get]
func (h *ProfessionHandler) GetProfessionItems(ctx context.Context, req *sharepb.RequestV3Proto) (*userpb.ProfessionPage, error) {
	pagable := &_dto.Pagable{
		Page: req.GetPage(),
		Size: req.GetSize(),
	}

	entities, total, err := h.usecase.ListVerified(ctx, pagable)
	if err != nil {
		return nil, err
	}

	return &userpb.ProfessionPage{
		Data:  h.mapper.EntitiesToProtos(entities),
		Total: uint32(total),
	}, nil
}
