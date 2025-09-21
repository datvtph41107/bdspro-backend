package handler

import (
	_utils "common/utils"
	"context"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
	"time"

	"organization/infrastructure/client"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/domain/entity"
	"organization/internal/usecase"
)

type OrganizationHandler struct {
	organizationpb.UnimplementedOrganizationServiceServer
	OrganizationUsecase     usecase.OrganizationUsecase
	OrganizationTransformer transformer.OrganizationTransformer
	OrganizationValidator   validator.OrganizationValidator
	UserClient              *client.UserClient
}

func NewOrganizationHandler(organizationUsecase usecase.OrganizationUsecase,
	organizationTransformer transformer.OrganizationTransformer,
	organizationValidator validator.OrganizationValidator,
	userClient *client.UserClient,
) *OrganizationHandler {
	return &OrganizationHandler{
		OrganizationUsecase:     organizationUsecase,
		OrganizationTransformer: organizationTransformer,
		OrganizationValidator:   organizationValidator,
		UserClient:              userClient,
	}
}

// @Summary Tạo tổ chức mới
// @Description Tạo tổ chức mới
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Param request body organizationpb.CreateOrganizationRequest true "Thông tin tổ chức"
// @Security BearerAuth
// @Router /organization/new [post]
func (h *OrganizationHandler) CreateOrganization(ctx context.Context, req *organizationpb.CreateOrganizationRequest) (*organizationpb.CreateOrganizationResponse, error) {
	organization := h.OrganizationTransformer.CreateOrganizationRequestToEntity(req)
	profileId := _utils.GetProfileIdWithContext(ctx)
	organization.CreatedBy = &profileId
	organization.OwnerId = profileId

	member := &entity.OrganizationMember{
		UserID:   profileId,
		Status:   entity.OrganizationMemberStatusActive,
		JoinedAt: time.Now(),
	}
	if err := h.OrganizationValidator.ValidateCreateOrganizationRequest(req); err != nil {
		return nil, err
	}
	organization, _, err := h.OrganizationUsecase.CreateOrganization(ctx, organization, member)
	if err != nil {
		return nil, err
	}
	return h.OrganizationTransformer.EntityToCreateOrganizationResponse(organization), nil
}

// @Summary Cập nhật tổ chức
// @Description Cập nhật tổ chức
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Param request body organizationpb.UpdateOrganizationRequest true "Thông tin tổ chức"
// @Security BearerAuth
// @Router /organization/{id} [put]
func (h *OrganizationHandler) UpdateOrganization(ctx context.Context, req *organizationpb.UpdateOrganizationRequest) (*organizationpb.UpdateOrganizationResponse, error) {
	if err := h.OrganizationValidator.ValidateUpdateOrganizationRequest(req); err != nil {
		return nil, err
	}

	organization := h.OrganizationTransformer.UpdateOrganizationRequestToEntity(req)

	organization, err := h.OrganizationUsecase.UpdateOrganization(ctx, organization)
	if err != nil {
		return nil, err
	}
	return h.OrganizationTransformer.EntityToUpdateOrganizationResponse(organization), nil
}

// @Summary Lấy danh sách tổ chức của người dùng hiện tại
// @Description Lấy danh sách tổ chức của người dùng hiện tại
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Param request body organizationpb.GetOrganizationsRequest true "Thông tin tổ chức"
// @Security BearerAuth
// @Router /organization/by-member [get]
func (h *OrganizationHandler) GetOrganizations(ctx context.Context, req *organizationpb.GetOrganizationsRequest) (*organizationpb.GetOrganizationsResponse, error) {
	// Get pagination parameters
	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	// Calculate offset
	offset := page * size

	// Get organizations with pagination
	organizations, total, err := h.OrganizationUsecase.GetOrganizations(ctx, offset, size)
	if err != nil {
		return nil, err
	}

	organizationPbs := h.OrganizationTransformer.EntityToGetOrganizationsResponse(organizations, total)
	h.UserClient.MapProfileToOrganizationPb(ctx, organizationPbs.Data)

	return organizationPbs, nil
}

// @Summary Lấy tổ chức của người dùng hiện tại
// @Description Lấy tổ chức của người dùng hiện tại
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /organization/by-member [get]
func (h *OrganizationHandler) GetOrganizationByUserId(ctx context.Context, req *organizationpb.GetOrganizationByUserIdRequest) (*organizationpb.GetOrganizationByUserIdResponse, error) {
	organizations, err := h.OrganizationUsecase.GetOrganizationByUserId(ctx)
	if err != nil {
		return nil, err
	}
	return h.OrganizationTransformer.EntitiesToGetOrganizationByUserIdResponse(organizations), nil
}

