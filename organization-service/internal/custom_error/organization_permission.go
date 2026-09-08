package custom_error

import (
	_errors "common/errors"
	"errors"

	"google.golang.org/grpc/codes"
)

var (
	PermissionKeyAlreadyExistsCode codes.Code = 5002
)

func PermissionKeyAlreadyExists() error {
	return _errors.ThrowError(int32(PermissionKeyAlreadyExistsCode), errors.New("permission key already exists"))
}
