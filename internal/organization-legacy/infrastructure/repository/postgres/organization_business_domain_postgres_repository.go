package postgres

import (
	"context"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type OrganizationBusinessDomainModel struct {
	ID               uint32    `gorm:"primaryKey;autoIncrement"`
	OrganizationID   uint32    `gorm:"not null;index:idx_org_business_domain_org_id"`
	BusinessDomainID uint32    `gorm:"not null;index:idx_org_business_domain_domain_id"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	CreatedBy        uint32    `gorm:"index"`
}

func (OrganizationBusinessDomainModel) TableName() string {
	return "organization_business_domains"
}

func OrganizationBusinessDomainModelToEntity(model *OrganizationBusinessDomainModel) *entity.OrganizationBusinessDomain {
	return &entity.OrganizationBusinessDomain{
		ID:               model.ID,
		OrganizationID:   model.OrganizationID,
		BusinessDomainID: model.BusinessDomainID,
		CreatedAt:        model.CreatedAt,
		CreatedBy:        model.CreatedBy,
	}
}

func OrganizationBusinessDomainEntityToModel(entity *entity.OrganizationBusinessDomain) *OrganizationBusinessDomainModel {
	return &OrganizationBusinessDomainModel{
		ID:               entity.ID,
		OrganizationID:   entity.OrganizationID,
		BusinessDomainID: entity.BusinessDomainID,
		CreatedAt:        entity.CreatedAt,
		CreatedBy:        entity.CreatedBy,
	}
}

func OrganizationBusinessDomainModelsToEntities(models []*OrganizationBusinessDomainModel) []*entity.OrganizationBusinessDomain {
	entities := make([]*entity.OrganizationBusinessDomain, len(models))
	for i, model := range models {
		entities[i] = OrganizationBusinessDomainModelToEntity(model)
	}
	return entities
}

// @bind: organization/internal/domain/repository.OrganizationBusinessDomainRepository
type OrganizationBusinessDomainPostgresRepository struct {
	db *gorm.DB
}

func NewOrganizationBusinessDomainPostgresRepository(db *gorm.DB) *OrganizationBusinessDomainPostgresRepository {
	return &OrganizationBusinessDomainPostgresRepository{db: db}
}

func (r *OrganizationBusinessDomainPostgresRepository) Create(ctx context.Context, organizationBusinessDomain *entity.OrganizationBusinessDomain) (*entity.OrganizationBusinessDomain, error) {
	model := OrganizationBusinessDomainEntityToModel(organizationBusinessDomain)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return OrganizationBusinessDomainModelToEntity(model), nil
}

func (r *OrganizationBusinessDomainPostgresRepository) FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationBusinessDomain, error) {
	var models []*OrganizationBusinessDomainModel
	if err := r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Find(&models).Error; err != nil {
		return nil, err
	}
	return OrganizationBusinessDomainModelsToEntities(models), nil
}

func (r *OrganizationBusinessDomainPostgresRepository) FindByBusinessDomainId(ctx context.Context, businessDomainId uint32) ([]*entity.OrganizationBusinessDomain, error) {
	var models []*OrganizationBusinessDomainModel
	if err := r.db.WithContext(ctx).Where("business_domain_id = ?", businessDomainId).Find(&models).Error; err != nil {
		return nil, err
	}
	return OrganizationBusinessDomainModelsToEntities(models), nil
}

func (r *OrganizationBusinessDomainPostgresRepository) DeleteByOrganizationId(ctx context.Context, organizationId uint32) error {
	return r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Delete(&OrganizationBusinessDomainModel{}).Error
}

func (r *OrganizationBusinessDomainPostgresRepository) DeleteByOrganizationIdAndBusinessDomainId(ctx context.Context, organizationId uint32, businessDomainId uint32) error {
	return r.db.WithContext(ctx).Where("organization_id = ? AND business_domain_id = ?", organizationId, businessDomainId).Delete(&OrganizationBusinessDomainModel{}).Error
}

func (r *OrganizationBusinessDomainPostgresRepository) FindByOrganizationIds(ctx context.Context, organizationIds []uint32) ([]*entity.OrganizationBusinessDomain, error) {
	var models []*OrganizationBusinessDomainModel
	if err := r.db.WithContext(ctx).Where("organization_id IN ?", organizationIds).Find(&models).Error; err != nil {
		return nil, err
	}
	return OrganizationBusinessDomainModelsToEntities(models), nil
}
