package iusecase

import "context"

type IGroupAuth interface {
	HasPermission(ctx context.Context, permission string) bool
	IsMemberGroup(ctx context.Context, groupID uint32, userID uint32) error
	IsLeaderGroup(ctx context.Context, groupID uint32, userID uint32) error
} 