package handler

import (
	"context"
	"organization/infrastructure/transformer"
	"organization/internal/domain/repository"
	"organization/internal/usecase"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
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
	dealMemberRepo         repository.DealMemberRepository
	dealUsecase            usecase.DealUsecase
}

func NewInternalOrganizationHandler(
	organizationRepo repository.OrganizationRepository,
	groupRepo repository.GroupRepository,
	dealInvitationRepo repository.DealInvitationRepository,
	groupMemberRepo repository.GroupMemberRepository,
	organizationMemberRepo repository.OrganizationMemberRepository,
	dealRepo repository.GroupDealRepository,
	dealTransactionRepo repository.DealTransactionRepository,
	dealMemberRepo repository.DealMemberRepository,
	dealUsecase usecase.DealUsecase,
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
		dealMemberRepo:         dealMemberRepo,
		dealUsecase:            dealUsecase,
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
