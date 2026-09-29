package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type OrganizationBranchValidator interface {
	ValidateCreateOrganizationBranch(request *organizationpb.CreateOrganizationBranchRequest) error
	ValidateUpdateOrganizationBranch(request *organizationpb.UpdateOrganizationBranchRequest) error
	ValidateDeleteOrganizationBranch(request *organizationpb.DeleteOrganizationBranchRequest) error
	ValidateGetOrganizationBranchDetail(request *organizationpb.GetOrganizationBranchDetailRequest) error
	ValidateAddMemberToOrganizationBranch(request *organizationpb.AddMemberToOrganizationBranchRequest) error
	ValidateRemoveMemberFromOrganizationBranch(request *organizationpb.RemoveMemberFromOrganizationBranchRequest) error
	ValidateGetOrganizationBranchMembers(request *organizationpb.GetOrganizationBranchMembersRequest) error
}

type organizationBranchValidator struct{}

func NewOrganizationBranchValidator() OrganizationBranchValidator {
	return &organizationBranchValidator{}
}

func (v *organizationBranchValidator) ValidateCreateOrganizationBranch(request *organizationpb.CreateOrganizationBranchRequest) error {
	details := []protoadapt.MessageV1{}

	if request.Name == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name is required",
				},
			},
		})
	}
	if request.ManagerId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "managerId",
					Description: "managerId is required",
				},
			},
		})
	}
	if request.Type == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "type",
					Description: "type is required",
				},
			},
		})
	} else {
		// Validate type enum
		validTypes := map[string]bool{
			"branch":     true, // Chi nhánh
			"department": true, // Phòng ban
			"store":      true, // Cửa hàng
		}
		if !validTypes[request.Type] {
			details = append(details, &errdetails.BadRequest{
				FieldViolations: []*errdetails.BadRequest_FieldViolation{
					{
						Field:       "type",
						Description: "type must be one of: branch, department, store",
					},
				},
			})
		}
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *organizationBranchValidator) ValidateUpdateOrganizationBranch(request *organizationpb.UpdateOrganizationBranchRequest) error {
	details := []protoadapt.MessageV1{}

	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}
	if request.Name == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *organizationBranchValidator) ValidateDeleteOrganizationBranch(request *organizationpb.DeleteOrganizationBranchRequest) error {
	details := []protoadapt.MessageV1{}

	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *organizationBranchValidator) ValidateAddMemberToOrganizationBranch(request *organizationpb.AddMemberToOrganizationBranchRequest) error {
	details := []protoadapt.MessageV1{}

	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}

	if request.UserId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "userId",
					Description: "userId is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *organizationBranchValidator) ValidateRemoveMemberFromOrganizationBranch(request *organizationpb.RemoveMemberFromOrganizationBranchRequest) error {
	details := []protoadapt.MessageV1{}

	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}

	if request.UserId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "userId",
					Description: "userId is required",
				},
			},
		})
	}
	return nil
}

func (v *organizationBranchValidator) ValidateGetOrganizationBranchMembers(request *organizationpb.GetOrganizationBranchMembersRequest) error {
	details := []protoadapt.MessageV1{}

	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *organizationBranchValidator) ValidateGetOrganizationBranchDetail(request *organizationpb.GetOrganizationBranchDetailRequest) error {
	details := []protoadapt.MessageV1{}

	if request.Id == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "id",
					Description: "id is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}
