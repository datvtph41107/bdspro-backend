package handler

import (
	"context"

	organizationpb "pb/types/organization"

	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/usecase"
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
	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}
	logs, total, err := h.OrganizationLogActivityUsecase.GetOrganizationLogActivitiesByOrganizationID(ctx, req.OrganizationId, page, size)
	if err != nil {
		return nil, err
	}
	return h.OrganizationLogActivityTransformer.EntityToGetOrganizationLogActivitiesByOrganizationIDResponse(logs, total), nil
}

func (h *OrganizationLogActivityHandler) GetOrganizationLogActivitiesByActorID(ctx context.Context, req *organizationpb.GetOrganizationLogActivitiesByActorIDRequest) (*organizationpb.GetOrganizationLogActivitiesByActorIDResponse, error) {
	if err := h.OrganizationLogActivityValidator.ValidateGetOrganizationLogActivitiesByActorIDRequest(req); err != nil {
		return nil, err
	}

	page := 0
	size := 10
	if req.Page != nil {
		page = int(*req.Page)
	}
	if req.Size != nil {
		size = int(*req.Size)
	}

	logs, total, err := h.OrganizationLogActivityUsecase.GetOrganizationLogActivitiesByActorID(ctx, req.ActorId, page, size)
	if err != nil {
		return nil, err
	}
	return h.OrganizationLogActivityTransformer.EntityToGetOrganizationLogActivitiesByActorIDResponse(logs, total), nil
}
