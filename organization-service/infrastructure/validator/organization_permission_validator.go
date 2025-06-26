package validator

import (
	"organization/internal/custom_error"
	organizationpb "pb/types/organization"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type OrganizationPermissionValidator interface {
	ValidateCreatePermissionRequest(req *organizationpb.CreatePermissionRequest) error
	ValidateUpdatePermissionRequest(req *organizationpb.UpdatePermissionRequest) error
	ValidateDeletePermissionRequest(req *organizationpb.DeletePermissionRequest) error
	ValidateGetPermissionsRequest(req *organizationpb.GetPermissionsRequest) error
}

type organizationPermissionValidator struct {
}

func NewOrganizationPermissionValidator() OrganizationPermissionValidator {
	return &organizationPermissionValidator{}
}

func (v *organizationPermissionValidator) ValidateCreatePermissionRequest(req *organizationpb.CreatePermissionRequest) error {
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

func (v *organizationPermissionValidator) ValidateUpdatePermissionRequest(req *organizationpb.UpdatePermissionRequest) error {
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

func (v *organizationPermissionValidator) ValidateDeletePermissionRequest(req *organizationpb.DeletePermissionRequest) error {
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

func (v *organizationPermissionValidator) ValidateGetPermissionsRequest(req *organizationpb.GetPermissionsRequest) error {
	details := []protoadapt.MessageV1{}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}

	return nil
}
