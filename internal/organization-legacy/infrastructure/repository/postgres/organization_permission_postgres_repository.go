package postgres

import (
	"context"
	"errors"
	"time"

	"organization/internal/domain/entity"
	"organization/internal/enums"

	"gorm.io/gorm"
)

type OrganizationPermissionModel struct {
	Id          uint32 `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null;index"`
	Key         string `gorm:"not null;index:idx_org_permissions_key"`
	Description string `gorm:"type:text"`
	// OrganizationId uint32               `gorm:"not null;index:idx_org_permissions_org_id"`
	ParentId       uint32               `json:"parent_id"`
	PermissionType enums.PermissionType `gorm:"not null;default:10" json:"permission_type"`
	CreatedAt      time.Time            `gorm:"autoCreateTime"`
	UpdatedAt      time.Time            `gorm:"autoUpdateTime"`
	DeletedAt      *time.Time           `gorm:"index"`
	IsDeleted      bool                 `gorm:"default:false"`
}

func (OrganizationPermissionModel) TableName() string {
	return "organization_permissions"
}

func OrganizationPermissionModelToEntity(model *OrganizationPermissionModel) *entity.OrganizationPermission {
	return &entity.OrganizationPermission{
		Id:          model.Id,
		Name:        model.Name,
		Key:         model.Key,
		Description: model.Description,
		// OrganizationId: model.OrganizationId,
		PermissionType: model.PermissionType,
		ParentId:       model.ParentId,
	}
}

func OrganizationPermissionEntityToModel(entity *entity.OrganizationPermission) *OrganizationPermissionModel {
	return &OrganizationPermissionModel{
		Id:          entity.Id,
		Name:        entity.Name,
		Key:         entity.Key,
		Description: entity.Description,
		// OrganizationId: entity.OrganizationId,
		PermissionType: entity.PermissionType,
		IsDeleted:      false,
	}
}

func OrganizationPermissionModelsToEntities(models []*OrganizationPermissionModel) []*entity.OrganizationPermission {
	entities := make([]*entity.OrganizationPermission, len(models))
	for i, model := range models {
		entities[i] = OrganizationPermissionModelToEntity(model)
	}
	return entities
}

// @bind: organization/internal/domain/repository.OrganizationPermissionRepository
type OrganizationPermissionRepository struct {
	db *gorm.DB
}

func NewOrganizationPermissionPostgresRepository(db *gorm.DB) *OrganizationPermissionRepository {
	return &OrganizationPermissionRepository{db: db}
}

func (r *OrganizationPermissionRepository) Create(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {
	model := OrganizationPermissionEntityToModel(permission)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return OrganizationPermissionModelToEntity(model), nil
}

func (r *OrganizationPermissionRepository) GetPermissionOrganization(ctx context.Context) ([]*entity.OrganizationPermission, error) {
	var models []*OrganizationPermissionModel
	if err := r.db.WithContext(ctx).Where("permission_type = ? AND is_deleted = ?", enums.PermissionTypeOrganization, false).Find(&models).Error; err != nil {
		return nil, err
	}
	return OrganizationPermissionModelsToEntities(models), nil
}

func (r *OrganizationPermissionRepository) FindById(ctx context.Context, id uint32) (*entity.OrganizationPermission, error) {
	var model OrganizationPermissionModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return OrganizationPermissionModelToEntity(&model), nil
}

func (r *OrganizationPermissionRepository) Update(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {
	model := OrganizationPermissionEntityToModel(permission)
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", model.Id, false).
		Updates(map[string]interface{}{
			"name":        model.Name,
			"key":         model.Key,
			"description": model.Description,
		}).Error; err != nil {
		return nil, err
	}
	return OrganizationPermissionModelToEntity(model), nil
}

func (r *OrganizationPermissionRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&OrganizationPermissionModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *OrganizationPermissionRepository) FindByKeyAndOrganizationId(ctx context.Context, key string, organizationId uint32) (*entity.OrganizationPermission, error) {
	var model OrganizationPermissionModel
	if err := r.db.WithContext(ctx).Where("key = ? AND organization_id = ? AND is_deleted = ?", key, organizationId, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return OrganizationPermissionModelToEntity(&model), nil
}

func (r *OrganizationPermissionRepository) FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationPermission, error) {
	var models []*OrganizationPermissionModel
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND is_deleted = ?", organizationId, false).Find(&models).Error; err != nil {
		return nil, err
	}
	return OrganizationPermissionModelsToEntities(models), nil
}

func (r *OrganizationPermissionRepository) FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationPermission, uint32, error) {
	var permissions []*entity.OrganizationPermission
	var total int64

	// Get total count
	err := r.db.Model(&OrganizationPermissionModel{}).Where("organization_id = ? AND is_deleted = ?", organizationId, false).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * size
	err = r.db.Where("organization_id = ? AND is_deleted = ?", organizationId, false).
		Offset(offset).
		Limit(size).
		Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, uint32(total), nil
}

func (r *OrganizationPermissionRepository) FindByOrganizationIdAndKey(ctx context.Context, organizationId uint32, key string) (*entity.OrganizationPermission, error) {
	var model OrganizationPermissionModel
	err := r.db.Where("organization_id = ? AND key = ?", organizationId, key).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return OrganizationPermissionModelToEntity(&model), nil
}

func (r *OrganizationPermissionRepository) GetByKey(ctx context.Context, key string) (*entity.OrganizationPermission, error) {
	var model OrganizationPermissionModel
	err := r.db.Where("key = ?", key).First(&model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationPermissionModelToEntity(&model), nil
}

func (r *OrganizationPermissionRepository) DeletePermission(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&OrganizationPermissionModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}
