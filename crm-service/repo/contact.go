package repo

import (
	_routes "common/routes"
	"crm/models"

	"gorm.io/gorm"
)

type ContactRepo struct {
	DB *gorm.DB
}

func NewContactRepo(DB *gorm.DB) *ContactRepo {
	return &ContactRepo{
		DB: DB,
	}
}

// ExistsBy kiểm tra xem một ContactEntity có tồn tại không dựa trên điều kiện
func (r *ContactRepo) ExistByProfile(profileId uint64) (bool, error) {
	var count int64
	result := r.DB.Model(&models.Profile{}).Where("profile_id = ?", profileId).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *ContactRepo) Existed(followId uint64) error {
	var count int64
	result := r.DB.Model(&models.Profile{}).Where("profile_id = ?", followId).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	// ok, err := r.contactRepo.ExistByProfile(followId)
	if count == 0 {
		return &_routes.Except{
			Code:    400,
			Message: "Người dùng không tồn tại",
		}
	}

	// if err != nil {
	// 	return &_routes.Except{
	// 		Code:    500,
	// 		Message: err.Error(),
	// 	}
	// }
	return nil
}
