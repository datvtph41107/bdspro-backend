package handler

import (
	"context"

	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
	"github.com/tuandq2112/organization-service/pkg/utils"
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

	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	member.CreatedBy = currentUserId
	member.UpdatedBy = currentUserId

	createdMember, err := h.groupMemberUsecase.CreateGroupMember(ctx, member)
	if err != nil {
		return nil, err
	}

	return h.groupMemberTransformer.EntityToCreateGroupMemberResponse(createdMember), nil
}

func (h *GroupMemberHandler) UpdateGroupMember(ctx context.Context, request *organizationpb.UpdateGroupMemberRequest) (*organizationpb.UpdateGroupMemberResponse, error) {
	if err := h.groupMemberValidator.ValidateUpdateGroupMemberRequest(request); err != nil {
		return nil, err
	}

	member := h.groupMemberTransformer.UpdateGroupMemberRequestToEntity(request)
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	member.UpdatedBy = currentUserId
	member.CreatedBy = currentUserId

	updatedMember, err := h.groupMemberUsecase.UpdateGroupMember(ctx, member)
	if err != nil {
		return nil, err
	}

	return h.groupMemberTransformer.EntityToUpdateGroupMemberResponse(updatedMember), nil
}

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

func (h *GroupMemberHandler) SwitchGroup(ctx context.Context, request *organizationpb.SwitchGroupRequest) (*organizationpb.SwitchGroupResponse, error) {
	// todo: thêm logic đổi nhóm check xem có thuộc nhóm không

	return &organizationpb.SwitchGroupResponse{
		IsSuccess: true,
	}, nil
}
