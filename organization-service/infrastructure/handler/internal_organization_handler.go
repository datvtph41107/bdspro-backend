package handler

import (
	"context"
	"organization/infrastructure/transformer"
	"organization/internal/domain/repository"
	organizationpb "pb/types/organization"
)

type InternalOrganizationHandler struct {
	organizationpb.UnimplementedInternalOrganizationServiceServer
	organizationRepo repository.OrganizationRepository
	groupRepo        repository.GroupRepository
	mapper           *transformer.InternalOrganizationMapper
}

func NewInternalOrganizationHandler(
	organizationRepo repository.OrganizationRepository,
	groupRepo repository.GroupRepository,
	mapper *transformer.InternalOrganizationMapper) *InternalOrganizationHandler {
	return &InternalOrganizationHandler{
		organizationRepo: organizationRepo,
		groupRepo:        groupRepo,
		mapper:           mapper,
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
