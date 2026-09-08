package client

import (
	base_enum "base/enum"
	"context"
	"crm/internal/enums"
)

// @bind: crm/internal/usecase.PermissionUsecase
type PermissionUsecase struct {
}

func NewPermissionUsecase() *PermissionUsecase {
	return &PermissionUsecase{}
}

func (u *PermissionUsecase) UserInOwner(c context.Context, ownerID uint64, ownerType base_enum.EOwnerOf) error {
	return nil
}

func (u *PermissionUsecase) UserCanAccessTarget(c context.Context, ownerID uint64, ownerType base_enum.EOwnerOf, targetID uint64, targetType int8) error {
	return nil
}

func (u *PermissionUsecase) HasRoleWithOwner(c context.Context, ownerID uint64, ownerType base_enum.EOwnerOf, roleValue enums.EAuth) error {
	return nil
}

func (u *PermissionUsecase) ClearProfileID(c context.Context, profileID uint64) error {
	return nil
}