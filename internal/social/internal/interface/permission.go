package iusecase

import "context"

type PermissionUsecase interface {
	OwnerComment(ctx context.Context, commentID uint64, userID uint64) error
	IsAdmin(ctx context.Context) bool
}
