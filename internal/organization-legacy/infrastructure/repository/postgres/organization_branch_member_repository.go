package postgres

import (
	"context"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type OrganizationBranchMemberModel struct {
	Id                   uint32    `gorm:"primaryKey;autoIncrement"`
	OrganizationBranchId uint32    `gorm:"not null;index:idx_org_branch_members_branch_id,priority:1;index:idx_org_branch_members_user_branch,priority:2"`
	UserId               uint32    `gorm:"not null;index:idx_org_branch_members_user_id;index:idx_org_branch_members_user_branch,priority:1"`
	CreatedAt            time.Time `gorm:"autoCreateTime;index:idx_org_branch_members_created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
	CreatedBy            uint32    `gorm:"index:idx_org_branch_members_created_by"`

	RoleId    uint64     `gorm:"column:role_id"`
	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}

func (m *OrganizationBranchMemberModel) TableName() string {
	return "organization_branch_members"
}

func OrganizationBranchMemberModelToEntity(m *OrganizationBranchMemberModel) *entity.OrganizationBranchMember {
	return &entity.OrganizationBranchMember{
		Id:                   m.Id,
		OrganizationBranchId: m.OrganizationBranchId,
		UserId:               m.UserId,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
		CreatedBy:            m.CreatedBy,
		RoleId:               m.RoleId,
	}
}

func OrganizationBranchMemberEntityToModel(e *entity.OrganizationBranchMember) *OrganizationBranchMemberModel {
	return &OrganizationBranchMemberModel{
		Id:                   e.Id,
		OrganizationBranchId: e.OrganizationBranchId,
		UserId:               e.UserId,
		CreatedAt:            e.CreatedAt,
		UpdatedAt:            e.UpdatedAt,
		CreatedBy:            e.CreatedBy,
		RoleId:               e.RoleId,
	}
}

// @bind: organization/internal/domain/repository.OrganizationBranchMemberRepository
type OrganizationBranchMemberRepository struct {
	db *gorm.DB
}

func NewOrganizationBranchMemberPostgresRepository(db *gorm.DB) *OrganizationBranchMemberRepository {
	return &OrganizationBranchMemberRepository{db: db}
}

func (r *OrganizationBranchMemberRepository) CreateOrganizationBranchMember(ctx context.Context, organizationBranchMember *entity.OrganizationBranchMember) (*entity.OrganizationBranchMember, error) {
	model := OrganizationBranchMemberEntityToModel(organizationBranchMember)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchMemberModelToEntity(model), nil
}

func (r *OrganizationBranchMemberRepository) DeleteOrganizationBranchMember(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&OrganizationBranchMemberModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *OrganizationBranchMemberRepository) GetOrganizationBranchMemberByOrganizationBranchId(ctx context.Context, organizationBranchId uint32) ([]*entity.OrganizationBranchMember, error) {
	var models []OrganizationBranchMemberModel
	err := r.db.WithContext(ctx).Where("organization_branch_id = ? AND is_deleted = ?", organizationBranchId, false).Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.OrganizationBranchMember, len(models))
	for i, model := range models {
		entities[i] = OrganizationBranchMemberModelToEntity(&model)
	}
	return entities, nil
}

func (r *OrganizationBranchMemberRepository) GetOrganizationBranchMemberByUserIdAndOrganizationBranchId(ctx context.Context, userId uint32, organizationBranchId uint32) (*entity.OrganizationBranchMember, error) {
	var model OrganizationBranchMemberModel
	err := r.db.WithContext(ctx).Where("user_id = ? AND organization_branch_id = ? AND is_deleted = ?", userId, organizationBranchId, false).First(&model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchMemberModelToEntity(&model), nil
}

func (r *OrganizationBranchMemberRepository) GetOrganizationBranchMemberByUserIdsAndOrganizationBranchId(ctx context.Context, userIds []uint64, organizationBranchId uint32) ([]*entity.OrganizationBranchMember, error) {
	var models []OrganizationBranchMemberModel
	err := r.db.WithContext(ctx).Where("user_id IN (?) AND organization_branch_id = ? AND deleted_at IS NULL", userIds, organizationBranchId).Find(&models).Error
	if err != nil {
		return nil, err
	}
	entities := make([]*entity.OrganizationBranchMember, len(models))
	for i, model := range models {
		entities[i] = OrganizationBranchMemberModelToEntity(&model)
	}
	return entities, nil
}

func (r *OrganizationBranchMemberRepository) GetOrganizationBranchMemberByOrganizationBranchIdWithPagination(ctx context.Context, organizationBranchId uint32, page, size int) ([]*entity.OrganizationBranchMember, uint32, error) {
	var models []OrganizationBranchMemberModel
	err := r.db.WithContext(ctx).Where("organization_branch_id = ? AND is_deleted = ?", organizationBranchId, false).Offset(page * size).Limit(size).Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	entities := make([]*entity.OrganizationBranchMember, len(models))
	for i, model := range models {
		entities[i] = OrganizationBranchMemberModelToEntity(&model)
	}

	var total int64
	if err := r.db.WithContext(ctx).Model(&OrganizationBranchMemberModel{}).Where("organization_branch_id = ? AND is_deleted = ?", organizationBranchId, false).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return entities, uint32(total), nil
}
