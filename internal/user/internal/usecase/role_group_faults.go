package usecase

import (
	_errors "common/errors"
	"user/internal"
	"user/internal/domain/access"
)

func roleGroupNotFoundFault() error {
	return _errors.ReturnError(
		service.RoleGroupNotFound,
		_errors.WithCause(access.ErrRoleGroupNotFound),
	)
}
