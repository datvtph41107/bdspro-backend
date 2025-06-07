package handler

import (
	"context"
	organizationpb "pb/types/organization"
	"time"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
	"github.com/tuandq2112/organization-service/pkg/utils"
)

type OrganizationHandler struct {
	organizationpb.UnimplementedOrganizationServiceServer
	OrganizationUsecase     usecase.OrganizationUsecase
	OrganizationTransformer transformer.OrganizationTransformer
	OrganizationValidator   validator.OrganizationValidator
}

func NewOrganizationHandler(organizationUsecase usecase.OrganizationUsecase,
	organizationTransformer transformer.OrganizationTransformer,
	organizationValidator validator.OrganizationValidator,
) *OrganizationHandler {
	return &OrganizationHandler{
		OrganizationUsecase:     organizationUsecase,
		OrganizationTransformer: organizationTransformer,
		OrganizationValidator:   organizationValidator,
	}
}

func (h *OrganizationHandler) CreateOrganization(ctx context.Context, req *organizationpb.CreateOrganizationRequest) (*organizationpb.CreateOrganizationResponse, error) {
	organization := h.OrganizationTransformer.CreateOrganizationRequestToEntity(req)
	organization.CreatedBy = utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	member := &entity.OrganizationMember{
		UserID:   organization.CreatedBy,
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

func (h *OrganizationHandler) UpdateOrganization(ctx context.Context, req *organizationpb.UpdateOrganizationRequest) (*organizationpb.UpdateOrganizationResponse, error) {
	if err := h.OrganizationValidator.ValidateUpdateOrganizationRequest(req); err != nil {
		return nil, err
	}

	organization := h.OrganizationTransformer.UpdateOrganizationRequestToEntity(req)
	organization.UpdatedBy = utils.GetUserID(ctx, appconfig.USER_CONTEXT)

	organization, err := h.OrganizationUsecase.UpdateOrganization(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT), organization)
	if err != nil {
		return nil, err
	}
	return h.OrganizationTransformer.EntityToUpdateOrganizationResponse(organization), nil
}

func (h *OrganizationHandler) GetOrganizations(ctx context.Context, req *organizationpb.GetOrganizationsRequest) (*organizationpb.GetOrganizationsResponse, error) {
	organizations, err := h.OrganizationUsecase.GetOrganizations(ctx, utils.GetUserID(ctx, appconfig.USER_CONTEXT))
	if err != nil {
		return nil, err
	}
	return h.OrganizationTransformer.EntityToGetOrganizationsResponse(organizations), nil
}