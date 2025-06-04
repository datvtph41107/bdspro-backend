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

type OrganizationBranchHandler struct {
	organizationpb.UnimplementedOrganizationBranchServiceServer
	usecase     usecase.OrganizationBranchUsecase
	validator   validator.OrganizationBranchValidator
	transformer transformer.OrganizationBranchTransformer
}

func NewOrganizationBranchHandler(
	usecase usecase.OrganizationBranchUsecase,
	validator validator.OrganizationBranchValidator,
	transformer transformer.OrganizationBranchTransformer,
) *OrganizationBranchHandler {
	return &OrganizationBranchHandler{
		usecase:     usecase,
		validator:   validator,
		transformer: transformer,
	}
}

func (h *OrganizationBranchHandler) CreateOrganizationBranch(ctx context.Context, req *organizationpb.CreateOrganizationBranchRequest) (*organizationpb.CreateOrganizationBranchResponse, error) {
	if err := h.validator.ValidateCreateOrganizationBranch(req); err != nil {
		return nil, err
	}
	reqEntity := h.transformer.CreateOrganizationBranchRequestToEntity(req)
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	branch, err := h.usecase.CreateOrganizationBranch(ctx, currentUserId, reqEntity)
	if err != nil {
		return nil, err
	}
	return &organizationpb.CreateOrganizationBranchResponse{
		Id: branch.Id,
	}, nil
}

func (h *OrganizationBranchHandler) UpdateOrganizationBranch(ctx context.Context, req *organizationpb.UpdateOrganizationBranchRequest) (*organizationpb.UpdateOrganizationBranchResponse, error) {
	if err := h.validator.ValidateUpdateOrganizationBranch(req); err != nil {
		return nil, err
	}
	reqEntity := h.transformer.UpdateOrganizationBranchRequestToEntity(req)
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	branch, err := h.usecase.UpdateOrganizationBranch(ctx, currentUserId, reqEntity)
	if err != nil {
		return nil, err
	}
	return &organizationpb.UpdateOrganizationBranchResponse{
		Id: branch.Id,
	}, nil
}

func (h *OrganizationBranchHandler) DeleteOrganizationBranch(ctx context.Context, req *organizationpb.DeleteOrganizationBranchRequest) (*organizationpb.DeleteOrganizationBranchResponse, error) {
	if err := h.validator.ValidateDeleteOrganizationBranch(req); err != nil {
		return nil, err
	}
	err := h.usecase.DeleteOrganizationBranch(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), req.Id)
	if err != nil {
		return nil, err
	}
	return &organizationpb.DeleteOrganizationBranchResponse{
		Success: true,
	}, nil
}

func (h *OrganizationBranchHandler) AddMemberToOrganizationBranch(ctx context.Context, req *organizationpb.AddMemberToOrganizationBranchRequest) (*organizationpb.AddMemberToOrganizationBranchResponse, error) {
	if err := h.validator.ValidateAddMemberToOrganizationBranch(req); err != nil {
		return nil, err
	}
	branch, err := h.usecase.AddMemberToOrganizationBranch(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), req.Id, req.UserId)
	if err != nil {
		return nil, err
	}
	return &organizationpb.AddMemberToOrganizationBranchResponse{
		Id: branch.Id,
	}, nil
}

func (h *OrganizationBranchHandler) RemoveMemberFromOrganizationBranch(ctx context.Context, req *organizationpb.RemoveMemberFromOrganizationBranchRequest) (*organizationpb.RemoveMemberFromOrganizationBranchResponse, error) {
	if err := h.validator.ValidateRemoveMemberFromOrganizationBranch(req); err != nil {
		return nil, err
	}
	
	err := h.usecase.RemoveMemberFromOrganizationBranch(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), req.Id, req.UserId)
	if err != nil {
		return nil, err
	}
	return &organizationpb.RemoveMemberFromOrganizationBranchResponse{
		Id: req.Id,
	}, nil
}

func (h *OrganizationBranchHandler) GetOrganizationBranchMembers(ctx context.Context, req *organizationpb.GetOrganizationBranchMembersRequest) (*organizationpb.GetOrganizationBranchMembersResponse, error) {
	if err := h.validator.ValidateGetOrganizationBranchMembers(req); err != nil {
		return nil, err
	}
	members, err := h.usecase.GetOrganizationBranchMembers(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), req.Id)
	if err != nil {
		return nil, err
	}
	membersResponse := make([]*organizationpb.OrganizationBranchMember, len(members))
	for i, member := range members {
		membersResponse[i] = &organizationpb.OrganizationBranchMember{
			Id: member.Id,
			UserId: member.UserId,
		}
	}
	return &organizationpb.GetOrganizationBranchMembersResponse{
		Members: membersResponse,
	}, nil
}

func (h *OrganizationBranchHandler) GetOrganizationBranches(ctx context.Context, req *organizationpb.GetOrganizationBranchesRequest) (*organizationpb.GetOrganizationBranchesResponse, error) {
	
	branches, err := h.usecase.GetOrganizationBranches(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), req.OrganizationId)
	if err != nil {
		return nil, err
	}
	branchesResponse := make([]*organizationpb.OrganizationBranch, len(branches))
	for i, branch := range branches {
		branchesResponse[i] = &organizationpb.OrganizationBranch{
			Id: branch.Id,
			Name: branch.Name,
			Address: branch.Address,
			Phone: branch.Phone,
			Email: branch.Email,
		}
	}
	return &organizationpb.GetOrganizationBranchesResponse{
		Branches: branchesResponse,
	}, nil
}