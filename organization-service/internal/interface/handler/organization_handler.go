package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
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
	if err := h.OrganizationValidator.ValidateCreateOrganizationRequest(req); err != nil {
		return nil, err
	}
	organization, err := h.OrganizationUsecase.CreateOrganization(ctx, organization)
	if err != nil {
		return nil, err
	}
	return h.OrganizationTransformer.EntityToCreateOrganizationResponse(organization), nil
}
