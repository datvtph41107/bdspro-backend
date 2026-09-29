package validator

import (
	_errors "common/errors"
	authpb "pb/types/auth"
	"user/internal"
)

type RoleGroupValidator struct{}

func NewRoleGroupValidator() *RoleGroupValidator {
	return &RoleGroupValidator{}
}

// ValidateCreateRoleGroup validate request tạo mới role group
func (v *RoleGroupValidator) ValidateCreateRoleGroup(req *authpb.RoleGroupRequest) error {
	if req.GroupName == "" {
		return _errors.ReturnError(service.RoleGroupNameRequired)
	}

	if len(req.GroupName) > 20 {
		return _errors.ReturnError(service.RoleGroupNameTooLong)
	}

	if req.Key == 0 {
		return _errors.ReturnError(service.RoleGroupKeyRequired)
	}

	if req.Description != "" && len(req.Description) > 500 {
		return _errors.ReturnError(service.RoleGroupDescriptionTooLong)
	}

	return nil
}

// ValidateUpdateRoleGroup validate request cập nhật role group
func (v *RoleGroupValidator) ValidateUpdateRoleGroup(req *authpb.RoleGroupRequest) error {
	if req.Id == 0 {
		return _errors.ReturnError(service.RoleGroupIDRequired)
	}

	return v.ValidateCreateRoleGroup(req)
}
