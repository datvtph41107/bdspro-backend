package postgres

import (
	"context"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type OrganizationBranchModel struct {
	Id             uint32 `gorm:"primaryKey;autoIncrement"`
	OrganizationId uint32 `gorm:"not null;index:idx_org_branches_org_id"`
	Name           string `gorm:"index:idx_org_branches_name"`
	Address        string
	Phone          string
	Email          string
	ManagerId      uint32     `gorm:"autoIncrement;index:idx_org_branches_manager_id"`
	CreatedAt      time.Time  `gorm:"autoCreateTime;index:idx_org_branches_created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
	CreatedBy      uint32     `gorm:"index:idx_org_branches_created_by"`
	DeletedAt      *time.Time `gorm:"index"`
	IsDeleted      bool       `gorm:"default:false"`
	IsActive       bool       `gorm:"default:true;column:is_active"`
	Type           string     `gorm:"column:type;index:idx_org_branches_type"`
	Description    string     `gorm:"column:description"`
	TotalMember    uint32     `gorm:"-"`
	TotalAsset     uint32     `gorm:"-"`
	TotalDeal      uint32     `gorm:"-"`
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
		IsActive:       m.IsActive,
		Type:           entity.OrganizationBranchType(m.Type),
		TotalMember:    m.TotalMember,
		TotalAsset:     m.TotalAsset,
		TotalDeal:      m.TotalDeal,
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
		IsActive:       e.IsActive,
		Type:           string(e.Type),
	}
}

// @bind: organization/internal/domain/repository.OrganizationBranchRepository
type OrganizationBranchPostgres struct {
	db *gorm.DB
}

func NewOrganizationBranchPostgresRepository(db *gorm.DB) *OrganizationBranchPostgres {
	return &OrganizationBranchPostgres{db: db}
}

func (r *OrganizationBranchPostgres) CreateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error) {
	model := OrganizationBranchEntityToModel(organizationBranch)
	model.IsDeleted = false
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchModelToEntity(model), nil
}

func (r *OrganizationBranchPostgres) UpdateOrganizationBranch(ctx context.Context, organizationBranch *entity.OrganizationBranch) (*entity.OrganizationBranch, error) {
	model := OrganizationBranchEntityToModel(organizationBranch)
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", organizationBranch.Id, false).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchModelToEntity(model), nil
}

func (r *OrganizationBranchPostgres) DeleteOrganizationBranch(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&OrganizationBranchModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *OrganizationBranchPostgres) GetOrganizationBranchById(ctx context.Context, id uint32) (*entity.OrganizationBranch, error) {
	var model OrganizationBranchModel
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchModelToEntity(&model), nil
}

func (r *OrganizationBranchPostgres) GetOrganizationBranchByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationBranch, error) {
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

func (r *OrganizationBranchPostgres) GetOrganizationBranchByOrganizationIdWithPagination(ctx context.Context, organizationId uint32, page, size int) ([]*entity.OrganizationBranch, uint32, error) {
	// Model mới để chứa kết quả từ query JOIN
	type OrganizationBranchWithMemberCount struct {
		OrganizationBranchModel
		TotalMember uint32 `gorm:"column:total_member"`
	}

	var models []OrganizationBranchWithMemberCount
	if err := r.db.WithContext(ctx).
		Select("organization_branches.*, COUNT(organization_branch_members.id) as total_member").
		Joins("LEFT JOIN organization_branch_members ON organization_branches.id = organization_branch_members.organization_branch_id AND organization_branch_members.is_deleted = ?", false).
		Where("organization_branches.organization_id = ? AND organization_branches.is_deleted = ?", organizationId, false).
		Group("organization_branches.id").
		Offset(page * size).
		Limit(size).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]*entity.OrganizationBranch, len(models))
	for i, model := range models {
		branch := OrganizationBranchModelToEntity(&model.OrganizationBranchModel)
		branch.TotalMember = model.TotalMember
		entities[i] = branch
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

func (r *OrganizationBranchPostgres) UpdateIsActive(ctx context.Context, id uint32, isActive bool) error {
	return r.db.WithContext(ctx).Model(&OrganizationBranchModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_active": isActive,
		}).Error
}

func (r *OrganizationBranchPostgres) CountBranchByOrganizationId(ctx context.Context, organizationId uint32) (uint32, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrganizationBranchModel{}).
		Where("organization_id = ? AND is_deleted = ?", organizationId, false).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}
