package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type OrganizationBranchModel struct {
	Id             uint32 `gorm:"primaryKey;autoIncrement"`
	OrganizationId uint32 `gorm:"not null;index:idx_org_branches_org_id"`
	Name           string `gorm:"index:idx_org_branches_name"`
	Address        string
	Phone          string
	Email          string
	ManagerId      uint32    `gorm:"autoIncrement;index:idx_org_branches_manager_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime;index:idx_org_branches_created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	CreatedBy      uint32    `gorm:"index:idx_org_branches_created_by"`
}

func (m *OrganizationBranchModel) TableName() string {
	return "organization_branches"
}

func OrganizationBranchModelToEntity(m *OrganizationBranchModel) *entity.OrganizationBranch {
	return &entity.OrganizationBranch{
		Id:             m.Id,
		OrganizationId: m.OrganizationId,
		Name:           m.Name,
		Address:        m.Address,
		Phone:          m.Phone,
		Email:          m.Email,
		ManagerId:      m.ManagerId,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		CreatedBy:      m.CreatedBy,
	}
}

func OrganizationBranchEntityToModel(e *entity.OrganizationBranch) *OrganizationBranchModel {
	return &OrganizationBranchModel{
		Id:             e.Id,
		OrganizationId: e.OrganizationId,
		Name:           e.Name,
		Address:        e.Address,
		Phone:          e.Phone,
		Email:          e.Email,
		ManagerId:      e.ManagerId,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		CreatedBy:      e.CreatedBy,
	}
}

type organizationBranchRepository struct {
	db *gorm.DB
}

func NewOrganizationBranchPostgresRepository(db *gorm.DB) repository.OrganizationBranchRepository {
	return &organizationBranchRepository{db: db}
}

func (r *organizationBranchRepository) CreateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error) {
	model := OrganizationBranchEntityToModel(organizationBranch)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchModelToEntity(model), nil
}

func (r *organizationBranchRepository) UpdateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error) {
	model := OrganizationBranchEntityToModel(organizationBranch)
	err := r.db.WithContext(ctx).Model(&OrganizationBranchModel{}).Where("id = ?", model.Id).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchModelToEntity(model), nil
}

func (r *organizationBranchRepository) DeleteOrganizationBranch(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&OrganizationBranchModel{}, id).Error
}

func (r *organizationBranchRepository) GetOrganizationBranchById(ctx context.Context, id uint32) (*entity.OrganizationBranch, error) {
	var model OrganizationBranchModel
	err := r.db.WithContext(ctx).First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchModelToEntity(&model), nil
}

func (r *organizationBranchRepository) GetOrganizationBranchByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationBranch, error) {
	var models []OrganizationBranchModel
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Find(&models).Error
	if err != nil {
		return nil, err
	}
	entities := make([]*entity.OrganizationBranch, len(models))
	for i, model := range models {
		entities[i] = OrganizationBranchModelToEntity(&model)
	}
	return entities, nil
}
