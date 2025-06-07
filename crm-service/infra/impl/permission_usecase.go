package impl

import (
	"context"
	"crm/internal/enums"
)

// @bind: crm/internal/usecase.PermissionUsecase
type PermissionUsecase struct {
}

func NewPermissionUsecase() *PermissionUsecase {
	return &PermissionUsecase{}
}

func (u *PermissionUsecase) UserInOwner(c context.Context, ownerID uint64, ownerType enums.EOwnerType) error {
	return nil
}

func (u *PermissionUsecase) UserCanAccessTarget(c context.Context, ownerID uint64, ownerType enums.EOwnerType, targetID uint64, targetType int8) error {
	return nil
}

func (u *PermissionUsecase) HasRoleWithOwner(c context.Context, ownerID *uint64, ownerType enums.EOwnerType, roleValue enums.EAuth) error {
	return nil
}

func (u *PermissionUsecase) ClearProfileID(c context.Context, profileID uint64) error {
	return nil
}
