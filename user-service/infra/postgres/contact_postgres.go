package postgres

import (
	_errors "common/errors"
	"user/internal"
	models "user/internal/models"

	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.IContactRepo
type ContactPostgres struct {
	DB *gorm.DB
}

func NewContactPostgres(DB *gorm.DB) *ContactPostgres {
	return &ContactPostgres{
		DB: DB,
	}
}

// ExistsBy kiểm tra xem một ContactEntity có tồn tại không dựa trên điều kiện
func (r *ContactPostgres) ExistByProfile(profileId uint64) (bool, error) {
	var count int64
	result := r.DB.Model(&models.Profile{}).Where("profile_id = ?", profileId).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *ContactPostgres) Existed(followId uint64) error {
	var count int64
	result := r.DB.Model(&models.Profile{}).Where("profile_id = ?", followId).Count(&count)
	if result.Error != nil {
		return result.Error
	}

	// ok, err := r.contactRepo.ExistByProfile(followId)
	if count == 0 {
		return _errors.ReturnError(service.UserNotFoundProfile)
	}

	// if err != nil {
	// 	return &_routes.Except{
	// 		Code:    500,
	// 		Message: err.Error(),
	// 	}
	// }
	return nil
}
