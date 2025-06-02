package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type organizationPostgresRepository struct {
	db *gorm.DB
}

type OrganizationModel struct {
	ID                 uint32 `gorm:"primaryKey;autoIncrement"`
	Name               string `gorm:"not null"`
	TaxCode            string `gorm:"not null"`
	BusinessLicenseUrl *string
	Address            *string
	Phone              *string
	Email              *string
	Website            *string
	LogoUrl            *string
	Status             string
	CreatedAt          time.Time `gorm:"autoCreateTime"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime"`
	CreatedBy          uint32
	UpdatedBy          uint32
}

func (OrganizationModel) TableName() string {
	return "organizations"
}

func OrganizationModelToEntity(model OrganizationModel) entity.Organization {
	return entity.Organization{
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

func OrganizationEntityToModel(entity entity.Organization) OrganizationModel {
	return OrganizationModel{
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
	}
}

func NewOrganizationPostgresRepository(db *gorm.DB) repository.OrganizationRepository {
	return &organizationPostgresRepository{db: db}
}

func (o *organizationPostgresRepository) CreateOrganizationWithMember(ctx context.Context,
	organization *entity.Organization,
	member *entity.OrganizationMember,
) (*entity.Organization, *entity.OrganizationMember, error) {
	tx := o.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	organizationModel := OrganizationEntityToModel(*organization)
	if err := tx.Create(&organizationModel).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	memberModel := OrganizationMemberEntityToModel(*member)
	memberModel.OrganizationID = organizationModel.ID
	if err := tx.Create(&memberModel).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}
	org := OrganizationModelToEntity(organizationModel)
	mem := OrganizationMemberModelToEntity(memberModel)
	return &org, &mem, nil
}

func (o *organizationPostgresRepository) ExistByTaxCode(ctx context.Context, taxCode string) (bool, error) {
	var count int64
	if err := o.db.WithContext(ctx).Model(&OrganizationModel{}).Where("tax_code = ?", taxCode).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
