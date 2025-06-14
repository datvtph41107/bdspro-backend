package repo

import (
	"context"
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

func (r *ProfileRepo) FriendSearchProfile(c context.Context, profileId uint64, dto data.ProfileSearch) (*[]user_model.UserProfileSearch, int64, error) {
	var friends []user_model.UserProfileSearch
	var total int64

	query := r.DB.
		Model(&user_model.FriendEntity{}).
		Select("user_profile.profile_id, user_profile.full_name, user_profile.avatar, user_profile.email").
		Preload("GroupUser").
		Preload("GroupReceiver").
		Joins("LEFT JOIN user_profile ON (friend.created_by = user_profile.profile_id OR friend.receiver_id = user_profile.profile_id) "+
			"						AND user_profile.profile_id <> ? ", profileId).
		Where("friend.status = 20 and friend.deleted_at is null")

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

func (r *ProfileRepo) UpdateProfile(c *gin.Context, profileId uint64, dto *data.UserInfoRequest) error {
	// Chỉ cập nhật các trường cụ thể
	return r.DB.WithContext(c).
		Model(&user_model.UserProfileEntity{}).
		Where("profile_id = ?", profileId).
		Updates(map[string]interface{}{
			"full_name":          dto.FullName,
			"email":              dto.Email,
			"avatar":             dto.Avatar,
			"background_image":   dto.BackgroundImage,
			"phone":              dto.Phone,
			"birth":              dto.Birth,
			"start_date":         dto.StartDate,
			"address":            dto.Address,
			"phone2":             dto.Phone2,
			"tax_code":           dto.TaxCode,
			"position":           dto.Position,
			"department_id":      dto.DepartmentID,
			"slogan":             dto.Slogan,
			"website":            dto.Website,
			"facebook":           dto.Facebook,
			"instagram":          dto.Instagram,
			"twitter":            dto.Twitter,
			"linkedin":           dto.Linkedin,
			"youtube":            dto.Youtube,
			"introduction":       dto.Introduction,
			"visibility":         dto.Visibility,
			"profile_visibility": dto.ProfileVisibility,
			"status_online":      dto.StatusOnline,
			"front_identify":     dto.FrontIdentify,
			"back_identify":      dto.BackIdentify,
			"tab_default":        dto.TabDefault,
			"role_real_estate":   dto.RoleRealEstate,
		}).Error
}

// GetByProfileId truy vấn UserProfileEntity theo profileId
func (r *ProfileRepo) GetByProfileID(profileId uint64) (*user_model.UserProfileEntity, error) {
	var userProfile user_model.UserProfileEntity
	if err := r.DB.Where("profile_id = ? and deleted_at is null", profileId).First(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func (r *ProfileRepo) GetProfileByIds(ids []uint64) ([]*user_model.UserProfileEntity, error) {
	var userProfiles []*user_model.UserProfileEntity
	if err := r.DB.Where("profile_id IN (?) and deleted_at is null", ids).Find(&userProfiles).Error; err != nil {
		return nil, err
	}
	return userProfiles, nil
}

func (r *ProfileRepo) GetProfileByPhones(phones []string) ([]*user_model.UserProfileEntity, error) {
	var userProfiles []*user_model.UserProfileEntity
	if err := r.DB.Where("phone IN (?) and deleted_at is null", phones).Find(&userProfiles).Error; err != nil {
		return nil, err
	}
	return userProfiles, nil
}
