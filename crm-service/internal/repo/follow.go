package repo

import (
	"crm/internal/domain"
)

type FollowRepo interface {
	FollowerUser(profileId uint64, page int, pageSize int) ([]domain.Profile, error)
	FollowingUser(profileId uint64, page int, pageSize int) ([]domain.Profile, error)
	FollowUser(profileId uint64, followId uint64) (*domain.FollowEntity, error)
	UnfollowUser(profileId uint64, followId uint64) (*domain.FollowEntity, error)
	BlockFollow(profileId uint64, targetId uint64) error
	GetFollowInfoCount(profileId uint64) (int64, int64, int64, error)
}
