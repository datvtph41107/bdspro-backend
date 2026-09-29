package postgres

import (
	"context"
	"errors"
	"time"

	"organization/internal/domain/entity"
	"organization/internal/enums"

	"gorm.io/gorm"
)

type OrganizationRoleModel struct {
	ID             uint32                         `gorm:"primaryKey;autoIncrement"`
	Name           string                         `gorm:"not null;index"`
	Key            string                         `gorm:"not null;index:idx_org_roles_key"`
	CreatedAt      time.Time                      `gorm:"autoCreateTime"`
	UpdatedAt      time.Time                      `gorm:"autoUpdateTime"`
	Permissions    []*OrganizationPermissionModel `gorm:"many2many:role_permissions;"`
	OrganizationId uint32                         `gorm:"not null;index:idx_org_roles_organization_id"`
	IsDefault      bool                           `gorm:"default:false;index"`
	DomainType     int                            `gorm:"default:0;index"`
	DeletedAt      *time.Time                     `gorm:"index"`
	IsDeleted      bool                           `gorm:"default:false"`
	Color          *entity.Color                  `gorm:"foreignKey:ColorId"`
	ColorId        *uint32                        ``
	RoleKey        uint32                         `gorm:"default:0;index"`
}

func (OrganizationRoleModel) TableName() string {
	return "organization_roles"
}

func RoleModelToEntity(model *OrganizationRoleModel) *entity.OrganizationRole {
	permissions := make([]*entity.OrganizationPermission, len(model.Permissions))
	for i, permission := range model.Permissions {
		permissions[i] = OrganizationPermissionModelToEntity(permission)
	}
	return &entity.OrganizationRole{
		Id:             model.ID,
		Name:           model.Name,
		Key:            model.Key,
		Permissions:    permissions,
		OrganizationId: model.OrganizationId,
		IsDefault:      model.IsDefault,
		DomainType:     enums.DomainType(model.DomainType),
		ColorId:        model.ColorId,
		Color:          model.Color,
		RoleKey:        model.RoleKey,
	}
}

func RoleEntityToModel(entity *entity.OrganizationRole) *OrganizationRoleModel {
	permissions := make([]*OrganizationPermissionModel, len(entity.PermissionIds))
	for i, permission := range entity.PermissionIds {
		permissions[i] = &OrganizationPermissionModel{
			Id: uint32(permission),
		}
	}
	return &OrganizationRoleModel{
		ID:             entity.Id,
		Name:           entity.Name,
		Key:            entity.Key,
		OrganizationId: entity.OrganizationId,
		Permissions:    permissions,
		IsDefault:      entity.IsDefault,
		DomainType:     int(entity.DomainType),
		IsDeleted:      false,
	}
}

func RoleModelsToEntities(models []*OrganizationRoleModel) []*entity.OrganizationRole {
	entities := make([]*entity.OrganizationRole, len(models))
	for i, model := range models {
		entities[i] = RoleModelToEntity(model)
	}
	return entities
}

// @bind: organization/internal/domain/repository.OrganizationRoleRepository
type OrganizationRoleRepository struct {
	db *gorm.DB
}

func NewOrganizationRolePostgresRepository(db *gorm.DB) *OrganizationRoleRepository {
	return &OrganizationRoleRepository{db: db}
}

func (r *OrganizationRoleRepository) Create(ctx context.Context, role *entity.OrganizationRole) (*entity.OrganizationRole, error) {
	model := &OrganizationRoleModel{
		Name:           role.Name,
		Key:            role.Key,
		OrganizationId: role.OrganizationId,
		IsDefault:      role.IsDefault,
		DomainType:     int(role.DomainType),
		IsDeleted:      false,
		// Permissions:    role.Permissions,
	}

	if err := GetDB(ctx, r.db).Create(model).Error; err != nil {
		return nil, err
	}

	return RoleModelToEntity(model), nil
}

