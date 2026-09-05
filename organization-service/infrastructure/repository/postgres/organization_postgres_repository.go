package postgres

import (
	"context"
	"errors"
	"time"

	"organization/env"
	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

// @bind: organization/internal/domain/repository.OrganizationRepository
type OrganizationPostgresRepository struct {
	db *gorm.DB
}

func NewOrganizationPostgresRepository(db *gorm.DB) *OrganizationPostgresRepository {
	return &OrganizationPostgresRepository{db: db}
}

func (o *OrganizationPostgresRepository) CreateOrganizationWithMember(ctx context.Context,
	organization *entity.Organization,
	member *entity.OrganizationMember,
) (*entity.Organization, *entity.OrganizationMember, error) {
	tx := o.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, nil, tx.Error
	}

	// entity.Organization := OrganizationEntityToModel(organization)
	// entity.Organization.IsDeleted = false
	if err := tx.Create(&organization).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	role := &OrganizationRoleModel{
		Name:           env.ADMIN_ROLE_NAME,
		Key:            env.ADMIN_ROLE_KEY,
		OrganizationId: uint32(organization.ID),
		IsDefault:      true,
	}
	if err := tx.Create(&role).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	member.OrganizationID = uint32(organization.ID)
	member.RoleId = uint64(role.ID)
	memberModel := OrganizationMemberEntityToModel(member)
	if err := tx.Create(memberModel).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if len(organization.BusinessDomainIds) > 0 {
		businessDomains := make([]*entity.BusinessDomain, len(organization.BusinessDomainIds))
		for i, domainId := range organization.BusinessDomainIds {
			businessDomains[i] = &entity.BusinessDomain{ID: domainId}
		}
		if err := tx.Model(&organization).Association("BusinessDomains").Append(businessDomains); err != nil {
			tx.Rollback()
			return nil, nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	return organization, OrganizationMemberModelToEntity(memberModel), nil
}

func (o *OrganizationPostgresRepository) FindByTaxCode(ctx context.Context, taxCode string) (*entity.Organization, error) {
	var organization *entity.Organization
	if err := o.db.WithContext(ctx).Where("tax_code = ? AND is_deleted = ?", taxCode, false).First(&organization).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return organization, nil
}

func (o *OrganizationPostgresRepository) UpdateOrganization(ctx context.Context, organization *entity.Organization) (*entity.Organization, error) {
	tx := o.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var model entity.Organization
	if err := tx.First(&model, "id = ? AND is_deleted = ?", organization.ID, false).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Cập nhật trường khác (ngoài BusinessDomains)
	updatedModel := organization
	updatedModel.ID = model.ID
	updatedModel.BusinessDomains = nil

	if err := tx.Model(&model).Updates(updatedModel).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Đảm bảo model đã attach vào transaction
	if err := tx.Save(&model).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Replace BusinessDomains
	if len(organization.BusinessDomainIds) > 0 {
		businessDomains := make([]*entity.BusinessDomain, len(organization.BusinessDomainIds))
		for i, id := range organization.BusinessDomainIds {
			businessDomains[i] = &entity.BusinessDomain{ID: id}
		}

		// ⚠️ Dùng Association từ tx.Model
		if err := tx.Model(&model).Association("BusinessDomains").Replace(businessDomains); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return updatedModel, nil
}

func (o *OrganizationPostgresRepository) FindById(ctx context.Context, id uint32) (*entity.Organization, error) {
	var org entity.Organization
	err := o.db.WithContext(ctx).
		Debug().
		Preload("BusinessDomains").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &org, nil
}

func (o *OrganizationPostgresRepository) FindByUserId(ctx context.Context, userId uint32) ([]*entity.Organization, error) {
	var organizations []*entity.Organization
	if err := o.db.WithContext(ctx).
		Joins("LEFT JOIN organization_members ON organizations.id = organization_members.organization_id").
		Where("(organizations.created_by = ? OR organization_members.user_id = ?) AND organizations.is_deleted = ?", userId, userId, false).
		Distinct().
		Order("organizations.created_at DESC").
		Find(&organizations).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

func (r *OrganizationPostgresRepository) FindByUserIdWithPagination(ctx context.Context, userId uint32, offset, limit int) ([]*entity.Organization, uint32, error) {
	var models []*entity.Organization
	if err := r.db.WithContext(ctx).
		Where("created_by = ? AND is_deleted = ?", userId, false).
		Offset(offset).
		Limit(limit).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}
	// entities := make([]*entity.Organization, len(models))
	// for i, model := range models {
	// 	entities[i] = entity.OrganizationToEntity(model)
	// }

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("created_by = ? AND is_deleted = ?", userId, false).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return models, uint32(count), nil
}

func (r *OrganizationPostgresRepository) GetAllOrganizationsWithPagination(ctx context.Context, offset, limit int) ([]*entity.Organization, uint32, error) {
	var models []*entity.Organization
	if err := r.db.WithContext(ctx).
		Preload("BusinessDomains").
		Where("is_deleted = ?", false).
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, 0, err
	}
	// entities := make([]*entity.Organization, len(models))
	// for i, model := range models {
	// 	entities[i] = entity.OrganizationToEntity(model)
	// }

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("is_deleted = ?", false).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return models, uint32(count), nil
}

// Add soft delete function
func (r *OrganizationPostgresRepository) DeleteOrganization(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&entity.Organization{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

// GetAllOrganizationsWithMemberCounts - Lấy organizations với total members trong 1 query
func (r *OrganizationPostgresRepository) GetAllOrganizationsWithMemberCounts(ctx context.Context, offset, limit int) ([]*entity.OrganizationWithMemberCount, uint32, error) {
	var results []*entity.OrganizationWithMemberCount

	// Query với LEFT JOIN để lấy count members
	err := r.db.WithContext(ctx).
		Table("organizations o").
		Select(`
			o.id, o.name, o.tax_code, o.business_license_url, o.address, o.detail_address,
			o.phone, o.email, o.website, o.logo_url, o.status, o.approve_investor,
			o.account_bank_investor, o.founded_at, o.description, o.created_at, o.updated_at,
			COALESCE(COUNT(om.id), 0) as total_members
		`).
		Joins("LEFT JOIN organization_members om ON o.id = om.organization_id AND om.is_deleted = false").
		Where("o.is_deleted = ?", false).
		Group("o.id, o.name, o.tax_code, o.business_license_url, o.address, o.detail_address, o.phone, o.email, o.website, o.logo_url, o.status, o.approve_investor, o.account_bank_investor, o.founded_at, o.description, o.created_at, o.updated_at").
		Order("o.created_at DESC").
		Offset(offset).
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, 0, err
	}

	// Lấy total count
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Organization{}).
		Where("is_deleted = ?", false).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}

	return results, uint32(count), nil
}

func (r *OrganizationPostgresRepository) GetByIds(ctx context.Context, ids []uint64) ([]*entity.Organization, error) {
	var orgs []*entity.Organization
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}
