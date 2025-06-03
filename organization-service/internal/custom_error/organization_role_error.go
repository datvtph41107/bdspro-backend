package custom_error

import (
	"errors"

	"google.golang.org/grpc/codes"
)


var (
	RoleKeyAlreadyExistsCode codes.Code = 4001
)

func RoleKeyAlreadyExists() error {
	return NewError(int32(RoleKeyAlreadyExistsCode), errors.New("role key already exists"))
}