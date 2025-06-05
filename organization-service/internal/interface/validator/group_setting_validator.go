package validator

import (
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/protoadapt"
)

type GroupSettingValidator interface {
	ValidateCreateGroupSettingRequest(request *grouppb.CreateGroupSettingRequest) error
	ValidateUpdateGroupSettingRequest(request *grouppb.UpdateGroupSettingRequest) error
	ValidateDeleteGroupSettingRequest(request *grouppb.DeleteGroupSettingRequest) error
	ValidateGetGroupSettingRequest(request *grouppb.GetGroupSettingRequest) error
	ValidateGetGroupSettingsRequest(request *grouppb.GetGroupSettingsRequest) error
}

type groupSettingValidator struct{}

func NewGroupSettingValidator() GroupSettingValidator {
	return &groupSettingValidator{}
}

func (v *groupSettingValidator) ValidateCreateGroupSettingRequest(request *grouppb.CreateGroupSettingRequest) error {
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

func (v *groupSettingValidator) ValidateUpdateGroupSettingRequest(request *grouppb.UpdateGroupSettingRequest) error {
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

func (v *groupSettingValidator) ValidateDeleteGroupSettingRequest(request *grouppb.DeleteGroupSettingRequest) error {
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

func (v *groupSettingValidator) ValidateGetGroupSettingRequest(request *grouppb.GetGroupSettingRequest) error {
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

func (v *groupSettingValidator) ValidateGetGroupSettingsRequest(request *grouppb.GetGroupSettingsRequest) error {
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