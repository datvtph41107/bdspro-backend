package iusecase

import "context"

type PermissionUsecase interface {
	OwnerPost(ctx context.Context, postID uint64, userID uint64) error
	OwnerComment(ctx context.Context, commentID uint64, userID uint64) error
	IsAdmin(ctx context.Context) bool
}
