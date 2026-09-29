package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"organization/infrastructure/client"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
)

type OrganizationBranchHandler struct {
	organizationpb.UnimplementedOrganizationBranchServiceServer
	usecase     usecase.OrganizationBranchUsecase
	validator   validator.OrganizationBranchValidator
	transformer transformer.OrganizationBranchTransformer
	userClient  *client.UserClient
}

func NewOrganizationBranchHandler(
	usecase usecase.OrganizationBranchUsecase,
	validator validator.OrganizationBranchValidator,
	transformer transformer.OrganizationBranchTransformer,
	userClient *client.UserClient,
) *OrganizationBranchHandler {
	return &OrganizationBranchHandler{
		usecase:     usecase,
		validator:   validator,
		transformer: transformer,
		userClient:  userClient,
	}
}

// @Summary Tạo chi nhánh
// @Description Tạo chi nhánh
// @Tags Chi nhánh
// @Accept json
// @Produce json
// @Param organization_branch body organizationpb.CreateOrganizationBranchRequest true "Thông tin chi nhánh"
// @Security BearerAuth
// @Router /organization/branch [post]
func (h *OrganizationBranchHandler) CreateOrganizationBranch(ctx context.Context, req *organizationpb.CreateOrganizationBranchRequest) (*organizationpb.CreateOrganizationBranchResponse, error) {
	if err := h.validator.ValidateCreateOrganizationBranch(req); err != nil {
		return nil, err
	}
	reqEntity := h.transformer.CreateOrganizationBranchRequestToEntity(req)
	branch, err := h.usecase.CreateOrganizationBranch(ctx, reqEntity)
	if err != nil {
		return nil, err
	}
	return &organizationpb.CreateOrganizationBranchResponse{
		Id: branch.Id,
	}, nil
}

// @Summary Cập nhật chi nhánh
// @Description Cập nhật chi nhánh
// @Tags Chi nhánh
// @Accept json
// @Produce json
// @Param organization_branch body organizationpb.UpdateOrganizationBranchRequest true "Thông tin chi nhánh"
// @Security BearerAuth
// @Router /organization/branch/{id} [put]
func (h *OrganizationBranchHandler) UpdateOrganizationBranch(ctx context.Context, req *organizationpb.UpdateOrganizationBranchRequest) (*organizationpb.UpdateOrganizationBranchResponse, error) {
	if err := h.validator.ValidateUpdateOrganizationBranch(req); err != nil {
		return nil, err
	}
	reqEntity := h.transformer.UpdateOrganizationBranchRequestToEntity(req)
	branch, err := h.usecase.UpdateOrganizationBranch(ctx, reqEntity)
	if err != nil {
		return nil, err
	}
	return &organizationpb.UpdateOrganizationBranchResponse{
		Id: branch.Id,
	}, nil
}

