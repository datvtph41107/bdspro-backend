package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type GroupValidator interface {
	ValidateCreateGroupRequest(request *organizationpb.CreateGroupRequest) error
	ValidateUpdateGroupRequest(request *organizationpb.UpdateGroupRequest) error
	ValidateDeleteGroupRequest(request *organizationpb.DeleteGroupRequest) error
	ValidateGetGroupRequest(request *organizationpb.GetGroupRequest) error
	ValidateGetGroupsRequest(request *organizationpb.GetGroupsRequest) error
}

type groupValidator struct{}

func NewGroupValidator() GroupValidator {
	return &groupValidator{}
}

func (v *groupValidator) ValidateCreateGroupRequest(request *organizationpb.CreateGroupRequest) error {
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

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupValidator) ValidateUpdateGroupRequest(request *organizationpb.UpdateGroupRequest) error {
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

func (v *groupValidator) ValidateDeleteGroupRequest(request *organizationpb.DeleteGroupRequest) error {
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

func (v *groupValidator) ValidateGetGroupRequest(request *organizationpb.GetGroupRequest) error {
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

func (v *groupValidator) ValidateGetGroupsRequest(request *organizationpb.GetGroupsRequest) error {
	details := []protoadapt.MessageV1{}

	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}
