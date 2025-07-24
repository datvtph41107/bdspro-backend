package validator

import (
	organizationpb "pb/types/organization"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrganizationMemberValidator interface {
	ValidateCreateOrganizationMemberRequest(request *organizationpb.CreateOrganizationMemberRequest) error
	ValidateCreateOrganizationMemberBatchRequest(request *organizationpb.CreateOrganizationMemberBatchRequest) error
	ValidateUpdateOrganizationMemberRequest(request *organizationpb.UpdateOrganizationMemberRequest) error
}

type organizationMemberValidator struct{}

func NewOrganizationMemberValidator() OrganizationMemberValidator {
	return &organizationMemberValidator{}
}

func (o *organizationMemberValidator) ValidateCreateOrganizationMemberRequest(request *organizationpb.CreateOrganizationMemberRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.UserId == 0 {
		return status.Errorf(codes.InvalidArgument, "userId is required")
	}

	if request.Role == 0 {
		return status.Errorf(codes.InvalidArgument, "role is required")
	}

	return nil
}

func (o *organizationMemberValidator) ValidateCreateOrganizationMemberBatchRequest(request *organizationpb.CreateOrganizationMemberBatchRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if len(request.Members) == 0 {
		return status.Errorf(codes.InvalidArgument, "members list cannot be empty")
	}

	if len(request.Members) > 100 {
		return status.Errorf(codes.InvalidArgument, "cannot create more than 100 members at once")
	}

	for i, member := range request.Members {
		if err := o.ValidateCreateOrganizationMemberRequest(member); err != nil {
			return status.Errorf(codes.InvalidArgument, "member at index %d: %v", i, err)
		}
	}

	return nil
}

func (o *organizationMemberValidator) ValidateUpdateOrganizationMemberRequest(request *organizationpb.UpdateOrganizationMemberRequest) error {
	if request == nil {
		return status.Errorf(codes.InvalidArgument, "request cannot be nil")
	}

	if request.Id == 0 {
		return status.Errorf(codes.InvalidArgument, "id is required")
	}

	if request.Role == 0 {
		return status.Errorf(codes.InvalidArgument, "role is required")
	}

	return nil
}

