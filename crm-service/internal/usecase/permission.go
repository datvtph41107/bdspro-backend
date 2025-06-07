package usecase

import (
	"context"
	"crm/internal/enums"
)

type PermissionUsecase interface {
	UserCanAccessTarget(c context.Context,
		ownerId uint64,
		ownerType enums.EOwnerType,
		targetId uint64,
		targetType int8,
	) error
	UserInOwner(
		c context.Context,
		ownerId uint64,
		ownerType enums.EOwnerType,
	) error
	HasRoleWithOwner(
		c context.Context,
		ownerId *uint64,
		ownerType enums.EOwnerType,
		roleValue enums.EAuth,
	) error
	// UserInOwner(
	// 	c context.Context,
	// 	ownerId uint64,
	// 	ownerType enums.EOwnerType,
	// ) error
}