func (r *OrganizationRoleRepository) GetAll(ctx context.Context) ([]*entity.OrganizationRole, error) {
	var models []*OrganizationRoleModel
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("is_deleted = ?", false).Find(&models).Error; err != nil {
		return nil, err
	}

	return RoleModelsToEntities(models), nil
}

func (r *OrganizationRoleRepository) FindById(ctx context.Context, id uint32) (*entity.OrganizationRole, error) {
	var model OrganizationRoleModel
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return RoleModelToEntity(&model), nil
}

func (r *OrganizationRoleRepository) Update(ctx context.Context, role *entity.OrganizationRole) (*entity.OrganizationRole, error) {
	model := RoleEntityToModel(role)
	if err := r.db.WithContext(ctx).Model(&OrganizationRoleModel{}).Where("id = ? AND is_deleted = ?", model.ID, false).Updates(model).Error; err != nil {
		return nil, err
	} else {
		if err := r.db.WithContext(ctx).Model(model).Association("Permissions").Replace(model.Permissions); err != nil {
			return nil, err
		}
	}

	return RoleModelToEntity(model), nil
}

func (r *OrganizationRoleRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&OrganizationRoleModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *OrganizationRoleRepository) GetByKey(ctx context.Context, key string) (*entity.OrganizationRole, error) {
	var model OrganizationRoleModel
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("key = ?", key).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return RoleModelToEntity(&model), nil
}

func (r *OrganizationRoleRepository) FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationRole, error) {
	var models []*OrganizationRoleModel
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("organization_id = ? AND is_deleted = ?", organizationId, false).Find(&models).Error; err != nil {
		return nil, err
	}
	return RoleModelsToEntities(models), nil
}

func (r *OrganizationRoleRepository) FindByOrganizationIdAndKey(ctx context.Context, organizationId uint32, key string) (*entity.OrganizationRole, error) {
	var model OrganizationRoleModel
	if err := r.db.WithContext(ctx).Preload("Permissions").Where("organization_id = ? AND key = ? AND is_deleted = ?", organizationId, key, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return RoleModelToEntity(&model), nil
}

func (r *OrganizationRoleRepository) FindByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationRole, uint32, error) {
	var models []*OrganizationRoleModel
	var total int64

	// Get total count
	err := r.db.Model(&OrganizationRoleModel{}).Where("organization_id = ? AND deleted_at IS NULL", organizationId).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated data with permissions
	offset := page * size
	err = r.db.WithContext(ctx).
		Model(&OrganizationRoleModel{}).
		Where("organization_id = ? AND deleted_at IS NULL", organizationId).
		Order("updated_at DESC").
		Offset(offset).
		Limit(size).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	return RoleModelsToEntities(models), uint32(total), nil
}

func (r *OrganizationRoleRepository) DeleteRole(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&OrganizationRoleModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *OrganizationRoleRepository) AddPermissionToRole(ctx context.Context, roleId uint32, permissionId uint32) error {
	return r.db.WithContext(ctx).Table("role_permissions").Create(&map[string]interface{}{
		"organization_role_model_id":       roleId,
		"organization_permission_model_id": permissionId,
	}).Error
}

func (r *OrganizationRoleRepository) FindByDomain(ctx context.Context, domain enums.DomainType) ([]*entity.OrganizationRole, error) {
	var models []*OrganizationRoleModel
	if err := r.db.WithContext(ctx).Preload("Color").Where("domain_type = ? AND is_deleted = ?", domain, false).Find(&models).Error; err != nil {
		return nil, err
	}
	return RoleModelsToEntities(models), nil
}

func (r *OrganizationRoleRepository) GetPermissionKeysByRoleId(ctx context.Context, roleId uint64) ([]string, error) {
	var permissionKeys []string

	err := r.db.WithContext(ctx).
		Table("role_permissions").
		Select("organization_permissions.key").
		Joins("JOIN organization_permissions ON organization_permissions.id = role_permissions.organization_permission_model_id").
		Where("role_permissions.organization_role_model_id = ?", roleId).
		Find(&permissionKeys).Error

	if err != nil {
		return nil, err
	}

	return permissionKeys, nil
}
