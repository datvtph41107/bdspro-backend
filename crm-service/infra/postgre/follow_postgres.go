package postgre

import (
	_errors "common/errors"
	_models "common/models"
	"context"
	"crm/internal"
	"crm/internal/domain"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.FollowRepo
type FollowPostgre struct {
	DB *gorm.DB
}

func NewFollowPostgre(DB *gorm.DB) *FollowPostgre {
	return &FollowPostgre{
		DB: DB,
	}
}

func (r *FollowPostgre) GetFollow(ctx context.Context, followerID, followingID uint64) (*domain.FollowEntity, error) {
	var follow domain.FollowEntity
	err := r.DB.WithContext(ctx).
		Where("created_by = ? AND following_id = ? AND deleted_at IS NULL", followerID, followingID).
		First(&follow).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &follow, err
}

func (r *FollowPostgre) CreateFollow(ctx context.Context, followerID, followingID uint64) (*domain.FollowEntity, error) {
	if followerID == followingID {
		return nil, _errors.ReturnError(service.SelfFollowNotAllowed, _errors.WithPublicMessage("cannot follow yourself"))
	}
	follow := &domain.FollowEntity{
		FollowingID: followingID,
		BaseEntity: _models.BaseEntity{
			AuditBase: _models.AuditBase{
				CreatedBy: &followerID,
			},
		},
		Status: 1,
	}
	err := r.DB.WithContext(ctx).Create(follow).Error
	return follow, err
}

func (r *FollowPostgre) Update(ctx context.Context, follow *domain.FollowEntity) error {
	return r.DB.WithContext(ctx).Save(follow).Error
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
func (repo *FollowPostgre) FollowerUser(c context.Context, currentId uint64, page int, pageSize int) ([]domain.Profile, error) {
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
func (repo *FollowPostgre) FollowingUser(c context.Context, currentId uint64, page int, pageSize int) ([]uint64, error) {
	var friends []uint64

	query := repo.DB.Model(&domain.FollowEntity{}).
		Select("following_id").
		Where("created_by = ? AND deleted_at IS NULL", currentId).
		Pluck("following_id", &friends)
	// query := repo.DB.Model(&domain.Profile{}).
	// 	Select("profile_transfer.profile_id").
	// 	Joins("LEFT JOIN follow ON (follow.following_id = profile_transfer.profile_id) "+
	// 		"						AND profile_transfer.profile_id <> ? ", currentId).
	// 	Where("follow.deleted_at is null and follow.created_by = ?", currentId)

	// if dto.Text != nil && *dto.Text != "" {
	// 	query = query.Where("full_name LIKE ? OR username LIKE ?", "%"+*dto.Text+"%", "%"+*dto.Text+"%")
	// }

	err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&friends).Error
	return friends, err
}

func (repo *FollowPostgre) FollowUser(c context.Context, profileId, followingID uint64) (*domain.FollowEntity, error) {
	if profileId == followingID {
		return nil, _errors.ReturnError(service.SelfFollowNotAllowed)
	}

	follow := domain.FollowEntity{
		FollowingID: followingID,
	}

	var exists bool
	// Kiểm tra nếu đã follow trước đó
	err := repo.DB.Raw(`SELECT EXISTS (
							SELECT 1 FROM follow
							WHERE created_by = ? AND following_id = ? AND deleted_at IS NULL)`,
		profileId, followingID).Scan(&exists).Error

	if err != nil || exists {
		return nil, _errors.ReturnError(service.AlreadyFollowing)
	}

	err = repo.DB.WithContext(c).Create(&follow).Error

	return &follow, err
}

func (repo *FollowPostgre) UnfollowUser(c context.Context, profileId uint64, followingId uint64) error {
	// Tìm bản ghi follow
	result := repo.DB.Model(&domain.FollowEntity{}).
		Where("created_by = ? AND following_id = ?", profileId, followingId).
		Update("deleted_at", time.Now())

	if result.RowsAffected == 0 {
		return _errors.ReturnError(service.NotFollowing)
	}
	return nil
}

func (repo *FollowPostgre) BlockFollow(c context.Context, profileId uint64, targetId uint64) error {
	repo.DB.Model(&domain.FollowEntity{}).
		Where("(created_by = ? AND following_id = ?) OR (created_by = ? AND following_id = ?)",
			profileId, targetId,
			targetId, profileId,
		).
		Update("deleted_at", time.Now())

	return nil
}

func (repo *FollowPostgre) GetFollowInfoCount(c context.Context, userID uint64) (followerCount int64, followingCount int64, friendCount int64, err error) {
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
		Where("(created_by = ? or receiver_id = ?) and status = 20", userID, userID).
		Count(&friendCount).
		Error
	if err != nil {
		return 0, 0, 0, err
	}

	return followerCount, followingCount, friendCount, nil
}

func (repo *FollowPostgre) GetFollowing(ctx context.Context, profileId uint64, targetId uint64) (bool, error) {
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

func (repo *FollowPostgre) GetFollowId(ctx context.Context, profileId uint64, targetId uint64) (*uint64, error) {
	var followId uint64

	err := repo.DB.Model(&domain.FollowEntity{}).
		Select("id").
		Where("created_by = ? AND following_id = ? AND deleted_at IS NULL", profileId, targetId).
		First(&followId).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &followId, nil
}
