package postgres

import (
	"context"
	"errors"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type BusinessDomainModel struct {
	ID          uint32 `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null;index"`
	Description *string
	Code        string    `gorm:"not null;uniqueIndex"`
	IsActive    bool      `gorm:"default:true;index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	CreatedBy   uint32    `gorm:"index"`
	UpdatedBy   uint32
	DeletedAt   *time.Time `gorm:"index"`
	IsDeleted   bool       `gorm:"default:false"`
}

func (BusinessDomainModel) TableName() string {
	return "business_domains"
}

func BusinessDomainModelToEntity(model *BusinessDomainModel) *entity.BusinessDomain {
	return &entity.BusinessDomain{
		ID:          model.ID,
		Name:        model.Name,
		Description: model.Description,
		Code:        model.Code,
		IsActive:    model.IsActive,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		CreatedBy:   model.CreatedBy,
		UpdatedBy:   model.UpdatedBy,
	}
}

func BusinessDomainEntityToModel(entity *entity.BusinessDomain) *BusinessDomainModel {
	return &BusinessDomainModel{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Code:        entity.Code,
		IsActive:    entity.IsActive,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
		IsDeleted:   false,
	}
}

func BusinessDomainModelsToEntities(models []*BusinessDomainModel) []*entity.BusinessDomain {
	entities := make([]*entity.BusinessDomain, len(models))
	for i, model := range models {
		entities[i] = BusinessDomainModelToEntity(model)
	}
	return entities
}

// @bind: organization/internal/domain/repository.BusinessDomainRepository
type BusinessDomainPostgresRepository struct {
	db *gorm.DB
}

func NewBusinessDomainPostgresRepository(db *gorm.DB) *BusinessDomainPostgresRepository {
	return &BusinessDomainPostgresRepository{db: db}
}

func (r *BusinessDomainPostgresRepository) Create(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error) {
	model := BusinessDomainEntityToModel(businessDomain)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return BusinessDomainModelToEntity(model), nil
}

func (r *BusinessDomainPostgresRepository) FindById(ctx context.Context, id uint32) (*entity.BusinessDomain, error) {
	var model BusinessDomainModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return BusinessDomainModelToEntity(&model), nil
}

func (r *BusinessDomainPostgresRepository) FindByCode(ctx context.Context, code string) (*entity.BusinessDomain, error) {
	var model BusinessDomainModel
	if err := r.db.WithContext(ctx).Where("code = ? AND is_deleted = ?", code, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return BusinessDomainModelToEntity(&model), nil
}

func (r *BusinessDomainPostgresRepository) FindAll(ctx context.Context) ([]*entity.BusinessDomain, error) {
	var models []*BusinessDomainModel
	if err := r.db.WithContext(ctx).Where("is_deleted = ?", false).Find(&models).Error; err != nil {
		return nil, err
	}
	return BusinessDomainModelsToEntities(models), nil
}

func (r *BusinessDomainPostgresRepository) FindActive(ctx context.Context) ([]*entity.BusinessDomain, error) {
	var models []*BusinessDomainModel
	if err := r.db.WithContext(ctx).Where("is_active = ? AND is_deleted = ?", true, false).Find(&models).Error; err != nil {
		return nil, err
	}
	return BusinessDomainModelsToEntities(models), nil
}

func (r *BusinessDomainPostgresRepository) Update(ctx context.Context, businessDomain *entity.BusinessDomain) (*entity.BusinessDomain, error) {
	model := BusinessDomainEntityToModel(businessDomain)
	if err := r.db.WithContext(ctx).Model(&model).Where("id = ? AND is_deleted = ?", model.ID, false).Updates(model).Error; err != nil {
		return nil, err
	}
	return BusinessDomainModelToEntity(model), nil
}

func (r *BusinessDomainPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&BusinessDomainModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *BusinessDomainPostgresRepository) FindByIds(ctx context.Context, ids []uint32) ([]*entity.BusinessDomain, error) {
	var models []*BusinessDomainModel
	if err := r.db.WithContext(ctx).Where("id IN ? AND is_deleted = ?", ids, false).Find(&models).Error; err != nil {
		return nil, err
	}
	return BusinessDomainModelsToEntities(models), nil
}
