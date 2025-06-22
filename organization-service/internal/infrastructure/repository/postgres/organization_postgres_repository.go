package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type organizationPostgresRepository struct {
	db *gorm.DB
}

type OrganizationModel struct {
	ID                 uint32    `gorm:"primaryKey;autoIncrement"`
	Name               string    `gorm:"not null;index"`
	TaxCode            string    `gorm:"not null;index:idx_organizations_tax_code"`
	BusinessLicenseUrl *string
	Address            *string
	Phone              *string
	Email              *string
	Website            *string
	LogoUrl            *string
	Status             string    `gorm:"index:idx_organizations_status_created_at"`
	CreatedAt          time.Time `gorm:"autoCreateTime;index:idx_organizations_status_created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
	CreatedBy          uint32    `gorm:"index:idx_organizations_created_by"`
	UpdatedBy          uint32
	DeletedAt          *time.Time `gorm:"index"`
	IsDeleted          bool       `gorm:"default:false"`
}

func (OrganizationModel) TableName() string {
	return "organizations"
}

func OrganizationModelToEntity(model *OrganizationModel) *entity.Organization {
	return &entity.Organization{
		ID:                 model.ID,
		Name:               model.Name,
		TaxCode:            model.TaxCode,
		BusinessLicenseUrl: model.BusinessLicenseUrl,
		Address:            model.Address,
		Phone:              model.Phone,
		Email:              model.Email,
		Website:            model.Website,
		LogoUrl:            model.LogoUrl,
		Status:             model.Status,
		CreatedAt:          model.CreatedAt,
		UpdatedAt:          model.UpdatedAt,
		CreatedBy:          model.CreatedBy,
		UpdatedBy:          model.UpdatedBy,
	}
}

func OrganizationEntityToModel(entity *entity.Organization) *OrganizationModel {
	return &OrganizationModel{
		ID:                 entity.ID,
		Name:               entity.Name,
		TaxCode:            entity.TaxCode,
		BusinessLicenseUrl: entity.BusinessLicenseUrl,
		Address:            entity.Address,
		Phone:              entity.Phone,
		Email:              entity.Email,
		Website:            entity.Website,
		LogoUrl:            entity.LogoUrl,
		Status:             entity.Status,
		CreatedAt:          entity.CreatedAt,
		UpdatedAt:          entity.UpdatedAt,
		CreatedBy:          entity.CreatedBy,
		UpdatedBy:          entity.UpdatedBy,
		IsDeleted:          false,
	}
}

func OrganizationModelsToEntities(models []*OrganizationModel) []*entity.Organization {
	entities := make([]*entity.Organization, len(models))
	for i, model := range models {
		entities[i] = OrganizationModelToEntity(model)
	}
	return entities
}

func NewOrganizationPostgresRepository(db *gorm.DB) repository.OrganizationRepository {
	return &organizationPostgresRepository{db: db}
}

func (o *organizationPostgresRepository) CreateOrganizationWithMember(ctx context.Context,
	organization *entity.Organization,
	member *entity.OrganizationMember,
) (*entity.Organization, *entity.OrganizationMember, error) {
	tx := o.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	organizationModel := OrganizationEntityToModel(organization)
	organizationModel.IsDeleted = false
	if err := tx.Create(&organizationModel).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	role := &OrganizationRoleModel{
		Name:           appconfig.ADMIN_ROLE_NAME,
		Key:            appconfig.ADMIN_ROLE_KEY,
		OrganizationId: organizationModel.ID,
	}
	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	memberModel := OrganizationMemberEntityToModel(member)
	memberModel.OrganizationID = organizationModel.ID
	memberModel.Role = role.ID
	if err := tx.Create(&memberModel).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	return OrganizationModelToEntity(organizationModel), OrganizationMemberModelToEntity(memberModel), nil
}

func (o *organizationPostgresRepository) FindByTaxCode(ctx context.Context, taxCode string) (*entity.Organization, error) {
	var organizationModel *OrganizationModel
	if err := o.db.WithContext(ctx).Where("tax_code = ? AND is_deleted = ?", taxCode, false).First(&organizationModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return OrganizationModelToEntity(organizationModel), nil
}

func (o *organizationPostgresRepository) UpdateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error) {
	model := OrganizationEntityToModel(organization)
	if err := o.db.WithContext(ctx).Model(&model).Where("id = ? AND is_deleted = ?", organization.ID, false).Updates(model).Error; err != nil {
		return nil, err
	}
	return OrganizationModelToEntity(model), nil
}

func (o *organizationPostgresRepository) FindById(ctx context.Context, id uint32) (*entity.Organization, error) {
	var organizationModel *OrganizationModel
	if err := o.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&organizationModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return OrganizationModelToEntity(organizationModel), nil
}

func (o *organizationPostgresRepository) FindByUserId(ctx context.Context, userId uint32) ([]*entity.Organization, error) {
	var organizations []*OrganizationModel
	if err := o.db.WithContext(ctx).
		Joins("LEFT JOIN organization_members ON organizations.id = organization_members.organization_id").
		Where("(organizations.created_by = ? OR organization_members.user_id = ?) AND organizations.is_deleted = ?", userId, userId, false).
		Distinct().
		Order("organizations.created_at DESC").
		Find(&organizations).Error; err != nil {
		return nil, err
	}
	return OrganizationModelsToEntities(organizations), nil
}

func (r *organizationPostgresRepository) FindByUserIdWithPagination(ctx context.Context, userId uint32, offset, limit int) ([]*entity.Organization, uint32, error) {
	var models []*OrganizationModel
	if err := r.db.WithContext(ctx).
		Where("created_by = ? AND is_deleted = ?", userId, false).
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.Organization, len(models))
	for i, model := range models {
		entities[i] = OrganizationModelToEntity(model)
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&OrganizationModel{}).
		Where("created_by = ? AND is_deleted = ?", userId, false).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}

// Add soft delete function
func (r *organizationPostgresRepository) DeleteOrganization(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&OrganizationModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted":  true,
			"deleted_at": time.Now(),
		}).Error
}
