package handler

import (
	"context"

	organizationpb "pb/types/organization"

	"organization/env"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
	"organization/pkg/utils"
)

type GroupMemberHandler struct {
	organizationpb.UnimplementedGroupMemberServiceServer
	groupMemberUsecase     usecase.GroupMemberUsecase
	groupMemberValidator   validator.GroupMemberValidator
	groupMemberTransformer transformer.GroupMemberTransformer
}

func NewGroupMemberHandler(
	groupMemberUsecase usecase.GroupMemberUsecase,
	groupMemberValidator validator.GroupMemberValidator,
	groupMemberTransformer transformer.GroupMemberTransformer,
) *GroupMemberHandler {
	return &GroupMemberHandler{
		groupMemberUsecase:     groupMemberUsecase,
		groupMemberValidator:   groupMemberValidator,
		groupMemberTransformer: groupMemberTransformer,
	}
}

// @Summary Tạo thành viên nhóm
// @Description Tạo thành viên nhóm
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param groupMember body organizationpb.CreateGroupMemberRequest true "Thông tin thành viên"
// @Security BearerAuth
// @Router /group-member [post]
func (h *GroupMemberHandler) CreateGroupMember(ctx context.Context, request *organizationpb.CreateGroupMemberRequest) (*organizationpb.CreateGroupMemberResponse, error) {
	if err := h.groupMemberValidator.ValidateCreateGroupMemberRequest(request); err != nil {
		return nil, err
	}

	member := h.groupMemberTransformer.CreateGroupMemberRequestToEntity(request)

	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	member.CreatedBy = currentUserId
	member.UpdatedBy = currentUserId

	createdMember, err := h.groupMemberUsecase.CreateGroupMember(ctx, member)
	if err != nil {
		return nil, err
	}

	return h.groupMemberTransformer.EntityToCreateGroupMemberResponse(createdMember), nil
}

// @Summary Tạo nhiều thành viên nhóm cùng lúc
// @Description Tạo nhiều thành viên nhóm cùng lúc
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param members body organizationpb.CreateGroupMemberBatchRequest true "Danh sách thành viên nhóm"
// @Security BearerAuth
// @Router /group-member/batch [post]
func (h *GroupMemberHandler) CreateGroupMemberBatch(ctx context.Context, req *organizationpb.CreateGroupMemberBatchRequest) (*organizationpb.CreateGroupMemberBatchResponse, error) {
	if err := h.groupMemberValidator.ValidateCreateGroupMemberBatchRequest(req); err != nil {
		return nil, err
	}

	createdMembers, err := h.groupMemberUsecase.CreateGroupMemberBatch(ctx, req.GroupId, req.UserIds)
	if err != nil {
		return nil, err
	}

	return h.groupMemberTransformer.EntitiesToCreateGroupMemberBatchResponse(createdMembers), nil
}

// @Summary Cập nhật thành viên nhóm
// @Description Cập nhật thành viên nhóm
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param groupMember body organizationpb.UpdateGroupMemberRequest true "Thông tin thành viên"
// @Security BearerAuth
// @Router /group-member/{id} [put]
func (h *GroupMemberHandler) UpdateGroupMember(ctx context.Context, request *organizationpb.UpdateGroupMemberRequest) (*organizationpb.UpdateGroupMemberResponse, error) {
	if err := h.groupMemberValidator.ValidateUpdateGroupMemberRequest(request); err != nil {
		return nil, err
	}

	member := h.groupMemberTransformer.UpdateGroupMemberRequestToEntity(request)
	currentUserId := utils.GetUserID(ctx, env.USER_CONTEXT)
	member.UpdatedBy = currentUserId
	member.CreatedBy = currentUserId

	updatedMember, err := h.groupMemberUsecase.UpdateGroupMember(ctx, member)
	if err != nil {
		return nil, err
	}

	return h.groupMemberTransformer.EntityToUpdateGroupMemberResponse(updatedMember), nil
}

// @Summary Xóa thành viên nhóm
// @Description Xóa thành viên nhóm
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param id path int true "ID của thành viên"
// @Security BearerAuth
// @Router /group-member/{id} [delete]
func (h *GroupMemberHandler) DeleteGroupMember(ctx context.Context, request *organizationpb.DeleteGroupMemberRequest) (*organizationpb.DeleteGroupMemberResponse, error) {
	if err := h.groupMemberValidator.ValidateDeleteGroupMemberRequest(request); err != nil {
		return nil, err
	}

	if err := h.groupMemberUsecase.DeleteGroupMember(ctx, request.Id); err != nil {
		return nil, err
	}

	return &organizationpb.DeleteGroupMemberResponse{
		Id: request.Id,
	}, nil
}

// @Summary Lấy thành viên nhóm
// @Description Lấy thành viên nhóm
// @Tags Thành viên
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param groupMemberId path int true "ID của thành viên"
// @Router /group-member/{groupMemberId} [get]
func (h *GroupMemberHandler) GetGroupMember(ctx context.Context, request *organizationpb.GetGroupMemberRequest) (*organizationpb.GetGroupMemberResponse, error) {
	if err := h.groupMemberValidator.ValidateGetGroupMemberRequest(request); err != nil {
		return nil, err
	}

	member, err := h.groupMemberUsecase.GetGroupMemberByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return h.groupMemberTransformer.EntityToGetGroupMemberResponse(member), nil
}

// @Summary Lấy danh sách thành viên nhóm
// @Description Lấy danh sách thành viên nhóm
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param groupId path int true "ID của nhóm"
// @Security BearerAuth
// @Router /group-member/list/{groupId} [get]
func (h *GroupMemberHandler) GetGroupMembers(ctx context.Context, request *organizationpb.GetGroupMembersRequest) (*organizationpb.GetGroupMembersResponse, error) {
	if err := h.groupMemberValidator.ValidateGetGroupMembersRequest(request); err != nil {
		return nil, err
	}
	members, err := h.groupMemberUsecase.GetGroupMembersByGroupID(ctx, request.GroupId)
	if err != nil {
		return nil, err
	}
	return h.groupMemberTransformer.EntityToGetGroupMembersResponse(members), nil
}

// @Summary Đổi nhóm
// @Description Đổi nhóm
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param groupId path int true "ID của nhóm"
// @Security BearerAuth
// @Router /group-member/switch [post]
func (h *GroupMemberHandler) SwitchGroup(ctx context.Context, request *organizationpb.SwitchGroupRequest) (*organizationpb.SwitchGroupResponse, error) {
	// todo: thêm logic đổi nhóm check xem có thuộc nhóm không

	return &organizationpb.SwitchGroupResponse{
		IsSuccess: true,
	}, nil
}
