package repo

import (
	"context"
	"crm/internal/domain"
)

type FollowRepo interface {
	FollowerUser(c context.Context, profileId uint64, page int, pageSize int) ([]domain.Profile, error)
	FollowingUser(c context.Context, profileId uint64, page int, pageSize int) ([]uint64, error)
	FollowUser(c context.Context, profileId uint64, followId uint64) (*domain.FollowEntity, error)
	BlockFollow(c context.Context, profileId uint64, targetId uint64) error
	GetFollowInfoCount(c context.Context, userID uint64) (int64, int64, int64, error)
	GetFollowing(ctx context.Context, profileId uint64, targetId uint64) (bool, error)
	GetFollowId(ctx context.Context, profileId uint64, targetId uint64) (*uint64, error)

	GetFollow(ctx context.Context, followerID, followingID uint64) (*domain.FollowEntity, error)
	CreateFollow(ctx context.Context, followerID, followingID uint64) (*domain.FollowEntity, error)
	Update(ctx context.Context, follow *domain.FollowEntity) error
	UnfollowUser(ctx context.Context, followerID, followingID uint64) error
}