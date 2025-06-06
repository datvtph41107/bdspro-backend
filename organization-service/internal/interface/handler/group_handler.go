package handler

import (
	"context"
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
)

type GroupHandler struct {
	grouppb.UnimplementedGroupServiceServer
	GroupUsecase     usecase.GroupUsecase
	GroupTransformer transformer.GroupTransformer
	GroupValidator   validator.GroupValidator
}

func NewGroupHandler(
	groupUsecase usecase.GroupUsecase,
	groupTransformer transformer.GroupTransformer,
	groupValidator validator.GroupValidator,
) *GroupHandler {
	return &GroupHandler{
		GroupUsecase:     groupUsecase,
		GroupTransformer: groupTransformer,
		GroupValidator:   groupValidator,
	}
}

func (h *GroupHandler) CreateGroup(ctx context.Context, req *grouppb.CreateGroupRequest) (*grouppb.CreateGroupResponse, error) {
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

func (h *GroupHandler) UpdateGroup(ctx context.Context, req *grouppb.UpdateGroupRequest) (*grouppb.UpdateGroupResponse, error) {
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

func (h *GroupHandler) DeleteGroup(ctx context.Context, req *grouppb.DeleteGroupRequest) (*grouppb.DeleteGroupResponse, error) {
	if err := h.GroupValidator.ValidateDeleteGroupRequest(req); err != nil {
		return nil, err
	}
	err := h.GroupUsecase.DeleteGroup(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &grouppb.DeleteGroupResponse{Id: req.Id}, nil
}

func (h *GroupHandler) GetGroup(ctx context.Context, req *grouppb.GetGroupRequest) (*grouppb.GetGroupResponse, error) {
	if err := h.GroupValidator.ValidateGetGroupRequest(req); err != nil {
		return nil, err
	}

	group, err := h.GroupUsecase.GetGroupByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return h.GroupTransformer.EntityToGetGroupResponse(group), nil
}

func (h *GroupHandler) GetGroups(ctx context.Context, req *grouppb.GetGroupsRequest) (*grouppb.GetGroupsResponse, error) {
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

	groups, total, err := h.GroupUsecase.GetGroups(ctx, page, size)
	if err != nil {
		return nil, err
	}

	return h.GroupTransformer.EntityToGetGroupsResponse(groups, total), nil
}
