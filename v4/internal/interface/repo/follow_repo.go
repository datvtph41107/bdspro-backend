package repo

import (
	"context"
	entity "mdn4/modules/crm/entity"
)

type IFollowRepo interface {
	FollowerUser(currentId uint64, page int, pageSize int) ([]entity.Profile, error)
	FollowingUser(currentId uint64, page int, pageSize int) ([]entity.Profile, error)
	FollowUser(ctx context.Context, profileId, followingID uint64) (*entity.FollowEntity, error)
	UnfollowUser(ctx context.Context, profileId, followingID uint64) error
	BlockFollow(profileId, targetId uint64) error
	GetFollowInfoCount(userID uint64) (followerCount int64, followingCount int64, friendCount int64, err error)
	GetFollowing(profileId uint64, targetId uint64) (bool, error)
}
