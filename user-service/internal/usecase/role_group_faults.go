package usecase

import (
	"common/fault"
	"user/internal/domain/access"
)

// RoleGroupNotFoundCode is the stable application identity for a missing IAM
// role group. Protocol mappings consume this identity; messages do not define
// semantics.
const RoleGroupNotFoundCode = "iam.role_group.not_found"

func roleGroupNotFoundFault() error {
	return fault.Wrap(
		access.ErrRoleGroupNotFound,
		fault.KindNotFound,
		RoleGroupNotFoundCode,
		"role group not found",
	)
}
