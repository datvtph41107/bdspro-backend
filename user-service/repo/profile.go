package repo

import (
	"user/data"
	user_model "user/models/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProfileRepo struct {
	DB *gorm.DB
}

func NewProfileRepo(db *gorm.DB) *ProfileRepo {
	return &ProfileRepo{
		DB: db,
	}
}

func (r *ProfileRepo) GlobalSearchProfile() {

}

func (r *ProfileRepo) FriendSearchProfile(c *gin.Context, profileId uint64, dto data.ProfileSearch) (*[]user_model.UserProfileSearch, int64, error) {
	var friends []user_model.UserProfileSearch
	var total int64

	query := r.DB.
		Model(&user_model.FriendEntity{}).
		Select("user_profile.profile_id, user_profile.full_name, user_profile.avatar, user_profile.email").
		Preload("GroupUser").
		Preload("GroupReceiver").
		Joins("LEFT JOIN user_profile ON (friend.created_by = user_profile.profile_id OR friend.receiver_id = user_profile.profile_id) "+
			"						AND user_profile.profile_id <> ? ", profileId).
		Where("friend.status = 2 and friend.deleted_at is null")

	if dto.Text != "" {
		query = query.Where("LOWER(user_profile.full_name) LIKE ? OR user_profile.email LIKE ?", "%"+dto.Text+"%", "%"+dto.Text+"%")
	}

	err := query.Debug().Offset(dto.GetOffset()).Limit(dto.GetLimit()).Scan(&friends).Error
	if err != nil {
		return nil, 0, err
	}
	err = query.Count(&total).Error
	return &friends, total, err
}
