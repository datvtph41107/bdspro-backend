package postgre

import (
	_routes "common/routes"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.ContactRepo
type PostgreContact struct {
	DB *gorm.DB
}

func NewContactRepo(DB *gorm.DB) *PostgreContact {
	return &PostgreContact{
		DB: DB,
	}
}

func (r *PostgreContact) Search(c context.Context, profileId uint64, dto dto.ContactSearchDTO) ([]domain.ContactEntity, int64, error) {
	var contacts []domain.ContactEntity
	result := r.DB.Model(&domain.ContactEntity{}).
		Where("deleted_at IS NULL AND owner_id = ?", profileId).
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&contacts)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	var total int64
	r.DB.Model(&domain.ContactEntity{}).Where("deleted_at IS NULL AND owner_id = ?", profileId).Count(&total)
	return contacts, total, nil
}

func (r *PostgreContact) Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	result := r.DB.Create(entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (r *PostgreContact) Update(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	result := r.DB.Save(entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (r *PostgreContact) Delete(c context.Context, id uint64) error {
	result := r.DB.Delete(&domain.ContactEntity{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PostgreContact) GetByID(c context.Context, id uint64) (*domain.ContactEntity, error) {
	var contact domain.ContactEntity
	result := r.DB.Model(&domain.ContactEntity{}).Where("id = ?", id).First(&contact)
	if result.Error != nil {
		return nil, result.Error
	}
	return &contact, nil
}

// ExistsBy kiểm tra xem một ContactEntity có tồn tại không dựa trên điều kiện
func (r *PostgreContact) ExistByProfile(profileId uint64) (bool, error) {
	var count int64
	result := r.DB.Model(&domain.Profile{}).Where("profile_id = ?", profileId).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

func (r *PostgreContact) Existed(followId uint64) error {
	var count int64
	result := r.DB.Model(&domain.Profile{}).Where("profile_id = ?", followId).Count(&count)
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