// @Summary Xóa chi nhánh
// @Description Xóa chi nhánh
// @Tags Chi nhánh
// @Accept json
// @Produce json
// @Param organization_branch body organizationpb.DeleteOrganizationBranchRequest true "Thông tin chi nhánh"
// @Security BearerAuth
// @Router /organization/branch/{id} [delete]
func (h *OrganizationBranchHandler) DeleteOrganizationBranch(ctx context.Context, req *organizationpb.DeleteOrganizationBranchRequest) (*organizationpb.DeleteOrganizationBranchResponse, error) {
	if err := h.validator.ValidateDeleteOrganizationBranch(req); err != nil {
		return nil, err
	}
	err := h.usecase.DeleteOrganizationBranch(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &organizationpb.DeleteOrganizationBranchResponse{
		Success: true,
	}, nil
}

// @Summary Lấy chi tiết chi nhánh
// @Description Lấy chi tiết chi nhánh
// @Tags Chi nhánh
// @Accept json
// @Produce json
// @Param id path uint32 true "ID chi nhánh"
// @Security BearerAuth
// @Router /organization/branch/{id} [get]
func (h *OrganizationBranchHandler) GetOrganizationBranchDetail(ctx context.Context, req *organizationpb.GetOrganizationBranchDetailRequest) (*organizationpb.OrganizationBranchDetail, error) {
	if err := h.validator.ValidateGetOrganizationBranchDetail(req); err != nil {
		return nil, err
	}

	branch, err := h.usecase.GetOrganizationBranchDetail(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	result := &organizationpb.OrganizationBranchDetail{
		Id:          branch.Id,
		Name:        branch.Name,
		Address:     branch.Address,
		Phone:       branch.Phone,
		Email:       branch.Email,
		ManagerId:   branch.ManagerId,
		IsActive:    branch.IsActive,
		Type:        string(branch.Type),
		TotalMember: branch.TotalMember,
		TotalAsset:  branch.TotalAsset,
		TotalDeal:   branch.TotalDeal,
		CreatedAt:   branch.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   branch.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:   branch.CreatedBy,
		Description: branch.Description,
	}
	h.userClient.MapToOrganizationBranchDetailPb(ctx, result)
	return result, nil
}

// @Summary Thêm thành viên vào chi nhánh
// @Description Thêm thành viên vào chi nhánh
// @Tags Chi nhánh
// @Accept json
// @Produce json
// @Param id path uint32 true "ID chi nhánh"
// @Param organization_branch body organizationpb.AddMemberToOrganizationBranchRequest true "Thông tin chi nhánh"
// @Security BearerAuth
// @Router /organization/branch/{id}/member [post]
func (h *OrganizationBranchHandler) AddMemberToOrganizationBranch(ctx context.Context, req *organizationpb.AddMemberToOrganizationBranchRequest) (*organizationpb.AddMemberToOrganizationBranchResponse, error) {
	if err := h.validator.ValidateAddMemberToOrganizationBranch(req); err != nil {
		return nil, err
	}
	branch, err := h.usecase.AddMemberToOrganizationBranch(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}
	return &organizationpb.AddMemberToOrganizationBranchResponse{
		Id: branch.Id,
	}, nil
}

// @Summary Xóa thành viên khỏi chi nhánh
// @Description Xóa thành viên khỏi chi nhánh
// @Tags Chi nhánh
// @Accept json
// @Produce json
// @Param organization_branch body organizationpb.RemoveMemberFromOrganizationBranchRequest true "Thông tin chi nhánh"
// @Security BearerAuth
// @Router /organization/branch/{id}/member/{user_id} [delete]
func (h *OrganizationBranchHandler) RemoveMemberFromOrganizationBranch(ctx context.Context, req *organizationpb.RemoveMemberFromOrganizationBranchRequest) (*organizationpb.RemoveMemberFromOrganizationBranchResponse, error) {
	if err := h.validator.ValidateRemoveMemberFromOrganizationBranch(req); err != nil {
		return nil, err
	}

	err := h.usecase.RemoveMemberFromOrganizationBranch(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}
	return &organizationpb.RemoveMemberFromOrganizationBranchResponse{
		Id: req.Id,
	}, nil
}

// @Summary Lấy danh sách thành viên của chi nhánh
// @Description Lấy danh sách thành viên của chi nhánh
// @Tags Chi nhánh
// @Accept json
// @Produce json
// @Param id path uint32 true "ID chi nhánh"
// @Param organization_branch query organizationpb.GetOrganizationBranchMembersRequest true "Thông tin chi nhánh"
// @Security BearerAuth
// @Router /organization/branch/{id}/members [get]
func (h *OrganizationBranchHandler) GetOrganizationBranchMembers(ctx context.Context, req *organizationpb.GetOrganizationBranchMembersRequest) (*organizationpb.GetOrganizationBranchMembersResponse, error) {
	if err := h.validator.ValidateGetOrganizationBranchMembers(req); err != nil {
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

	members, total, err := h.usecase.GetOrganizationBranchMembers(ctx, req.Id, page, size)
	if err != nil {
		return nil, err
	}
	membersResponse := make([]*organizationpb.OrganizationBranchMember, len(members))
	for i, member := range members {
		membersResponse[i] = &organizationpb.OrganizationBranchMember{
			Id:     member.Id,
			UserId: member.UserId,
		}
	}

	h.userClient.MapToOrganizationBranchMemberPb(ctx, membersResponse)
	return &organizationpb.GetOrganizationBranchMembersResponse{
		Data:  membersResponse,
		Total: total,
	}, nil
}

// @Summary Lấy danh sách chi nhánh
// @Description Lấy danh sách chi nhánh
// @Tags Chi nhánh
// @Accept json
// @Produce json
// @Param organization_branch query organizationpb.GetOrganizationBranchesRequest true "Thông tin chi nhánh"
// @Security BearerAuth
// @Router /organization/branch [get]
func (h *OrganizationBranchHandler) GetOrganizationBranches(ctx context.Context, req *organizationpb.GetOrganizationBranchesRequest) (*organizationpb.GetOrganizationBranchesResponse, error) {
	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}
	branches, total, err := h.usecase.GetOrganizationBranches(ctx, page, size)
	if err != nil {
		return nil, err
	}
	branchesResponse := make([]*organizationpb.OrganizationBranch, len(branches))
	for i, branch := range branches {
		branchesResponse[i] = &organizationpb.OrganizationBranch{
			Id:          branch.Id,
			Name:        branch.Name,
			Address:     branch.Address,
			Phone:       branch.Phone,
			Email:       branch.Email,
			IsActive:    branch.IsActive,
			Type:        string(branch.Type),
			TotalMember: branch.TotalMember,
			TotalAsset:  branch.TotalAsset,
			TotalDeal:   branch.TotalDeal,
		}
	}
	return &organizationpb.GetOrganizationBranchesResponse{
		Data:  branchesResponse,
		Total: total,
	}, nil
}

func (h *OrganizationBranchHandler) UpdateIsActive(ctx context.Context, req *organizationpb.UpdateIsActiveRequest) (*organizationpb.UpdateIsActiveResponse, error) {
	err := h.usecase.UpdateIsActive(ctx, req.Id, req.IsActive)
	if err != nil {
		return nil, err
	}
	return &organizationpb.UpdateIsActiveResponse{
		Success: true,
	}, nil
}
