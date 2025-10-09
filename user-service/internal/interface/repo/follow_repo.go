package repo

import (
	models "user/internal/models"

	"github.com/gin-gonic/gin"
)

type IFollowRepo interface {
	FollowerUser(currentId uint64, page int, pageSize int) ([]models.Profile, error)
	FollowingUser(currentId uint64, page int, pageSize int) ([]models.Profile, error)
	FollowUser(c *gin.Context, profileId, followingID uint64) (*models.FollowEntity, error)
	UnfollowUser(c *gin.Context, profileId, followingID uint64) error
	BlockFollow(profileId, targetId uint64) error
	GetFollowInfoCount(userID uint64) (followerCount int64, followingCount int64, friendCount int64, err error)
	GetFollowing(profileId uint64, targetId uint64) (bool, error)
}
