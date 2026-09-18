package custom_error

import (
	"errors"

	"google.golang.org/grpc/codes"
)

func PermissionKeyAlreadyExists() error {
	return organizationError(codes.AlreadyExists, "ORGANIZATION_PERMISSION_EXISTS", errors.New("permission key already exists"))
}
