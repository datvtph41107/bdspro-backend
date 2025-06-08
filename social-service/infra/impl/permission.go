package impl

import (
	"context"
)

// @bind: social/internal/interface.PermissionUsecase
type PermissionImpl struct {
}

func NewPermissionImpl() *PermissionImpl {
	return &PermissionImpl{}
}

func (r *PermissionImpl) OwnerPost(ctx context.Context, postID uint64, profileID uint64) error {
	return nil
}

func (r *PermissionImpl) IsAdmin(ctx context.Context) bool {
	return true
}

func (r *PermissionImpl) OwnerComment(ctx context.Context, commentID uint64, profileID uint64) error {
	return nil
}
