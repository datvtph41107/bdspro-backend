package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type OrganizationRoleModel struct {
	ID            uint32    `gorm:"primaryKey;autoIncrement"`
	Name          string    `gorm:"not null"`
	Key           string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	Permissions   []OrganizationPermissionModel `gorm:"many2many:role_permissions;"`
	OrganizationId uint32 `gorm:"not null"`
}

func (OrganizationRoleModel) TableName() string {
	return "organization_roles"
}

func RoleModelToEntity(model *OrganizationRoleModel) *entity.OrganizationRole {
	permissions := make([]*entity.OrganizationPermission, len(model.Permissions))
	for i, permission := range model.Permissions {
		permissions[i] = OrganizationPermissionModelToEntity(&permission)
	}
	return &entity.OrganizationRole{
		Id:            model.ID,
		Name:          model.Name,
		Key:           model.Key,
		Permissions:   permissions,
		OrganizationId: model.OrganizationId,
	}
}

func RoleEntityToModel(entity *entity.OrganizationRole) *OrganizationRoleModel {
	return &OrganizationRoleModel{
		ID:            entity.Id,
		Name:          entity.Name,
		Key:           entity.Key,
		OrganizationId: entity.OrganizationId,
	}
}

func RoleModelsToEntities(models []*OrganizationRoleModel) []*entity.OrganizationRole {
	entities := make([]*entity.OrganizationRole, len(models))
	for i, model := range models {
		entities[i] = RoleModelToEntity(model)
	}
	return entities
}

type organizationRoleRepository struct {
	db *gorm.DB
}

func NewOrganizationRolePostgresRepository(db *gorm.DB) repository.OrganizationRoleRepository {
	return &organizationRoleRepository{db: db}
}

func (r *organizationRoleRepository) Create(ctx context.Context, role *entity.OrganizationRole) (*entity.OrganizationRole, error) {
	model := &OrganizationRoleModel{
		Name: role.Name,
		Key:  role.Key,
		OrganizationId: role.OrganizationId,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	return RoleModelToEntity(model), nil
}

func (r *organizationRoleRepository) GetAll(ctx context.Context) ([]*entity.OrganizationRole, error) {
	var models []*OrganizationRoleModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	return RoleModelsToEntities(models), nil
}

func (r *organizationRoleRepository) FindById(ctx context.Context, id uint32) (*entity.OrganizationRole, error) {
	var model OrganizationRoleModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return RoleModelToEntity(&model), nil
}

func (r *organizationRoleRepository) Update(ctx context.Context, role *entity.OrganizationRole) (*entity.OrganizationRole, error) {
	model := RoleEntityToModel(role)

	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, err
	}

	return RoleModelToEntity(model), nil
}

func (r *organizationRoleRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&OrganizationRoleModel{}, id).Error
}

func (r *organizationRoleRepository) GetByKey(ctx context.Context, key string) (*entity.OrganizationRole, error) {
	var model OrganizationRoleModel
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return RoleModelToEntity(&model), nil
}

func (r *organizationRoleRepository) FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationRole, error) {
	var models []*OrganizationRoleModel
	if err := r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Find(&models).Error; err != nil {
		return nil, err
	}
	return RoleModelsToEntities(models), nil
}

func (r *organizationRoleRepository) FindByOrganizationIdAndKey(ctx context.Context, organizationId uint32, key string) (*entity.OrganizationRole, error) {
	var model OrganizationRoleModel
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND key = ?", organizationId, key).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return RoleModelToEntity(&model), nil
}
