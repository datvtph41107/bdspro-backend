package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type OrganizationPermissionModel struct {
	Id   uint32 `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null;index"`
	Key  string `gorm:"not null;index:idx_org_permissions_key"`
	OrganizationId uint32 `gorm:"not null;index:idx_org_permissions_org_id"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (OrganizationPermissionModel) TableName() string {
	return "organization_permissions"
}

func OrganizationPermissionModelToEntity(model *OrganizationPermissionModel) *entity.OrganizationPermission {
	return &entity.OrganizationPermission{
		Id:   model.Id,
		Name: model.Name,
		Key:  model.Key,
		OrganizationId: model.OrganizationId,
	}
}

func OrganizationPermissionEntityToModel(entity *entity.OrganizationPermission) *OrganizationPermissionModel {
	return &OrganizationPermissionModel{
		Id:   entity.Id,
		Name: entity.Name,
		Key:  entity.Key,
		OrganizationId: entity.OrganizationId,
	}
}

func OrganizationPermissionModelsToEntities(models []*OrganizationPermissionModel) []*entity.OrganizationPermission {
	entities := make([]*entity.OrganizationPermission, len(models))
	for i, model := range models {
		entities[i] = OrganizationPermissionModelToEntity(model)
	}
	return entities
}

type organizationPermissionRepository struct {
	db *gorm.DB
}

func NewOrganizationPermissionPostgresRepository(db *gorm.DB) repository.OrganizationPermissionRepository {
	return &organizationPermissionRepository{db: db}
}
 
func (r *organizationPermissionRepository) Create(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {
	model := OrganizationPermissionEntityToModel(permission)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return OrganizationPermissionModelToEntity(model), nil
}



func (r *organizationPermissionRepository) FindById(ctx context.Context, id uint32) (*entity.OrganizationPermission, error) {
	var model OrganizationPermissionModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return OrganizationPermissionModelToEntity(&model), nil
}

func (r *organizationPermissionRepository) Update(ctx context.Context, permission *entity.OrganizationPermission) (*entity.OrganizationPermission, error) {	
	model := OrganizationPermissionEntityToModel(permission)
	if err := r.db.WithContext(ctx).Model(&OrganizationPermissionModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return OrganizationPermissionModelToEntity(model), nil
}

func (r *organizationPermissionRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&OrganizationPermissionModel{}, id).Error
}

func (r *organizationPermissionRepository) FindByKeyAndOrganizationId(ctx context.Context, key string, organizationId uint32) (*entity.OrganizationPermission, error) {
	var model OrganizationPermissionModel
	if err := r.db.WithContext(ctx).Where("key = ? AND organization_id = ?", key, organizationId).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return OrganizationPermissionModelToEntity(&model), nil
}

func (r *organizationPermissionRepository) FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationPermission, error) {
	var models []*OrganizationPermissionModel
	if err := r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Find(&models).Error; err != nil {
		return nil, err
	}
	return OrganizationPermissionModelsToEntities(models), nil
}

func (r *organizationPermissionRepository) FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationPermission, uint32, error) {
	var permissions []*entity.OrganizationPermission
	var total int64

	// Get total count
	err := r.db.Model(&OrganizationPermissionModel{}).Where("organization_id = ?", organizationId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data
	offset := (page - 1) * size
	err = r.db.Where("organization_id = ?", organizationId).
		Offset(offset).
		Limit(size).
		Find(&permissions).Error
	if err != nil {
		return nil, 0, err
	}

	return permissions, uint32(total), nil
}

func (r *organizationPermissionRepository) FindByOrganizationIdAndKey(ctx context.Context, organizationId uint32, key string) (*entity.OrganizationPermission, error) {
	var model OrganizationPermissionModel
	err := r.db.Where("organization_id = ? AND key = ?", organizationId, key).First(&model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationPermissionModelToEntity(&model), nil
}

func (r *organizationPermissionRepository) GetAll(ctx context.Context) ([]*entity.OrganizationPermission, error) {
	var models []*OrganizationPermissionModel
	err := r.db.Find(&models).Error
	if err != nil {
		return nil, err
	}
	return OrganizationPermissionModelsToEntities(models), nil
}

func (r *organizationPermissionRepository) GetByKey(ctx context.Context, key string) (*entity.OrganizationPermission, error) {
	var model OrganizationPermissionModel
	err := r.db.Where("key = ?", key).First(&model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationPermissionModelToEntity(&model), nil
}