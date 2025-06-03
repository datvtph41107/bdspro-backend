package custom_error

import (
	"errors"

	"google.golang.org/grpc/codes"
)


var (
	PermissionKeyAlreadyExistsCode codes.Code = 5002
)

func PermissionKeyAlreadyExists() error {
	return NewError(int32(PermissionKeyAlreadyExistsCode), errors.New("permission key already exists"))
}