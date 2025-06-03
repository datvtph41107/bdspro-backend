package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type OrganizationMemberModel struct {
	ID             uint32    `gorm:"primaryKey;autoIncrement"`
	OrganizationID uint32    `gorm:"not null"`
	UserID         uint32    `gorm:"not null"`
	Role           uint32    `gorm:"not null"`
	Status         uint32    `gorm:"default:1"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
	RemovedAt      *time.Time

	Organization *OrganizationModel
}

func (OrganizationMemberModel) TableName() string {
	return "organization_members"
}

func OrganizationMemberModelToEntity(model *OrganizationMemberModel) *entity.OrganizationMember {
	return &entity.OrganizationMember{
		ID:             model.ID,
		OrganizationID: model.OrganizationID,
		UserID:         model.UserID,
		RoleId:           model.Role,
		Status:         entity.OrganizationMemberStatus(model.Status),
		JoinedAt:       model.JoinedAt,
		RemovedAt:      model.RemovedAt,
	}
}

func OrganizationMemberEntityToModel(entity *entity.OrganizationMember) *OrganizationMemberModel {
	return &OrganizationMemberModel{
		ID:             entity.ID,
		OrganizationID: entity.OrganizationID,
		UserID:         entity.UserID,
		Role:           uint32(entity.RoleId),
		Status:         uint32(entity.Status),
		JoinedAt:       entity.JoinedAt,
		RemovedAt:      entity.RemovedAt,
	}
}

func OrganizationMemberModelsToEntities(models []*OrganizationMemberModel) []*entity.OrganizationMember {
	entities := make([]*entity.OrganizationMember, len(models))
	for i, model := range models {
		entities[i] = OrganizationMemberModelToEntity(model)
	}
	return entities
}

type OrganizationMemberPostgresRepository struct {
	db *gorm.DB
}

func NewOrganizationMemberPostgresRepository(db *gorm.DB) repository.OrganizationMemberRepository {
	return &OrganizationMemberPostgresRepository{db: db}
}

func (r *OrganizationMemberPostgresRepository) ExistByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).Where("user_id = ? AND organization_id = ?", userId, organizationId).Count(&count).Error
	return count > 0, err
}

func (r *OrganizationMemberPostgresRepository) FindByUserIdAndOrganizationId(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error) {
	var organizationMemberModel *OrganizationMemberModel
	err := r.db.WithContext(ctx).Where("user_id = ? AND organization_id = ?", userId, organizationId).First(&organizationMemberModel).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelToEntity(organizationMemberModel), nil
}

func (r *OrganizationMemberPostgresRepository) CreateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	model := OrganizationMemberEntityToModel(organizationMember)
	err := r.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelToEntity(model), nil
}

func (r *OrganizationMemberPostgresRepository) FindById(ctx context.Context, id uint32) (*entity.OrganizationMember, error) {
	var organizationMemberModel *OrganizationMemberModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&organizationMemberModel).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelToEntity(organizationMemberModel), nil
}	

func (r *OrganizationMemberPostgresRepository) UpdateOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	model := OrganizationMemberEntityToModel(organizationMember)
	err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).Where("id = ?", organizationMember.ID).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelToEntity(model), nil
}

func (r *OrganizationMemberPostgresRepository) DeleteOrganizationMember(ctx context.Context, organizationMember *entity.OrganizationMember) (*entity.OrganizationMember, error) {
	err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).Where("id = ?", organizationMember.ID).Delete(&OrganizationMemberModel{}).Error
	if err != nil {
		return nil, err
	}
	return organizationMember, nil
}
func (r *OrganizationMemberPostgresRepository) FindByOrganizationId(ctx context.Context, organizationId uint32) ([]*entity.OrganizationMember, error) {
	var organizationMembers []*OrganizationMemberModel
	err := r.db.WithContext(ctx).Where("organization_id = ?", organizationId).Find(&organizationMembers).Error
	if err != nil {
		return nil, err
	}
	return OrganizationMemberModelsToEntities(organizationMembers), nil
}
