package handler

import (
	"common/identity"
	"context"
	"fmt"
	"math"
	"organization/env"
	"organization/infrastructure/transformer"
	"organization/internal/domain/repository"
	"organization/internal/usecase"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InternalOrganizationHandler struct {
	organizationpb.UnimplementedInternalOrganizationServiceServer
	organizationRepo repository.OrganizationRepository
	groupRepo        repository.GroupRepository
	mapper           *transformer.InternalOrganizationMapper
	// dealInvitationRepo     repository.DealMemberRepository
	groupMemberRepo           repository.GroupMemberRepository
	organizationMemberRepo    repository.OrganizationMemberRepository
	organizationMemberUsecase usecase.OrganizationMemberUsecase
	organizationAuthUsecase   usecase.OrganizationAuthUsecase
	dealRepo                  repository.DealRepository
	dealMemberRepo            repository.DealMemberRepository
	dealUsecase               usecase.DealUsecase
	branchMemberRepo          repository.OrganizationBranchMemberRepository
}

func NewInternalOrganizationHandler(
	organizationRepo repository.OrganizationRepository,
	groupRepo repository.GroupRepository,
	// dealInvitationRepo repository.DealMemberRepository,
	groupMemberRepo repository.GroupMemberRepository,
	organizationMemberRepo repository.OrganizationMemberRepository,
	organizationMemberUsecase usecase.OrganizationMemberUsecase,
	organizationAuthUsecase usecase.OrganizationAuthUsecase,
	dealRepo repository.DealRepository,
	dealMemberRepo repository.DealMemberRepository,
	dealUsecase usecase.DealUsecase,
	branchMemberRepo repository.OrganizationBranchMemberRepository,
	mapper *transformer.InternalOrganizationMapper) *InternalOrganizationHandler {
	return &InternalOrganizationHandler{
		organizationRepo: organizationRepo,
		groupRepo:        groupRepo,
		// dealInvitationRepo:     dealInvitationRepo,
		mapper:                    mapper,
		groupMemberRepo:           groupMemberRepo,
		organizationMemberRepo:    organizationMemberRepo,
		dealRepo:                  dealRepo,
		dealMemberRepo:            dealMemberRepo,
		dealUsecase:               dealUsecase,
		branchMemberRepo:          branchMemberRepo,
		organizationMemberUsecase: organizationMemberUsecase,
		organizationAuthUsecase:   organizationAuthUsecase,
	}
}

// AuthorizeOrganizationAction là transport internal duy nhất công bố quyết
// định permission thuộc organization. Caller chỉ gửi danh tính profile đã
// được xác thực; organization-service tự kiểm membership, role và permission.
func (h *InternalOrganizationHandler) AuthorizeOrganizationAction(ctx context.Context, req *organizationpb.AuthorizeOrganizationActionRequest) (*organizationpb.AuthorizeOrganizationActionResponse, error) {
	if req == nil || req.GetActorProfileId() == 0 || req.GetOrganizationId() == 0 || req.GetOrganizationId() > math.MaxUint32 {
		return nil, status.Error(codes.InvalidArgument, "valid actor profile and organization are required")
	}
	caller, ok := identity.ServiceCallerFromContext(ctx)
	if !ok || caller.ServiceID != env.USER_SERVICE_ID {
		return nil, status.Error(codes.PermissionDenied, "organization action caller is not allowed")
	}
	actor, ok := identity.ActorFromContext(ctx)
	if !ok || actor.ProfileID == 0 || actor.ProfileID != req.GetActorProfileId() {
		return nil, status.Error(codes.Unauthenticated, "signed actor profile does not match request")
	}

	var permission string
	switch req.GetAction() {
	case organizationpb.OrganizationAction_ORGANIZATION_ACTION_SUBSCRIPTION_CHECKOUT:
		permission = env.CHECKOUT_ORGANIZATION_SUBSCRIPTION_PERMISSION
	default:
		return nil, status.Error(codes.InvalidArgument, "unsupported organization action")
	}

	allowed, err := h.organizationAuthUsecase.Authorize(ctx, req.GetActorProfileId(), uint32(req.GetOrganizationId()), permission)
	if err != nil {
		return nil, status.Error(codes.Internal, "organization authorization failed")
	}
	return &organizationpb.AuthorizeOrganizationActionResponse{Allowed: allowed}, nil
}

func (h *InternalOrganizationHandler) GetOrganizationByIds(ctx context.Context, req *organizationpb.GetOrganizationByIdsRequest) (*organizationpb.GetOrganizationByIdsResponse, error) {
	organizations, err := h.organizationRepo.GetByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	pbOrganizations := h.mapper.OrganizationsToPb(organizations)
	return &organizationpb.GetOrganizationByIdsResponse{
		Organizations: pbOrganizations,
	}, nil
}

func (h *InternalOrganizationHandler) GetGroupByIds(ctx context.Context, req *organizationpb.GetGroupByIdsRequest) (*organizationpb.GetGroupByIdsResponse, error) {
	groups, err := h.groupRepo.GetByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	pbGroups := h.mapper.GroupsToPb(groups)
	return &organizationpb.GetGroupByIdsResponse{
		Groups: pbGroups,
	}, nil
}

func (h *InternalOrganizationHandler) GetDealMemberByIds(ctx context.Context, req *organizationpb.GetMemberByIdsRequest) (*organizationpb.GetDealMemberByIdsResponse, error) {
	dealMembers, err := h.dealMemberRepo.GetByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	pbDealMembers := h.mapper.DealMembersToPb(dealMembers)
	return &organizationpb.GetDealMemberByIdsResponse{
		Data: pbDealMembers,
	}, nil
}

func (h *InternalOrganizationHandler) GetGroupMemberByIds(ctx context.Context, req *organizationpb.GetMemberByIdsRequest) (*organizationpb.GetGroupMemberByIdsResponse, error) {
	groupMembers, err := h.groupMemberRepo.GetByGroupIDAndUserIDs(ctx, uint32(req.GroupId), req.Ids)
	if err != nil {
		return nil, err
	}
	pbGroupMembers := h.mapper.GroupMembersToPb(groupMembers)
	return &organizationpb.GetGroupMemberByIdsResponse{
		Data: pbGroupMembers,
	}, nil
}

func (h *InternalOrganizationHandler) GetOrganizationMemberByIds(ctx context.Context, req *organizationpb.GetMemberByIdsRequest) (*organizationpb.GetOrganizationMemberByIdsResponse, error) {
	organizationMembers, err := h.organizationMemberRepo.FindByUserIdsAndOrganizationId(ctx, req.Ids, uint32(req.GroupId))
	if err != nil {
		return nil, err
	}
	pbOrganizationMembers := h.mapper.OrganizationMembersToPb(organizationMembers)
	return &organizationpb.GetOrganizationMemberByIdsResponse{
		Data: pbOrganizationMembers,
	}, nil
}

func (h *InternalOrganizationHandler) GetDealMembers(ctx context.Context, req *organizationpb.DealMemberRequest) (*organizationpb.GetDealMemberByIdsResponse, error) {
	dealMembers, err := h.dealMemberRepo.GetDealMembersByUserIds(ctx, req.DealId, req.UserIds)
	if err != nil {
		return nil, err
	}

	data := h.mapper.DealMembersToPb(dealMembers)
	return &organizationpb.GetDealMemberByIdsResponse{
		Data: data,
	}, nil
}

func (h *InternalOrganizationHandler) GetBranchMembers(ctx context.Context, req *organizationpb.GetBranchMembersRequest) (*organizationpb.GetBranchMemberByIdsResponse, error) {
	branchMembers, err := h.branchMemberRepo.GetOrganizationBranchMemberByUserIdsAndOrganizationBranchId(ctx, req.UserIds, uint32(req.BranchId))

	if err != nil {
		return nil, err
	}
	pbBranchMembers := h.mapper.BranchMembersToPb(branchMembers)
	return &organizationpb.GetBranchMemberByIdsResponse{
		Data: pbBranchMembers,
	}, nil
}

func (h *InternalOrganizationHandler) GetDealById(ctx context.Context, req *organizationpb.GetDealByIdRequest) (*organizationpb.GetDealByIdResponse, error) {
	deal, err := h.dealRepo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	pbDeal := h.mapper.GroupDealToPb(deal)
	return &organizationpb.GetDealByIdResponse{
		Deal: pbDeal,
	}, nil
}

func (h *InternalOrganizationHandler) AddProductToDeal(ctx context.Context, req *organizationpb.AddProductToDealRequest) (*sharepb.SubmitResponse, error) {
	_, err := h.dealUsecase.SaveDealProduct(ctx, req.DealId, []uint64{req.ProductId})
	if err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Message: "success",
	}, nil
}

func (h *InternalOrganizationHandler) CheckUsersInOrganization(ctx context.Context, req *organizationpb.CheckUsersInOrganizationRequest) (*organizationpb.CheckUsersInOrganizationResponse, error) {
	users, err := h.organizationMemberUsecase.CheckUsersInOrganization(ctx, req.OrganizationId, req.UserIds)
	if err != nil {
		return nil, err
	}
	checkUsers := make([]*organizationpb.CheckUserInOrganization, len(users))
	for i, user := range users {
		fmt.Println("user", user)
		checkUsers[i] = &organizationpb.CheckUserInOrganization{UserId: user.UserId, IsExist: user.IsExist, JoinedAt: user.JoinedAt}
	}
	return &organizationpb.CheckUsersInOrganizationResponse{Data: checkUsers}, nil
}
