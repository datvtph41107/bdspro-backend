package validator

import (
	organizationpb "pb/types/organization"

	"organization/internal/custom_error"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type GroupSettingValidator interface {
	ValidateCreateGroupSettingRequest(request *organizationpb.CreateGroupSettingRequest) error
	ValidateUpdateGroupSettingRequest(request *organizationpb.UpdateGroupSettingRequest) error
	ValidateDeleteGroupSettingRequest(request *organizationpb.DeleteGroupSettingRequest) error
	ValidateGetGroupSettingRequest(request *organizationpb.GetGroupSettingRequest) error
	ValidateGetGroupSettingsRequest(request *organizationpb.GetGroupSettingsRequest) error
}

type groupSettingValidator struct{}

func NewGroupSettingValidator() GroupSettingValidator {
	return &groupSettingValidator{}
}

func (v *groupSettingValidator) ValidateCreateGroupSettingRequest(request *organizationpb.CreateGroupSettingRequest) error {
	details := []protoadapt.MessageV1{}
	if request.GroupId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "group_id",
					Description: "group_id is required",
				},
			},
		})
	}
	if request.ConfigKey == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "config_key",
					Description: "config_key is required",
				},
			},
		})
	}
	if request.ConfigValue == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "config_value",
					Description: "config_value is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupSettingValidator) ValidateUpdateGroupSettingRequest(request *organizationpb.UpdateGroupSettingRequest) error {
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
	if request.ConfigKey == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "config_key",
					Description: "config_key is required",
				},
			},
		})
	}
	if request.ConfigValue == "" {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "config_value",
					Description: "config_value is required",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}

func (v *groupSettingValidator) ValidateDeleteGroupSettingRequest(request *organizationpb.DeleteGroupSettingRequest) error {
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

func (v *groupSettingValidator) ValidateGetGroupSettingRequest(request *organizationpb.GetGroupSettingRequest) error {
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

func (v *groupSettingValidator) ValidateGetGroupSettingsRequest(request *organizationpb.GetGroupSettingsRequest) error {
	details := []protoadapt.MessageV1{}
	if request.GroupId == 0 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "group_id",
					Description: "group_id is required",
				},
			},
		})
	}
	if request.Page != nil && *request.Page < 1 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "page",
					Description: "page must be greater than 0",
				},
			},
		})
	}
	if request.Size != nil && *request.Size < 1 {
		details = append(details, &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       "size",
					Description: "size must be greater than 0",
				},
			},
		})
	}
	if len(details) > 0 {
		return custom_error.InvalidRequest(details...)
	}
	return nil
}
