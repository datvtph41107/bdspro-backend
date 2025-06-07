package validator

import (
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type GroupValidator interface {
	ValidateCreateGroupRequest(request *grouppb.CreateGroupRequest) error
	ValidateUpdateGroupRequest(request *grouppb.UpdateGroupRequest) error
	ValidateDeleteGroupRequest(request *grouppb.DeleteGroupRequest) error
	ValidateGetGroupRequest(request *grouppb.GetGroupRequest) error
	ValidateGetGroupsRequest(request *grouppb.GetGroupsRequest) error
}

type groupValidator struct{}

func NewGroupValidator() GroupValidator {
	return &groupValidator{}
}

func (v *groupValidator) ValidateCreateGroupRequest(request *grouppb.CreateGroupRequest) error {
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

func (v *groupValidator) ValidateUpdateGroupRequest(request *grouppb.UpdateGroupRequest) error {
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

func (v *groupValidator) ValidateDeleteGroupRequest(request *grouppb.DeleteGroupRequest) error {
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

func (v *groupValidator) ValidateGetGroupRequest(request *grouppb.GetGroupRequest) error {
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

func (v *groupValidator) ValidateGetGroupsRequest(request *grouppb.GetGroupsRequest) error {
	details := []protoadapt.MessageV1{}
	
	
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
} 