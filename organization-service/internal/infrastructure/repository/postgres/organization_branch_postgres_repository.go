package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type OrganizationBranchModel struct {
	Id             uint32    `gorm:"primaryKey;autoIncrement"`
	OrganizationId uint32    `gorm:"not null;index:idx_org_branches_org_id"`
	Name           string    `gorm:"index:idx_org_branches_name"`
	Address        string
	Phone          string
	Email          string
	ManagerId      uint32    `gorm:"autoIncrement;index:idx_org_branches_manager_id"`
	CreatedAt      time.Time `gorm:"autoCreateTime;index:idx_org_branches_created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
	CreatedBy      uint32    `gorm:"index:idx_org_branches_created_by"`
	DeletedAt      *time.Time `gorm:"index"`
	IsDeleted      bool       `gorm:"default:false"`
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
		IsDeleted:      false,
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
	model.IsDeleted = false
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchModelToEntity(model), nil
}

func (r *organizationBranchRepository) UpdateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error) {
	model := OrganizationBranchEntityToModel(organizationBranch)
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", organizationBranch.Id, false).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchModelToEntity(model), nil
}

func (r *organizationBranchRepository) DeleteOrganizationBranch(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&OrganizationBranchModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted":  true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *organizationBranchRepository) GetOrganizationBranchById(ctx context.Context, id uint32) (*entity.OrganizationBranch, error) {
	var model OrganizationBranchModel
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchModelToEntity(&model), nil
}

func (r *organizationBranchRepository) GetOrganizationBranchByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationBranch, error) {
	var models []*OrganizationBranchModel
	err := r.db.WithContext(ctx).Where("organization_id = ? AND is_deleted = ?", organizationId, false).Find(&models).Error
	if err != nil {
		return nil, err
	}
	entities := make([]*entity.OrganizationBranch, len(models))
	for i, model := range models {
		entities[i] = OrganizationBranchModelToEntity(model)
	}
	return entities, nil
}

func (r *organizationBranchRepository) GetOrganizationBranchByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationBranch, uint32, error) {
	var models []OrganizationBranchModel
	if err := r.db.WithContext(ctx).
		Where("organization_id = ? AND is_deleted = ?", organizationId, false).
		Offset(page * size).
		Limit(size).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.OrganizationBranch, len(models))
	for i, model := range models {
		entities[i] = OrganizationBranchModelToEntity(&model)
	}
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&OrganizationBranchModel{}).
		Where("organization_id = ? AND is_deleted = ?", organizationId, false).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}