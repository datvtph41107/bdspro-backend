package handler

import (
	_dto "common/domain/dto"
	"common/fault"
	"context"
	authpb "pb/types/auth"
	sharepb "pb/types/shared"
	"user/infra/mapper"
	"user/internal/enums"
	"user/internal/usecase"
	"user/validator"
)

type RoleGroupHandler struct {
	authpb.UnimplementedRoleGroupServiceServer
	RoleGroupUsecase   *usecase.RoleGroupUsecase
	Mapper             *mapper.RoleGroupMapper
	RoleMapper         *mapper.RoleMapper
	RoleGroupValidator *validator.RoleGroupValidator
	InternalHandler    *InternalHandler
}

func NewRoleGroupHandler(
	roleGroupUsecase *usecase.RoleGroupUsecase,
	mapper *mapper.RoleGroupMapper,
	roleMapper *mapper.RoleMapper,
	roleGroupValidator *validator.RoleGroupValidator,
	internalHandler *InternalHandler,
) *RoleGroupHandler {
	return &RoleGroupHandler{
		RoleGroupUsecase:   roleGroupUsecase,
		Mapper:             mapper,
		RoleMapper:         roleMapper,
		RoleGroupValidator: roleGroupValidator,
		InternalHandler:    internalHandler,
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
// @Router /role-group/list [get]
func (h *RoleGroupHandler) GetListRoleGroup(ctx context.Context, req *sharepb.IdRequest) (*authpb.RoleGroupListResponse, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{permissionIAMRoleView}); err != nil {
		return nil, err
	}
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
// @Router /role-group [post]
func (h *RoleGroupHandler) CreateRoleGroup(ctx context.Context, req *authpb.RoleGroupRequest) (*authpb.RoleGroup, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{permissionIAMRoleManage}); err != nil {
		return nil, err
	}
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
// @Router /role-group/{id} [put]
func (h *RoleGroupHandler) UpdateRoleGroup(ctx context.Context, req *authpb.RoleGroupRequest) (*authpb.RoleGroup, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{permissionIAMRoleManage}); err != nil {
		return nil, err
	}
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
// @Router /role-group/{id} [delete]
func (h *RoleGroupHandler) DeleteRoleGroup(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{permissionIAMRoleManage}); err != nil {
		return nil, err
	}
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
// @Router /role-group/detail/{id} [get]
func (h *RoleGroupHandler) GetRoleGroupById(ctx context.Context, req *sharepb.IdRequest) (*authpb.RoleGroup, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{permissionIAMRoleView}); err != nil {
		return nil, err
	}
	result, err := h.RoleGroupUsecase.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.Mapper.MapToPb(result), nil
}

func mapUpdateGroupPermissionsError(err error) error {
	if _, ok := fault.As(err); !ok {
		return err
	}
	return fault.ToGRPC(err)
}

// @Summary Cập nhật permission cho group
// @Description Cập nhật toàn bộ permission cho group
// @Tags RoleGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body authpb.UpdateGroupPermissionsRequest true "Thông tin cập nhật permission"
// @Success 200 {object} sharepb.SubmitResponse
// @Router /role-group/permissions [put]
func (h *RoleGroupHandler) UpdateGroupPermissions(ctx context.Context, req *authpb.UpdateGroupPermissionsRequest) (*sharepb.SubmitResponse, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{permissionIAMRoleManage}); err != nil {
		return nil, err
	}
	err := h.RoleGroupUsecase.UpdateGroupPermissions(ctx, req.GroupId, req.PermissionIds)
	if err != nil {
		return nil, mapUpdateGroupPermissionsError(err)
	}

	return &sharepb.SubmitResponse{
		Id:      req.GroupId,
		Message: "Cập nhật permission cho group thành công",
	}, nil
}

// @Summary Lấy danh sách permission theo group key
// @Description Lấy danh sách permission theo group key
// @Tags RoleGroup
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param groupKey path int true "Group key"
// @Param page query int false "Page"
// @Param size query int false "Size"
// @Success 200 {object} authpb.PermissionListResponse
// @Router /role-group/{groupKey}/permissions [get]
func (h *RoleGroupHandler) GetPermissionsByGroupKey(ctx context.Context, req *authpb.GetPermissionsByGroupKeyRequest) (*authpb.PermissionListResponse, error) {
	if err := h.InternalHandler.HasPermissions(ctx, []string{permissionIAMRoleView}); err != nil {
		return nil, err
	}
	permissions, err := h.RoleGroupUsecase.GetPermissionsByGroupKey(ctx, enums.GroupRoleKeyEnum(req.GroupKey))
	if err != nil {
		return nil, err
	}

	// Sử dụng mapper giống như GetAllPermissions
	pbPermissions := make([]*sharepb.Permission, len(permissions))
	for i, permission := range permissions {
		pbPermissions[i] = h.RoleMapper.MapPermissionListWithChild(permission)
	}

	return &authpb.PermissionListResponse{
		Data:          pbPermissions,
		TotalElements: 0,
	}, nil
}
