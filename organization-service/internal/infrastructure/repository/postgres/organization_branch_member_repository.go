package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type OrganizationBranchMemberModel struct {
	Id                   uint32    `gorm:"primaryKey;autoIncrement"`
	OrganizationBranchId uint32    `gorm:"not null"`	
	UserId               uint32    `gorm:"not null"`
	CreatedAt            time.Time `gorm:"autoCreateTime"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime"`
	CreatedBy            uint32
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
	}
}

type organizationBranchMemberRepository struct {
	db *gorm.DB
}

func NewOrganizationBranchMemberPostgresRepository(db *gorm.DB) repository.OrganizationBranchMemberRepository {
	return &organizationBranchMemberRepository{db: db}
}

func (r *organizationBranchMemberRepository) CreateOrganizationBranchMember(ctx context.Context, organizationBranchMember *entity.OrganizationBranchMember) (*entity.OrganizationBranchMember, error) {
	model := OrganizationBranchMemberEntityToModel(organizationBranchMember)
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchMemberModelToEntity(model), nil
}

func (r *organizationBranchMemberRepository) DeleteOrganizationBranchMember(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&OrganizationBranchMemberModel{}, id).Error
}

func (r *organizationBranchMemberRepository) GetOrganizationBranchMemberByOrganizationBranchId(ctx context.Context, organizationBranchId uint32) ([]*entity.OrganizationBranchMember, error) {
	var models []OrganizationBranchMemberModel
	err := r.db.WithContext(ctx).Where("organization_branch_id = ?", organizationBranchId).Find(&models).Error
	if err != nil {
		return nil, err
	}

	entities := make([]*entity.OrganizationBranchMember, len(models))
	for i, model := range models {
		entities[i] = OrganizationBranchMemberModelToEntity(&model)
	}
	return entities, nil
}

func (r *organizationBranchMemberRepository) GetOrganizationBranchMemberByUserIdAndOrganizationBranchId(ctx context.Context, userId uint32, organizationBranchId uint32) (*entity.OrganizationBranchMember, error) {
	var model OrganizationBranchMemberModel
	err := r.db.WithContext(ctx).Where("user_id = ? AND organization_branch_id = ?", userId, organizationBranchId).First(&model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationBranchMemberModelToEntity(&model), nil
}