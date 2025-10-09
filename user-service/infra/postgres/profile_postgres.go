package postgres

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"user/internal/dto"
	models "user/internal/models"

	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.IProfileRepo
type ProfilePostgres struct {
	DB *gorm.DB
}

func NewProfilePostgres(db *gorm.DB) *ProfilePostgres {
	return &ProfilePostgres{
		DB: db,
	}
}

func (r *ProfilePostgres) GlobalSearchProfile(c context.Context, text string, pageable _dto.Pagable) (*[]models.UserProfileSearch, int64, error) {
	var profiles []models.UserProfileSearch
	var total int64
	currentUserId := _utils.GetProfileIdWithContext(c)

	query := r.DB.WithContext(c).
		Model(&models.UserProfileEntity{}).
		Select("profile_id, full_name, avatar, email").
		Where("deleted_at is null AND profile_id != ?", currentUserId)

	if text != "" {
		query = query.Where("LOWER(full_name) LIKE ? OR email LIKE ?", "%"+text+"%", "%"+text+"%")
	}

	err := query.Debug().
		Offset(pageable.GetOffset()).
		Limit(pageable.GetLimit()).
		Order("updated_at DESC").
		Find(&profiles).
		Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return &profiles, total, nil
}

func (r *ProfilePostgres) FriendSearchProfile(c context.Context, profileId uint64, dto dto.ProfileSearch) (*[]models.UserProfileSearch, int64, error) {
	var friends []models.UserProfileSearch
	var total int64

	query := r.DB.
		Model(&models.FriendEntity{}).
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

func (r *ProfilePostgres) UpdateProfile(c context.Context, profileId uint64, dto *dto.UserInfoRequest) error {
	// Chỉ cập nhật các trường cụ thể
	return r.DB.WithContext(c).
		Model(&models.UserProfileEntity{}).
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
func (r *ProfilePostgres) GetByProfileID(profileId uint64) (*models.UserProfileEntity, error) {
	var userProfile models.UserProfileEntity
	if err := r.DB.Where("profile_id = ? and deleted_at is null", profileId).First(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func (r *ProfilePostgres) GetProfileByIds(ids []uint64) ([]*models.UserProfileEntity, error) {
	var userProfiles []*models.UserProfileEntity
	if err := r.DB.Where("profile_id IN (?) and deleted_at is null", ids).Find(&userProfiles).Error; err != nil {
		return nil, err
	}
	return userProfiles, nil
}

func (r *ProfilePostgres) GetProfileByPhones(phones []string) ([]*models.UserProfileEntity, error) {
	var userProfiles []*models.UserProfileEntity
	if err := r.DB.Where("phone IN (?) and deleted_at is null", phones).Find(&userProfiles).Error; err != nil {
		return nil, err
	}
	return userProfiles, nil
}
