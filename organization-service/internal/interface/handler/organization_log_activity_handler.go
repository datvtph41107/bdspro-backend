package handler

import (
	"context"

	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/interface/transformer"
	"github.com/tuandq2112/organization-service/internal/interface/validator"
	"github.com/tuandq2112/organization-service/internal/usecase"
)

type OrganizationLogActivityHandler struct {
	organizationpb.UnimplementedOrganizationLogActivityServiceServer
	OrganizationLogActivityUsecase     usecase.OrganizationLogActivityUsecase
	OrganizationLogActivityTransformer transformer.OrganizationLogActivityTransformer
	OrganizationLogActivityValidator   validator.OrganizationLogActivityValidator
}

func NewOrganizationLogActivityHandler(
	organizationLogActivityUsecase usecase.OrganizationLogActivityUsecase,
	organizationLogActivityTransformer transformer.OrganizationLogActivityTransformer,
	organizationLogActivityValidator validator.OrganizationLogActivityValidator,
) *OrganizationLogActivityHandler {
	return &OrganizationLogActivityHandler{
		OrganizationLogActivityUsecase:     organizationLogActivityUsecase,
		OrganizationLogActivityTransformer: organizationLogActivityTransformer,
		OrganizationLogActivityValidator:   organizationLogActivityValidator,
	}
}

func (h *OrganizationLogActivityHandler) GetOrganizationLogActivityByID(ctx context.Context, req *organizationpb.GetOrganizationLogActivityByIDRequest) (*organizationpb.GetOrganizationLogActivityByIDResponse, error) {
	if err := h.OrganizationLogActivityValidator.ValidateGetOrganizationLogActivityByIDRequest(req); err != nil {
		return nil, err
	}

	log, err := h.OrganizationLogActivityUsecase.GetOrganizationLogActivityByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.OrganizationLogActivityTransformer.EntityToGetOrganizationLogActivityByIDResponse(log), nil
}

func (h *OrganizationLogActivityHandler) GetOrganizationLogActivitiesByOrganizationID(ctx context.Context, req *organizationpb.GetOrganizationLogActivitiesByOrganizationIDRequest) (*organizationpb.GetOrganizationLogActivitiesByOrganizationIDResponse, error) {
	if err := h.OrganizationLogActivityValidator.ValidateGetOrganizationLogActivitiesByOrganizationIDRequest(req); err != nil {
		return nil, err
	}

	logs, err := h.OrganizationLogActivityUsecase.GetOrganizationLogActivitiesByOrganizationID(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return h.OrganizationLogActivityTransformer.EntityToGetOrganizationLogActivitiesByOrganizationIDResponse(logs), nil
}

func (h *OrganizationLogActivityHandler) GetOrganizationLogActivitiesByActorID(ctx context.Context, req *organizationpb.GetOrganizationLogActivitiesByActorIDRequest) (*organizationpb.GetOrganizationLogActivitiesByActorIDResponse, error) {
	if err := h.OrganizationLogActivityValidator.ValidateGetOrganizationLogActivitiesByActorIDRequest(req); err != nil {
		return nil, err
	}

	logs, err := h.OrganizationLogActivityUsecase.GetOrganizationLogActivitiesByActorID(ctx, req.ActorId)
	if err != nil {
		return nil, err
	}
	return h.OrganizationLogActivityTransformer.EntityToGetOrganizationLogActivitiesByActorIDResponse(logs), nil
} 