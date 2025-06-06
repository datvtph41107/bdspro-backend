package handler

import (
	"context"

	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
	"github.com/tuandq2112/organization-service/pkg/utils"
)

type GroupMemberHandler struct {
	grouppb.UnimplementedGroupMemberServiceServer
	groupMemberUsecase    usecase.GroupMemberUsecase
	groupMemberValidator  validator.GroupMemberValidator
	groupMemberTransformer transformer.GroupMemberTransformer
}

func NewGroupMemberHandler(
	groupMemberUsecase usecase.GroupMemberUsecase,
	groupMemberValidator validator.GroupMemberValidator,
	groupMemberTransformer transformer.GroupMemberTransformer,
) *GroupMemberHandler {
	return &GroupMemberHandler{
		groupMemberUsecase:    groupMemberUsecase,
		groupMemberValidator:  groupMemberValidator,
		groupMemberTransformer: groupMemberTransformer,
	}
}

func (h *GroupMemberHandler) CreateGroupMember(ctx context.Context, request *grouppb.CreateGroupMemberRequest) (*grouppb.CreateGroupMemberResponse, error) {
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

func (h *GroupMemberHandler) UpdateGroupMember(ctx context.Context, request *grouppb.UpdateGroupMemberRequest) (*grouppb.UpdateGroupMemberResponse, error) {
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

func (h *GroupMemberHandler) DeleteGroupMember(ctx context.Context, request *grouppb.DeleteGroupMemberRequest) (*grouppb.DeleteGroupMemberResponse, error) {
	if err := h.groupMemberValidator.ValidateDeleteGroupMemberRequest(request); err != nil {
		return nil, err
	}

	if err := h.groupMemberUsecase.DeleteGroupMember(ctx, request.Id); err != nil {
		return nil, err
	}

	return &grouppb.DeleteGroupMemberResponse{
		Id: request.Id,
	}, nil
}

func (h *GroupMemberHandler) GetGroupMember(ctx context.Context, request *grouppb.GetGroupMemberRequest) (*grouppb.GetGroupMemberResponse, error) {
	if err := h.groupMemberValidator.ValidateGetGroupMemberRequest(request); err != nil {
		return nil, err
	}

	member, err := h.groupMemberUsecase.GetGroupMemberByID(ctx, request.Id)
	if err != nil {
		return nil, err
	}

	return h.groupMemberTransformer.EntityToGetGroupMemberResponse(member), nil
}

func (h *GroupMemberHandler) GetGroupMembers(ctx context.Context, request *grouppb.GetGroupMembersRequest) (*grouppb.GetGroupMembersResponse, error) {
	if err := h.groupMemberValidator.ValidateGetGroupMembersRequest(request); err != nil {
		return nil, err
	}

	members, err := h.groupMemberUsecase.GetGroupMembersByGroupID(ctx, request.GroupId)
	if err != nil {
		return nil, err
	}

	return h.groupMemberTransformer.EntityToGetGroupMembersResponse(members), nil
} 