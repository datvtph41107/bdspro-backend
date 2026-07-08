package usecase

import (
	base_enum "base/enum"
	"context"
	"crm/internal/enums"
)

type PermissionUsecase interface {
	UserCanAccessTarget(c context.Context,
		ownerId uint64,
		ownerType base_enum.EOwnerOf,
		targetId uint64,
		targetType int8,
	) error
	UserInOwner(
		c context.Context,
		ownerId uint64,
		ownerType base_enum.EOwnerOf,
	) error
	HasRoleWithOwner(
		c context.Context,
		ownerId uint64,
		ownerType base_enum.EOwnerOf,
		roleValue enums.EAuth,
	) error
	// UserInOwner(
	// 	c context.Context,
	// 	ownerId uint64,
	// 	ownerType enums.EOwnerType,
	// ) error
}