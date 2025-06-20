package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationMemberHandler struct {
	organizationpb.UnimplementedOrganizationMemberServiceServer
	OrganizationMemberUsecase     usecase.OrganizationMemberUsecase
	OrganizationMemberTransformer transformer.OrganizationMemberTransformer
	OrganizationMemberValidator   validator.OrganizationMemberValidator
}

func NewOrganizationMemberHandler(organizationMemberUsecase usecase.OrganizationMemberUsecase,
	organizationMemberTransformer transformer.OrganizationMemberTransformer,
	organizationMemberValidator validator.OrganizationMemberValidator,
) *OrganizationMemberHandler {
	return &OrganizationMemberHandler{
		OrganizationMemberUsecase:     organizationMemberUsecase,
		OrganizationMemberTransformer: organizationMemberTransformer,
		OrganizationMemberValidator:   organizationMemberValidator,
	}
}

func (h *OrganizationMemberHandler) SwitchOrganization(ctx context.Context, req *organizationpb.SwitchOrganizationRequest) (*organizationpb.SwitchOrganizationResponse, error) {
	member, err := h.OrganizationMemberUsecase.FindByUserIdAndOrganizationIdWithRole(ctx, req.UserId, req.OrganizationId)
	if err != nil {
		return nil, err
	}

	if member == nil {
		return nil, status.Errorf(codes.NotFound, "user not in organization")
	}
	permissions := make([]string, len(member.Role.Permissions))
	for i, permission := range member.Role.Permissions {
		permissions[i] = permission.Key
	}
	return &organizationpb.SwitchOrganizationResponse{
		IsSuccess:   true,
		Id:          req.OrganizationId,
		UserId:      req.UserId,
		Role:        member.Role.Key,
		Permissions: permissions,
	}, nil
}

// @Summary Tạo thành viên tổ chức
// @Description Tạo thành viên tổ chức
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param organizationMember body organizationpb.CreateOrganizationMemberRequest true "Thông tin thành viên"
// @Security BearerAuth
// @Router /organization/member [post]
func (h *OrganizationMemberHandler) CreateOrganizationMember(ctx context.Context, req *organizationpb.CreateOrganizationMemberRequest) (*organizationpb.CreateOrganizationMemberResponse, error) {
	organizationMember := h.OrganizationMemberTransformer.CreateOrganizationMemberRequestToEntity(req)
	if err := h.OrganizationMemberValidator.ValidateCreateOrganizationMemberRequest(req); err != nil {
		return nil, err
	}
	organizationMember, err := h.OrganizationMemberUsecase.CreateOrganizationMember(ctx, organizationMember)
	if err != nil {
		return nil, err
	}
	return h.OrganizationMemberTransformer.EntityToCreateOrganizationMemberResponse(organizationMember), nil
}

func (h *OrganizationMemberHandler) UpdateOrganizationMember(ctx context.Context, req *organizationpb.UpdateOrganizationMemberRequest) (*organizationpb.UpdateOrganizationMemberResponse, error) {
	organizationMember := h.OrganizationMemberTransformer.UpdateOrganizationMemberRequestToEntity(req)
	if err := h.OrganizationMemberValidator.ValidateUpdateOrganizationMemberRequest(req); err != nil {
		return nil, err
	}
	organizationMember, err := h.OrganizationMemberUsecase.UpdateOrganizationMember(ctx, organizationMember)
	if err != nil {
		return nil, err
	}
	return h.OrganizationMemberTransformer.EntityToUpdateOrganizationMemberResponse(organizationMember), nil
}

func (h *OrganizationMemberHandler) DeleteOrganizationMember(ctx context.Context, req *organizationpb.DeleteOrganizationMemberRequest) (*organizationpb.DeleteOrganizationMemberResponse, error) {
	organizationMember, err := h.OrganizationMemberUsecase.DeleteOrganizationMember(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.OrganizationMemberTransformer.EntityToDeleteOrganizationMemberResponse(organizationMember), nil
}

// @Summary Lấy danh sách thành viên của tổ chức hiện tại mà người dùng đang đăng nhập
// @Description Lấy danh sách thành viên của tổ chức hiện tại mà người dùng đang đăng nhập
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param text query string false "Tên thành viên"
// @Param organizationId path int true "ID của tổ chức"
// @Security BearerAuth
// @Router /organization/member/{organizationId} [get]
func (h *OrganizationMemberHandler) GetOrganizationMembers(ctx context.Context, req *organizationpb.GetOrganizationMembersRequest) (*organizationpb.GetOrganizationMembersResponse, error) {
	// Get pagination parameters
	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	// Get members with pagination
	members, total, err := h.OrganizationMemberUsecase.GetOrganizationMembers(ctx, req.OrganizationId, page, size)
	if err != nil {
		return nil, err
	}

	return h.OrganizationMemberTransformer.EntityToGetOrganizationMembersResponse(members, total), nil
}

// @Summary Lấy danh sách thành viên của tổ chức hiện tại mà người dùng đang đăng nhập
// @Description Lấy danh sách thành viên của tổ chức hiện tại mà người dùng đang đăng nhập
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Security BearerAuth
// @Router /organization/members/me [get]
func (h *OrganizationMemberHandler) GetCurrentOrganizationMembers(ctx context.Context, req *organizationpb.Empty) (*organizationpb.GetOrganizationMembersResponse, error) {
	members, total, err := h.OrganizationMemberUsecase.GetCurrentOrganizationMember(ctx, 0, 50)
	if err != nil {
		return nil, err
	}
	return h.OrganizationMemberTransformer.EntityToGetOrganizationMembersResponse(members, total), nil
}
