package postgre

import (
	_db "common/db"
	_routes "common/routes"
	"crm/internal/domain"
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

// func (repo *FriendRepo) SearchFriend(currentId uint64, page int, pageSize int) ([]domain.ContactEntity, error) {
// 	var friends []domain.ContactEntity

// 	query := repo.DB.Model(&domain.ContactEntity{}).
// 		Joins("JOIN follow ON follow.created_by = contact.id").
// 		Where("(friend.created_by = ? OR friend.receiver_id = ?) AND friend.status = 2", currentId, currentId)

// 	// if dto.Text != nil && *dto.Text != "" {
// 	// 	query = query.Where("full_name LIKE ? OR username LIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
// 	// }

//		err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
//		return friends, err
//	}
func (repo *FollowRepo) FollowerUser(currentId uint64, page int, pageSize int) ([]domain.Profile, error) {
	var friends []domain.Profile

	query := repo.DB.Model(&domain.Profile{}).
		Joins("LEFT JOIN follow ON (follow.created_by = profile_transfer.profile_id) "+
			"						AND profile_transfer.profile_id <> ? ", currentId).
		Where("follow.deleted_at is null and follow.following_id = ?", currentId)

	// if dto.Text != nil && *dto.Text != "" {
	// 	query = query.Where("full_name LIKE ? OR username LIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
	// }

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}
func (repo *FollowRepo) FollowingUser(currentId uint64, page int, pageSize int) ([]domain.Profile, error) {
	var friends []domain.Profile

	query := repo.DB.Model(&domain.Profile{}).
		Joins("LEFT JOIN follow ON (follow.following_id = profile_transfer.profile_id) "+
			"						AND profile_transfer.profile_id <> ? ", currentId).
		Where("follow.deleted_at is null and follow.created_by = ?", currentId)

	// if dto.Text != nil && *dto.Text != "" {
	// 	query = query.Where("full_name LIKE ? OR username LIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
	// }

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}

func (repo *FollowRepo) FollowUser(c *gin.Context, profileId, followingID uint64) (*domain.FollowEntity, error) {
	if profileId == followingID {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Không thể tự follow chính mình",
		}
	}

	follow := domain.FollowEntity{
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
	result := repo.DB.Model(&domain.FollowEntity{}).
		Where("created_by = ? AND following_id = ?", profileId, followingID).
		Update("deleted_at", time.Now())

	if result.RowsAffected == 0 {
		return errors.New("bạn chưa follow người này")
	}
	return nil
}

func (repo *FollowRepo) BlockFollow(profileId, targetId uint64) error {
	repo.DB.Model(&domain.FollowEntity{}).
		Where("(created_by = ? AND following_id = ?) OR (created_by = ? AND following_id = ?)",
			profileId, targetId,
			targetId, profileId,
		).
		Update("deleted_at", time.Now())

	return nil
}

func (repo *FollowRepo) GetFollowInfoCount(userID uint64) (followerCount int64, followingCount int64, friendCount int64, err error) {
	err = repo.DB.Model(&domain.FollowEntity{}).
		Where("following_id = ?", userID).
		Count(&followerCount).
		Error
	if err != nil {
		return 0, 0, 0, err
	}

	err = repo.DB.Model(&domain.FollowEntity{}).
		Where("created_by = ?", userID).
		Count(&followingCount).
		Error
	if err != nil {
		return 0, 0, 0, err
	}

	err = repo.DB.Model(&domain.FriendEntity{}).
		Where("(created_by = ? or receiver_id = ?) and status = 2", userID, userID).
		Count(&friendCount).
		Error
	if err != nil {
		return 0, 0, 0, err
	}

	return followerCount, followingCount, friendCount, nil
}
