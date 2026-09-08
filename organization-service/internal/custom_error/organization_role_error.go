package custom_error

import (
	_errors "common/errors"
	"errors"

	"google.golang.org/grpc/codes"
)

var (
	RoleKeyAlreadyExistsCode codes.Code = 4001
)

func RoleKeyAlreadyExists() error {
	return _errors.ThrowError(int32(RoleKeyAlreadyExistsCode), errors.New("role key already exists"))
}
