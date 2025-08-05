package handler

import (
	"context"
	"organization/infrastructure/transformer"
	"organization/internal/domain/repository"
	organizationpb "pb/types/organization"
)

type InternalOrganizationHandler struct {
	organizationpb.UnimplementedInternalOrganizationServiceServer
	organizationRepo       repository.OrganizationRepository
	groupRepo              repository.GroupRepository
	mapper                 *transformer.InternalOrganizationMapper
	dealInvitationRepo     repository.DealInvitationRepository
	groupMemberRepo        repository.GroupMemberRepository
	organizationMemberRepo repository.OrganizationMemberRepository
	dealRepo               repository.GroupDealRepository
	dealTransactionRepo    repository.DealTransactionRepository
}

func NewInternalOrganizationHandler(
	organizationRepo repository.OrganizationRepository,
	groupRepo repository.GroupRepository,
	dealInvitationRepo repository.DealInvitationRepository,
	groupMemberRepo repository.GroupMemberRepository,
	organizationMemberRepo repository.OrganizationMemberRepository,
	dealRepo repository.GroupDealRepository,
	dealTransactionRepo repository.DealTransactionRepository,
	mapper *transformer.InternalOrganizationMapper) *InternalOrganizationHandler {
	return &InternalOrganizationHandler{
		organizationRepo:       organizationRepo,
		groupRepo:              groupRepo,
		dealInvitationRepo:     dealInvitationRepo,
		mapper:                 mapper,
		groupMemberRepo:        groupMemberRepo,
		organizationMemberRepo: organizationMemberRepo,
		dealRepo:               dealRepo,
		dealTransactionRepo:    dealTransactionRepo,
	}
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
	dealMembers, err := h.dealInvitationRepo.GetByIds(ctx, req.GroupId, req.Ids)
	if err != nil {
		return nil, err
	}
	pbDealMembers := h.mapper.DealInvitationsToPb(dealMembers)
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

func (h *InternalOrganizationHandler) GetDealTransactionsByIds(ctx context.Context, req *organizationpb.GetDealTransactionsByIdsRequest) (*organizationpb.GetDealTransactionsByIdsResponse, error) {
	dealTransactions, err := h.dealTransactionRepo.GetByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	pbDealTransactions := h.mapper.DealTransactionsToPb(dealTransactions)
	return &organizationpb.GetDealTransactionsByIdsResponse{
		DealTransactions: pbDealTransactions,
	}, nil
}
