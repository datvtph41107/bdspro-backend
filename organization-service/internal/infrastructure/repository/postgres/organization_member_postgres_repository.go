package postgres

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type OrganizationMemberModel struct {
	ID             uint32    `gorm:"primaryKey;autoIncrement"`
	OrganizationID uint32    `gorm:"not null;index:idx_org_members_organization_id,priority:1;index:idx_org_members_user_org,priority:2;index:idx_org_members_status_org,priority:2"`
	UserID         uint32    `gorm:"not null;index:idx_org_members_user_org,priority:1"`
	Role           uint32    `gorm:"not null;index:idx_org_members_role"`
	Status         uint32    `gorm:"default:1;index:idx_org_members_status_org,priority:1"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
	RemovedAt      *time.Time

	Organization *OrganizationModel
}

type OrganizationMemberRoleModel struct {
	ID          uint32   `gorm:"column:id"`
	UserId      uint32   `gorm:"column:user_id"`
	Role        string   `gorm:"column:role"`
	Permissions []string `gorm:"column:permissions"`
}

func (OrganizationMemberModel) TableName() string {
	return "organization_members"
}

func OrganizationMemberModelToEntity(model *OrganizationMemberModel) *entity.OrganizationMember {
	return &entity.OrganizationMember{
		ID:             model.ID,
		OrganizationID: model.OrganizationID,
		UserID:         model.UserID,
		RoleId:         model.Role,
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

func (r *OrganizationMemberPostgresRepository) FindByUserIdAndOrganizationIdWithRole(ctx context.Context, userId, organizationId uint32) (*entity.OrganizationMember, error) {
	var result struct {
		ID             uint32 `gorm:"column:id"`
		UserID         uint32 `gorm:"column:user_id"`
		OrganizationID uint32 `gorm:"column:organization_id"`
		RoleKey        string `gorm:"column:role_key"`
		Permissions    string `gorm:"column:permissions"`
	}

	err := r.db.WithContext(ctx).
		Table("organization_members om").
		Select("om.id, om.user_id, om.organization_id, or2.key as role_key, string_agg(op.key, ',') as permissions").
		Joins("left join organization_roles or2 on om.role = or2.id").
		Joins("left join role_permissions rp on or2.id = rp.organization_role_model_id").
		Joins("left join organization_permissions op on rp.organization_permission_model_id = op.id").
		Where("om.user_id = ? AND om.organization_id = ?", userId, organizationId).
		Group("om.id, om.user_id, om.organization_id, or2.key").
		First(&result).Error
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	permissions := strings.Split(result.Permissions, ",")
	return &entity.OrganizationMember{
		ID:             result.ID,
		UserID:         result.UserID,
		OrganizationID: result.OrganizationID,
		Role:           result.RoleKey,
		Permissions:    permissions,
	}, nil
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

func (r *OrganizationMemberPostgresRepository) CountMemberByOrganizationId(ctx context.Context, organizationId uint32) (uint32, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&OrganizationMemberModel{}).Where("organization_id = ?", organizationId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return uint32(count), nil
}