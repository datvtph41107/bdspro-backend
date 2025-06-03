package validator

import (
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type OrganizationRoleValidator interface {
	ValidateCreateRoleRequest(req *organizationpb.CreateRoleRequest) error
	ValidateUpdateRoleRequest(req *organizationpb.UpdateRoleRequest) error
	ValidateDeleteRoleRequest(req *organizationpb.DeleteRoleRequest) error
	ValidateGetRolesRequest(req *organizationpb.GetRolesRequest) error
}

type organizationRoleValidator struct {
}

func NewOrganizationRoleValidator() OrganizationRoleValidator {
	return &organizationRoleValidator{}
}

func (v *organizationRoleValidator) ValidateCreateRoleRequest(req *organizationpb.CreateRoleRequest) error {

	details := []protoadapt.MessageV1{}

	if req.Name == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name is required",
				},
			},
		})
	}
	if req.Key == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "key",
					Description: "key is required",
				},
			},
		})
	}

	if req.OrganizationId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "organizationId",
					Description: "organizationId is required",
				},
			},
		})
	}


	return nil
}

func (v *organizationRoleValidator) ValidateUpdateRoleRequest(req *organizationpb.UpdateRoleRequest) error {
	details := []protoadapt.MessageV1{}

	if req.Name == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "name",
					Description: "name is required",
				},
			},
		})
	}
	if req.Key == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "key",
					Description: "key is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
}

func (v *organizationRoleValidator) ValidateDeleteRoleRequest(req *organizationpb.DeleteRoleRequest) error {
	details := []protoadapt.MessageV1{}

	if req.Id == 0 {
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

func (v *organizationRoleValidator) ValidateGetRolesRequest(req *organizationpb.GetRolesRequest) error {
	details := []protoadapt.MessageV1{}

	if req.OrganizationId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "organizationId",	
					Description: "organizationId is required",
				},
			},
		})
	}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)	
	}

	return nil
}