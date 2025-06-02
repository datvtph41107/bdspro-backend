package handler

import (
	"context"
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationMemberHandler struct {
	organizationpb.UnimplementedOrganizationMemberServiceServer
	OrganizationMemberUsecase usecase.OrganizationMemberUsecase
}

func NewOrganizationMemberHandler(organizationMemberUsecase usecase.OrganizationMemberUsecase) *OrganizationMemberHandler {
	return &OrganizationMemberHandler{OrganizationMemberUsecase: organizationMemberUsecase}
}

func (h *OrganizationMemberHandler) SwitchOrganization(ctx context.Context, req *organizationpb.SwitchOrganizationRequest) (*organizationpb.SwitchOrganizationResponse, error) {
	exist, err := h.OrganizationMemberUsecase.ExistByUserIdAndOrganizationId(ctx, req.UserId, req.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, status.Errorf(codes.NotFound, "user not in organization")
	}
	return &organizationpb.SwitchOrganizationResponse{
		IsSuccess: true,
		Id:        req.OrganizationId,
		UserId:    req.UserId,
	}, nil
}
