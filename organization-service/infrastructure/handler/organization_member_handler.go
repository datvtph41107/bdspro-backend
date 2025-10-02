package handler

import (
	_dto "common/domain/dto"
	"context"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"

	"organization/infrastructure/client"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationMemberHandler struct {
	organizationpb.UnimplementedOrganizationMemberServiceServer
	OrganizationMemberUsecase     usecase.OrganizationMemberUsecase
	OrganizationMemberTransformer transformer.OrganizationMemberTransformer
	OrganizationMemberValidator   validator.OrganizationMemberValidator
	OrganizationUsecase           usecase.OrganizationUsecase
	OrganizationTransformer       transformer.OrganizationTransformer
	internalOrganizationService   *InternalOrganizationHandler
	authClient                    *client.AuthClient
}

func NewOrganizationMemberHandler(organizationMemberUsecase usecase.OrganizationMemberUsecase,
	organizationMemberTransformer transformer.OrganizationMemberTransformer,
	organizationMemberValidator validator.OrganizationMemberValidator,
	organizationUsecase usecase.OrganizationUsecase,
	organizationTransformer transformer.OrganizationTransformer,
	internalOrganizationService *InternalOrganizationHandler,
	authClient *client.AuthClient,
) *OrganizationMemberHandler {
	return &OrganizationMemberHandler{
		OrganizationMemberUsecase:     organizationMemberUsecase,
		OrganizationMemberTransformer: organizationMemberTransformer,
		OrganizationMemberValidator:   organizationMemberValidator,
		OrganizationUsecase:           organizationUsecase,
		OrganizationTransformer:       organizationTransformer,
		internalOrganizationService:   internalOrganizationService,
		authClient:                    authClient,
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
	// permissions := make([]string, len(member.Role.Permissions))
	// for i, permission := range member.Role.Permissions {
	// 	permissions[i] = permission.Key
	// }
	return &organizationpb.SwitchOrganizationResponse{
		IsSuccess: true,
		Id:        req.OrganizationId,
		UserId:    req.UserId,
		// Role:        member.Role.Key,
		// Permissions: permissions,
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

// @Summary Tạo nhiều thành viên tổ chức cùng lúc
// @Description Tạo nhiều thành viên tổ chức cùng lúc
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param organizationMembers body organizationpb.CreateOrganizationMemberBatchRequest true "Danh sách thông tin thành viên"
// @Security BearerAuth
// @Router /organization/member/batch [post]
func (h *OrganizationMemberHandler) CreateOrganizationMemberBatch(ctx context.Context, req *organizationpb.CreateOrganizationMemberBatchRequest) (*organizationpb.CreateOrganizationMemberBatchResponse, error) {
	if err := h.OrganizationMemberValidator.ValidateCreateOrganizationMemberBatchRequest(req); err != nil {
		return nil, err
	}

	createdMembers, err := h.OrganizationMemberUsecase.CreateOrganizationMemberBatch(ctx, req.UserIds)
	if err != nil {
		return nil, err
	}

	return h.OrganizationMemberTransformer.EntitiesToCreateOrganizationMemberBatchResponse(createdMembers), nil
}

// @Summary Cập nhật thành viên tổ chức
// @Description Cập nhật thành viên tổ chức
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param organizationMember body organizationpb.UpdateOrganizationMemberRequest true "Thông tin thành viên"
// @Security BearerAuth
// @Router /organization/member/{id} [put]
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

// @Summary Xóa thành viên tổ chức
// @Description Xóa thành viên tổ chức
// @Tags Thành viên
// @Accept json
// @Produce json
// @Param id path int true "ID của thành viên"
// @Security BearerAuth
// @Router /organization/member/{id} [delete]
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
// @Param roleId query int false "Role ID để lọc thành viên"
// @Param page query int false "Page number"
// @Param size query int false "Page size"
// @Param text query string false "Tên thành viên"
// @Security BearerAuth
// @Router /organization/member [get]
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

	// Get roleId parameter for filtering
	var roleId *uint32
	if req.RoleId != nil {
		roleId = req.RoleId
	}

	// Get members with pagination and role filtering
	members, total, err := h.OrganizationMemberUsecase.GetOrganizationMembers(ctx, page, size, roleId)
	if err != nil {
		return nil, err
	}

	return h.OrganizationMemberTransformer.EntityToPbsResponse(members, total, true), nil
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
func (h *OrganizationMemberHandler) GetCurrentOrganizationMembers(ctx context.Context, req *sharepb.FriendListRequest) (*organizationpb.GetOrganizationMembersResponse, error) {
	dto := _dto.Pagable{
		Page: uint32(req.Page),
		Size: uint32(req.Size),
	}
	members, total, err := h.OrganizationMemberUsecase.GetCurrentOrganizationMember(ctx, dto, req.DealId)
	if err != nil {
		return nil, err
	}

	results := h.OrganizationMemberTransformer.EntityToPbsResponse(members, total, false)
	h.authClient.MapRoleToOrganizationMemberPb(ctx, results.Data)

	// if req.DealId != nil {
	// 	dealMembers, err := h.internalOrganizationService.GetDealMemberByIds(ctx, &organizationpb.GetMemberByIdsRequest{
	// 		Ids: []uint64{*req.DealId},
	// 	})

	// 	// if err != nil {
	// 	// 	return nil, err
	// 	// }
	// 	members = append(members, dealMembers...)
	// }

	return results, nil
}

func (h *OrganizationMemberHandler) GetCurrentOrganizationMemberRole(ctx context.Context, req *organizationpb.Empty) (*organizationpb.GetCurrentOrganizationMemberRoleResponse, error) {
	roleId, joinDate, roleName, err := h.OrganizationMemberUsecase.GetCurrentOrganizationMemberRole(ctx)
	if err != nil {
		return nil, err
	}
	return &organizationpb.GetCurrentOrganizationMemberRoleResponse{RoleId: roleId, JoinDate: joinDate, RoleName: roleName}, nil
}
