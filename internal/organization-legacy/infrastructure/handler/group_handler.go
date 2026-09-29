package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"organization/infrastructure/client"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
)

type GroupHandler struct {
	organizationpb.UnimplementedGroupServiceServer
	GroupUsecase     usecase.GroupUsecase
	GroupTransformer transformer.GroupTransformer
	GroupValidator   validator.GroupValidator
	AuthClient       *client.AuthClient
}

func NewGroupHandler(
	groupUsecase usecase.GroupUsecase,
	groupTransformer transformer.GroupTransformer,
	groupValidator validator.GroupValidator,
	authClient *client.AuthClient,
) *GroupHandler {
	return &GroupHandler{
		GroupUsecase:     groupUsecase,
		GroupTransformer: groupTransformer,
		GroupValidator:   groupValidator,
		AuthClient:       authClient,
	}
}

// @Summary Tạo nhóm
// @Description Tạo nhóm
// @Tags Nhóm
// @Accept json
// @Produce json
// @Param group body organizationpb.CreateGroupRequest true "Thông tin nhóm"
// @Security BearerAuth
// @Router /group [post]
func (h *GroupHandler) CreateGroup(ctx context.Context, req *organizationpb.CreateGroupRequest) (*organizationpb.CreateGroupResponse, error) {
	if err := h.GroupValidator.ValidateCreateGroupRequest(req); err != nil {
		return nil, err
	}

	group := h.GroupTransformer.CreateGroupRequestToEntity(req)

	group, err := h.GroupUsecase.CreateGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return h.GroupTransformer.EntityToCreateGroupResponse(group), nil
}

// @Summary Cập nhật nhóm
// @Description Cập nhật nhóm
// @Tags Nhóm
// @Accept json
// @Produce json
// @Param group body organizationpb.UpdateGroupRequest true "Thông tin nhóm"
// @Security BearerAuth
// @Router /group/{id} [put]
func (h *GroupHandler) UpdateGroup(ctx context.Context, req *organizationpb.UpdateGroupRequest) (*organizationpb.UpdateGroupResponse, error) {
	if err := h.GroupValidator.ValidateUpdateGroupRequest(req); err != nil {
		return nil, err
	}

	group := h.GroupTransformer.UpdateGroupRequestToEntity(req)

	group, err := h.GroupUsecase.UpdateGroup(ctx, group)
	if err != nil {
		return nil, err
	}

	return h.GroupTransformer.EntityToUpdateGroupResponse(group), nil
}

// @Summary Xóa nhóm
// @Description Xóa nhóm
// @Tags Nhóm
// @Accept json
// @Produce json
// @Param id path int true "ID của nhóm"
// @Security BearerAuth
// @Router /group/{id} [delete]
func (h *GroupHandler) DeleteGroup(ctx context.Context, req *organizationpb.DeleteGroupRequest) (*organizationpb.DeleteGroupResponse, error) {
	if err := h.GroupValidator.ValidateDeleteGroupRequest(req); err != nil {
		return nil, err
	}
	err := h.GroupUsecase.DeleteGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &organizationpb.DeleteGroupResponse{Id: req.Id}, nil
}

// @Summary Lấy nhóm
// @Description Lấy nhóm
// @Tags Nhóm
// @Accept json
// @Produce json
// @Param id path int true "ID của nhóm"
// @Security BearerAuth
// @Router /group [get]
func (h *GroupHandler) GetGroup(ctx context.Context, req *organizationpb.GetGroupRequest) (*organizationpb.GetGroupResponse, error) {
	if err := h.GroupValidator.ValidateGetGroupRequest(req); err != nil {
		return nil, err
	}

	group, err := h.GroupUsecase.GetGroupByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.GroupTransformer.EntityToGetGroupResponse(group), nil
}

// @Summary Lấy danh sách nhóm của user hiện tại
// @Description Lấy danh sách nhóm của user hiện tại
// @Tags Nhóm
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Security BearerAuth
// @Router /group [get]
func (h *GroupHandler) GetGroups(ctx context.Context, req *organizationpb.GetGroupsRequest) (*organizationpb.GetGroupsResponse, error) {
	if err := h.GroupValidator.ValidateGetGroupsRequest(req); err != nil {
		return nil, err
	}

	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	groups, total, err := h.GroupUsecase.GetGroupsWithDetails(ctx, page, size)
	if err != nil {
		return nil, err
	}

	h.AuthClient.MapRoleToGroups(ctx, groups)

	return h.GroupTransformer.EntityToGetGroupsWithDetailsResponse(groups, uint32(total)), nil
}

// @Summary Lấy danh sách nhóm của user hiện tại
// @Description Lấy danh sách nhóm của user hiện tại
// @Tags Nhóm
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Security BearerAuth
// @Router /group/by-member [get]
func (h *GroupHandler) GetGroupByUserId(ctx context.Context, req *organizationpb.GetGroupByUserIdRequest) (*organizationpb.GetGroupsResponse, error) {
	groups, err := h.GroupUsecase.GetGroupByUserIdWithDetails(ctx)
	if err != nil {
		return nil, err
	}
	return h.GroupTransformer.EntityToGetGroupsWithDetailsResponse(groups, uint32(len(groups))), nil
}
