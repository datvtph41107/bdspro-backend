package handler

import (
	"auth/infra/mapper"
	"auth/internal/usecase"
	"auth/validator"
	_dto "common/dto"
	"context"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
)

type RoleGroupHandler struct {
	authpb.UnimplementedRoleGroupServiceServer
	RoleGroupUsecase   *usecase.RoleGroupUsecase
	Mapper             *mapper.RoleGroupMapper
	RoleGroupValidator *validator.RoleGroupValidator
}

func NewRoleGroupHandler(
	roleGroupUsecase *usecase.RoleGroupUsecase,
	mapper *mapper.RoleGroupMapper,
	roleGroupValidator *validator.RoleGroupValidator,
) *RoleGroupHandler {
	return &RoleGroupHandler{
		RoleGroupUsecase:   roleGroupUsecase,
		Mapper:             mapper,
		RoleGroupValidator: roleGroupValidator,
	}
}

// @Summary Lấy danh sách role group
// @Description Lấy danh sách role group với phân trang
// @Tags RoleGroup
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Success 200 {object} authpb.RoleGroupListResponse
// @Router /v2/auth/role-group/list [get]
func (h *RoleGroupHandler) GetListRoleGroup(ctx context.Context, req *sharepb.IdRequest) (*authpb.RoleGroupListResponse, error) {
	result, total, err := h.RoleGroupUsecase.GetList(ctx, &_dto.Pagable{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}

	// Convert DTO to protobuf
	response := &authpb.RoleGroupListResponse{
		Data:          h.Mapper.MapToPbList(result),
		TotalElements: uint32(total),
	}

	return response, nil
}

// @Summary Tạo mới role group
// @Description Tạo mới role group
// @Tags RoleGroup
// @Accept json
// @Produce json
// @Param body body authpb.RoleGroupRequest true "Thông tin role group"
// @Success 200 {object} authpb.RoleGroup
// @Router /v2/auth/role-group [post]
func (h *RoleGroupHandler) CreateRoleGroup(ctx context.Context, req *authpb.RoleGroupRequest) (*authpb.RoleGroup, error) {
	domain := h.Mapper.MapToDomain(req)

	err := h.RoleGroupValidator.ValidateCreateRoleGroup(req)
	if err != nil {
		return nil, err
	}

	result, err := h.RoleGroupUsecase.Create(ctx, domain)
	if err != nil {
		return nil, err
	}

	return h.Mapper.MapToPb(result), nil
}

// @Summary Cập nhật role group
// @Description Cập nhật role group
// @Tags RoleGroup
// @Accept json
// @Produce json
// @Param id path int true "ID role group"
// @Param body body authpb.RoleGroupRequest true "Thông tin role group"
// @Success 200 {object} authpb.RoleGroup
// @Router /v2/auth/role-group/{id} [put]
func (h *RoleGroupHandler) UpdateRoleGroup(ctx context.Context, req *authpb.RoleGroupRequest) (*authpb.RoleGroup, error) {
	domain := h.Mapper.MapToDomain(req)

	err := h.RoleGroupValidator.ValidateUpdateRoleGroup(req)
	if err != nil {
		return nil, err
	}

	result, err := h.RoleGroupUsecase.Update(ctx, req.Id, domain)
	if err != nil {
		return nil, err
	}

	return h.Mapper.MapToPb(result), nil
}

// @Summary Xóa role group
// @Description Xóa role group
// @Tags RoleGroup
// @Accept json
// @Produce json
// @Param id path int true "ID role group"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /v2/auth/role-group/{id} [delete]
func (h *RoleGroupHandler) DeleteRoleGroup(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	_, err := h.RoleGroupUsecase.Delete(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Xóa role group thành công",
	}, nil
}

// @Summary Lấy chi tiết role group
// @Description Lấy chi tiết role group theo ID
// @Tags RoleGroup
// @Accept json
// @Produce json
// @Param id path int true "ID role group"
// @Success 200 {object} authpb.RoleGroup
// @Router /v2/auth/role-group/detail/{id} [get]
func (h *RoleGroupHandler) GetRoleGroupById(ctx context.Context, req *sharepb.IdRequest) (*authpb.RoleGroup, error) {
	result, err := h.RoleGroupUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.Mapper.MapToPb(result), nil
}
