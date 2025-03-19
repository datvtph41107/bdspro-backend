package repo

import (
	_db "common/db"
	_routes "common/routes"
	"crm/models"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FollowRepo struct {
	DB *gorm.DB
}

func NewFollowRepo(DB *gorm.DB) *FollowRepo {
	return &FollowRepo{
		DB: DB,
	}
}

// func (repo *FriendRepo) SearchFriend(currentId uint64, page int, pageSize int) ([]models.ContactEntity, error) {
// 	var friends []models.ContactEntity

// 	query := repo.DB.Model(&models.ContactEntity{}).
// 		Joins("JOIN follow ON follow.created_by = contact.id").
// 		Where("(friend.created_by = ? OR friend.receiver_id = ?) AND friend.status = 2", currentId, currentId)

// 	// if dto.Text != nil && *dto.Text != "" {
// 	// 	query = query.Where("full_name LIKE ? OR username LIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
// 	// }

//		err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
//		return friends, err
//	}
func (repo *FollowRepo) FollowerUser(currentId uint64, page int, pageSize int) ([]models.Profile, error) {
	var friends []models.Profile

	query := repo.DB.Model(&models.Profile{}).
		Joins("LEFT JOIN follow ON (follow.created_by = profile_transfer.profile_id) "+
			"						AND profile_transfer.profile_id <> ? ", currentId).
		Where("follow.deleted_at is null and follow.following_id = ?", currentId)

	// if dto.Text != nil && *dto.Text != "" {
	// 	query = query.Where("full_name LIKE ? OR username LIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
	// }

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}
func (repo *FollowRepo) FollowingUser(currentId uint64, page int, pageSize int) ([]models.Profile, error) {
	var friends []models.Profile

	query := repo.DB.Model(&models.Profile{}).
		Joins("LEFT JOIN follow ON (follow.following_id = profile_transfer.profile_id) "+
			"						AND profile_transfer.profile_id <> ? ", currentId).
		Where("follow.deleted_at is null and follow.created_by = ?", currentId)

	// if dto.Text != nil && *dto.Text != "" {
	// 	query = query.Where("full_name LIKE ? OR username LIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
	// }

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}

func (repo *FollowRepo) FollowUser(c *gin.Context, profileId, followingID uint64) (*models.FollowEntity, error) {
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
							WHERE created_by = ? AND following_id = ? AND unfollowed_at IS NULL)`,
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

func (repo *FollowRepo) UnfollowUser(c *gin.Context, profileId, followingID uint64) error {
	// Tìm bản ghi follow
	result := repo.DB.Model(&models.FollowEntity{}).
		Where("created_by = ? AND following_id = ?", profileId, followingID).
		Update("deleted_at", time.Now())

	if result.RowsAffected == 0 {
		return errors.New("bạn chưa follow người này")
	}
	return nil
}

func (repo *FollowRepo) BlockFollow(profileId, targetId uint64) error {
	repo.DB.Model(&models.FollowEntity{}).
		Where("(created_by = ? AND following_id = ?) OR (created_by = ? AND following_id = ?)",
			profileId, targetId,
			targetId, profileId,
		).
		Update("deleted_at", time.Now())

	return nil
}
