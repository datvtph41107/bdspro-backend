package repo

import (
	"context"
	"crm/internal/domain"
)

type FollowRepo interface {
	FollowerUser(c context.Context, profileId uint64, page int, pageSize int) ([]domain.Profile, error)
	FollowingUser(c context.Context, profileId uint64, page int, pageSize int) ([]domain.Profile, error)
	FollowUser(c context.Context, profileId uint64, followId uint64) (*domain.FollowEntity, error)
	UnfollowUser(c context.Context, profileId uint64, followId uint64) error
	BlockFollow(c context.Context, profileId uint64, targetId uint64) error
	GetFollowInfoCount(c context.Context, userID uint64) (int64, int64, int64, error)
}
