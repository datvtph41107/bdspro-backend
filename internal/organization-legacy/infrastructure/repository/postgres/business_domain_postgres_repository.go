package postgres

import (
	"context"
	"errors"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

// type entity.BusinessDomain struct {
// 	ID          uint32 `gorm:"primaryKey;autoIncrement"`
// 	Name        string `gorm:"not null;index"`
// 	Description *string
// 	Code        string    `gorm:"not null;uniqueIndex"`
// 	IsActive    bool      `gorm:"default:true;index"`
// 	CreatedAt   time.Time `gorm:"autoCreateTime"`
// 	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
// 	CreatedBy   uint32    `gorm:"index"`
// 	UpdatedBy   uint32
// 	DeletedAt   *time.Time `gorm:"index"`
// 	IsDeleted   bool       `gorm:"default:false"`
// }

// func (entity.BusinessDomain) TableName() string {
// 	return "business_domains"
// }

// @bind: organization/internal/domain/repository.BusinessDomainRepository
type BusinessDomainPostgresRepository struct {
	db *gorm.DB
}

func NewBusinessDomainPostgresRepository(db *gorm.DB) *BusinessDomainPostgresRepository {
	return &BusinessDomainPostgresRepository{db: db}
}

func (r *BusinessDomainPostgresRepository) Create(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error) {
	if err := r.db.WithContext(ctx).Create(businessDomain).Error; err != nil {
		return nil, err
	}
	return businessDomain, nil
}

func (r *BusinessDomainPostgresRepository) FindById(ctx context.Context, id uint32) (*entity.BusinessDomain, error) {
	var model entity.BusinessDomain
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

func (r *BusinessDomainPostgresRepository) FindByCode(ctx context.Context, code string) (*entity.BusinessDomain, error) {
	var model entity.BusinessDomain
	if err := r.db.WithContext(ctx).Where("code = ? AND is_deleted = ?", code, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

func (r *BusinessDomainPostgresRepository) FindAll(ctx context.Context) ([]*entity.BusinessDomain, error) {
	var models []*entity.BusinessDomain
	if err := r.db.WithContext(ctx).Where("is_deleted = ?", false).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *BusinessDomainPostgresRepository) FindActive(ctx context.Context) ([]*entity.BusinessDomain, error) {
	var models []*entity.BusinessDomain
	if err := r.db.WithContext(ctx).Where("is_active = ? AND is_deleted = ?", true, false).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *BusinessDomainPostgresRepository) Update(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error) {
	model := businessDomain
	if err := r.db.WithContext(ctx).Model(&model).Where("id = ? AND is_deleted = ?", model.ID, false).Updates(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (r *BusinessDomainPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&entity.BusinessDomain{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *BusinessDomainPostgresRepository) FindByIds(ctx context.Context, ids []uint32) ([]*entity.BusinessDomain, error) {
	var models []*entity.BusinessDomain
	if err := r.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", ids, false).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}
