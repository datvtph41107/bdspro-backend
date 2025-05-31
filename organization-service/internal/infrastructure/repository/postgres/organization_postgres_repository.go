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
}

func (OrganizationModel) TableName() string {
	return "organizations"
}

func ModelToEntity(model OrganizationModel) entity.Organization {
	return entity.Organization{
		ID: model.ID,
		Name: model.Name,
		TaxCode: model.TaxCode,
		BusinessLicenseUrl: model.BusinessLicenseUrl,
		Address: model.Address,
		Phone: model.Phone,
		Email: model.Email,
		Website: model.Website,
		LogoUrl: model.LogoUrl,
		Status: model.Status,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func EntityToModel(entity entity.Organization) OrganizationModel {
	return OrganizationModel{
		ID: entity.ID,
		Name: entity.Name,
		TaxCode: entity.TaxCode,
		BusinessLicenseUrl: entity.BusinessLicenseUrl,
		Address: entity.Address,
		Phone: entity.Phone,
		Email: entity.Email,
		Website: entity.Website,
		LogoUrl: entity.LogoUrl,
		Status: entity.Status,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,	
	}
}

func NewOrganizationPostgresRepository(db *gorm.DB) repository.OrganizationRepository {
	return &organizationPostgresRepository{db: db}
}

func (o *organizationPostgresRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := o.db.WithContext(ctx).Model(&OrganizationModel{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (o *organizationPostgresRepository) Create(ctx context.Context, organization *entity.Organization) (*entity.Organization, error) {
	model := EntityToModel(*organization)
	if err := o.db.WithContext(ctx).Create(&model).Error; err != nil {
		return nil, err
	}
	entity := ModelToEntity(model)
	return &entity, nil
}

func (o *organizationPostgresRepository) Delete(ctx context.Context, id uint32) (*entity.Organization, error) {
	var organization entity.Organization
	if err := o.db.WithContext(ctx).Delete(&OrganizationModel{}, id).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

func (o *organizationPostgresRepository) GetByEmail(ctx context.Context, email string) (*entity.Organization, error) {
	var organization entity.Organization
	if err := o.db.WithContext(ctx).Where("email = ?", email).First(&OrganizationModel{}).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

func (o *organizationPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.Organization, error) {
	var organization entity.Organization
	if err := o.db.WithContext(ctx).First(&OrganizationModel{}, id).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

func (o *organizationPostgresRepository) GetByTaxCode(ctx context.Context, taxCode string) (*entity.Organization, error) {
	var organization entity.Organization
	if err := o.db.WithContext(ctx).Where("tax_code = ?", taxCode).First(&OrganizationModel{}).Error; err != nil {
		return nil, err
	}
	return &organization, nil
}

func (o *organizationPostgresRepository) List(ctx context.Context, offset int, limit int) ([]*entity.Organization, error) {
	var organizations []*entity.Organization
	if err := o.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&[]OrganizationModel{}).Error; err != nil {
		return nil, err
	}
	return organizations, nil
}

func (o *organizationPostgresRepository) Update(ctx context.Context, organization *entity.Organization) (*entity.Organization, error) {
	model := EntityToModel(*organization)
	if err := o.db.WithContext(ctx).Save(&model).Error; err != nil {
		return nil, err
	}
	entity := ModelToEntity(model)
	return &entity, nil
}

func (o *organizationPostgresRepository) UpdateStatus(ctx context.Context, id uint32, status string) error {
	if err := o.db.WithContext(ctx).Model(&OrganizationModel{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}