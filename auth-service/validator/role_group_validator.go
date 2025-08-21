package validator

import (
	_errors "common/errors"
	authpb "pb/types/auth"
)

type RoleGroupValidator struct{}

func NewRoleGroupValidator() *RoleGroupValidator {
	return &RoleGroupValidator{}
}

// ValidateCreateRoleGroup validate request tạo mới role group
func (v *RoleGroupValidator) ValidateCreateRoleGroup(req *authpb.RoleGroupRequest) error {
	if req.GroupName == "" {
		return _errors.ReturnError(400, "Tên role group không được để trống")
	}

	if len(req.GroupName) > 20 {
		return _errors.ReturnError(400, "Tên role group không được quá 20 ký tự")
	}

	if req.Key == 0 {
		return _errors.ReturnError(400, "Key role group không được để trống")
	}

	if req.Description != "" && len(req.Description) > 500 {
		return _errors.ReturnError(400, "Mô tả role group không được quá 500 ký tự")
	}

	return nil
}

// ValidateUpdateRoleGroup validate request cập nhật role group
func (v *RoleGroupValidator) ValidateUpdateRoleGroup(req *authpb.RoleGroupRequest) error {
	if req.Id == 0 {
		return _errors.ReturnError(400, "ID role group không được để trống")
	}

	return v.ValidateCreateRoleGroup(req)
}
