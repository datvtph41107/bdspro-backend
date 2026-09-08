package postgres

import (
	_db "common/db"
	_routes "common/routes"
	"errors"
	"time"
	models "user/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.IFollowRepo
type FollowPostgres struct {
	DB *gorm.DB
}

func NewFollowPostgres(DB *gorm.DB) *FollowPostgres {
	return &FollowPostgres{
		DB: DB,
	}
}

func (repo *FollowPostgres) FollowerUser(currentId uint64, page int, pageSize int) ([]models.Profile, error) {
	var friends []models.Profile

	query := repo.DB.Model(&models.Profile{}).
		Joins("LEFT JOIN follow ON (follow.created_by = profile_transfer.profile_id) "+
			"						AND profile_transfer.profile_id <> ? ", currentId).
		Where("follow.deleted_at is null and follow.following_id = ?", currentId)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}

func (repo *FollowPostgres) FollowingUser(currentId uint64, page int, pageSize int) ([]models.Profile, error) {
	var friends []models.Profile

	query := repo.DB.Model(&models.Profile{}).
		Joins("LEFT JOIN follow ON (follow.following_id = profile_transfer.profile_id) "+
			"						AND profile_transfer.profile_id <> ? ", currentId).
		Where("follow.deleted_at is null and follow.created_by = ?", currentId)

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}

func (repo *FollowPostgres) FollowUser(c *gin.Context, profileId, followingID uint64) (*models.FollowEntity, error) {
	if profileId == followingID {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể tự follow chính mình",
		}
	}

	follow := models.FollowEntity{
		FollowingID: followingID,
	}

	var exists bool
	// Kiểm tra nếu đã follow trước đó
	err := repo.DB.Raw(`SELECT EXISTS (
							SELECT 1 FROM follow
							WHERE created_by = ? AND following_id = ? AND deleted_at IS NULL)`,
		profileId, followingID).Scan(&exists).Error

	if err != nil || exists {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Đã follow người này rồi",
		}
	}

	err = _db.SaveWithAudit(c, &follow)

	return &follow, err
}

func (repo *FollowPostgres) UnfollowUser(c *gin.Context, profileId, followingID uint64) error {
	// Tìm bản ghi follow
	result := repo.DB.Model(&models.FollowEntity{}).
		Where("created_by = ? AND following_id = ?", profileId, followingID).
		Update("deleted_at", time.Now())

	if result.RowsAffected == 0 {
		return errors.New("bạn chưa follow người này")
	}
	return nil
}

func (repo *FollowPostgres) BlockFollow(profileId, targetId uint64) error {
	repo.DB.Model(&models.FollowEntity{}).
		Where("(created_by = ? AND following_id = ?) OR (created_by = ? AND following_id = ?)",
			profileId, targetId,
			targetId, profileId,
		).
		Update("deleted_at", time.Now())

	return nil
}

func (repo *FollowPostgres) GetFollowInfoCount(userID uint64) (followerCount int64, followingCount int64, friendCount int64, err error) {
	err = repo.DB.Model(&models.FollowEntity{}).
		Where("following_id = ?", userID).
		Count(&followerCount).
		Error
	if err != nil {
		return 0, 0, 0, err
	}

	err = repo.DB.Model(&models.FollowEntity{}).
		Where("created_by = ?", userID).
		Count(&followingCount).
		Error
	if err != nil {
		return 0, 0, 0, err
	}

	err = repo.DB.Model(&models.FriendEntity{}).
		Where("(created_by = ? or receiver_id = ?) and status = 2", userID, userID).
		Count(&friendCount).
		Error
	if err != nil {
		return 0, 0, 0, err
	}

	return followerCount, followingCount, friendCount, nil
}

func (repo *FollowPostgres) GetFollowing(profileId uint64, targetId uint64) (bool, error) {
	var exists bool

	err := repo.DB.
		Raw(`
			SELECT EXISTS (
				SELECT 1 FROM follow
				WHERE created_by = ? AND following_id = ? AND deleted_at IS NULL
			)
		`, profileId, targetId).
		Scan(&exists).Error

	return exists, err
}