// @Summary Kiểm tra tổ chức có tồn tại không
// @Description Kiểm tra tổ chức có tồn tại không
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Param request body organizationpb.CheckOrganizationIsExistRequest true "Thông tin tổ chức"
// @Security BearerAuth
// @Router /organization/check/{id} [get]
func (h *OrganizationHandler) CheckOrganizationIsExist(ctx context.Context, req *organizationpb.CheckOrganizationIsExistRequest) (*organizationpb.CheckOrganizationIsExistResponse, error) {
	isExist, err := h.OrganizationUsecase.CheckOrganizationIsExist(ctx)
	if err != nil {
		return nil, err
	}
	return &organizationpb.CheckOrganizationIsExistResponse{IsExist: isExist}, nil
}

// @Summary Lấy tổ chức hiện tại
// @Description Lấy tổ chức hiện tại
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /organization/current [get]
func (h *OrganizationHandler) GetOwnedOrganization(ctx context.Context, req *organizationpb.GetOwnedOrganizationRequest) (*organizationpb.GetOwnedOrganizationResponse, error) {
	organization, err := h.OrganizationUsecase.GetOwnedOrganization(ctx)
	if err != nil {
		return nil, err
	}
	return h.OrganizationTransformer.EntityToGetOwnedOrganizationResponse(organization), nil
}

// @Summary Lấy danh sách tất cả doanh nghiệp với thành viên (Admin only)
// @Description Lấy danh sách tất cả doanh nghiệp với thành viên cho admin
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Param request body organizationpb.GetAllOrganizationsWithMembersForAdminRequest true "Thông tin request"
// @Security BearerAuth
// @Router /organization/admin/all-with-members [get]
func (h *OrganizationHandler) GetAllOrganizationsWithMembersForAdmin(ctx context.Context, req *organizationpb.GetAllOrganizationsWithMembersForAdminRequest) (*organizationpb.GetAllOrganizationsWithMembersForAdminResponse, error) {
	// Get pagination parameters
	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	// Calculate offset
	offset := page * size

	// Get organizations with members for admin
	organizationsWithMembers, total, err := h.OrganizationUsecase.GetAllOrganizationsWithMembersForAdmin(ctx, offset, size)
	if err != nil {
		return nil, err
	}

	result, err := h.OrganizationTransformer.EntitiesToGetAllOrganizationsWithMembersForAdminResponse(organizationsWithMembers, total), nil
	if err != nil {
		return nil, err
	}

	h.UserClient.MapProfileToOrganizationPbWithMembers(ctx, result.Data)

	return result, nil
}

// @Summary Lấy dashboard của tổ chức hiện tại
// @Description Lấy dashboard của tổ chức hiện tại
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /organization/current/dashboard [get]
func (h *OrganizationHandler) GetCurrentDashboard(ctx context.Context, req *sharepb.IdRequest) (*organizationpb.DashboardResponse, error) {
	dashboard, err := h.OrganizationUsecase.GetCurrentDashboard(ctx)
	if err != nil {
		return nil, err
	}
	return dashboard, nil
}

// @Summary Lấy profile của tổ chức global
// @Description Lấy profile của tổ chức global
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Param id path uint64 true "ID tổ chức"
// @Router /organization/global/profile/{id} [get]
func (h *OrganizationHandler) GetGlobalProfile(ctx context.Context, req *sharepb.IdRequest) (*organizationpb.Organization, error) {
	organization, err := h.OrganizationUsecase.GetProfilePublic(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.OrganizationTransformer.EntityToPb(organization), nil
}

// @Summary Lấy danh sách thành viên của tổ chức (Admin only)
// @Description Lấy danh sách thành viên của tổ chức cho admin
// @Tags Tổ chức
// @Accept json
// @Produce json
// @Param page query uint32 false "Số trang"
// @Param size query uint32 false "Số lượng mỗi trang"
// @Security BearerAuth
// @Router /organization/member [get]
func (h *OrganizationHandler) GetOrganizationMembers(ctx context.Context, req *organizationpb.GetOrganizationMembersRequest) (*organizationpb.GetOrganizationMembersResponse, error) {
	// Get pagination parameters
	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	// Calculate offset
	offset := page * size

	// Get organizationId from context
	organizationId := _utils.GetOrganizationIdFromContext(ctx)

	// Get organization members
	membersWithProfiles, total, err := h.OrganizationUsecase.GetOrganizationMembers(ctx, uint32(organizationId), offset, size)
	if err != nil {
		return nil, err
	}

	return h.OrganizationTransformer.EntitiesToGetOrganizationMembersResponse(membersWithProfiles, total), nil
}

