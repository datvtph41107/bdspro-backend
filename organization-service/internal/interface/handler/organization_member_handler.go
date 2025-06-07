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
	return &organizationpb.SwitchOrganizationResponse{
		IsSuccess:   true,
		Id:          req.OrganizationId,
		UserId:      req.UserId,
		Role:        member.Role,
		Permissions: member.Permissions,
	}, nil
}

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

func (h *OrganizationMemberHandler) GetOrganizationMembers(ctx context.Context, req *organizationpb.GetOrganizationMembersRequest) (*organizationpb.GetOrganizationMembersResponse, error) {
	organizationMembers, err := h.OrganizationMemberUsecase.GetOrganizationMembers(ctx, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	return h.OrganizationMemberTransformer.EntityToGetOrganizationMembersResponse(organizationMembers), nil
}
