package custom_error

import (
	"errors"

	"google.golang.org/grpc/codes"
)

func RoleKeyAlreadyExists() error {
	return organizationError(codes.AlreadyExists, "ORGANIZATION_ROLE_EXISTS", errors.New("role key already exists"))
}
