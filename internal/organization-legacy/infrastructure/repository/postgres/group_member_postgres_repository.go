package postgres

import (
	"context"
	"errors"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type GroupMemberModel struct {
	Id        uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_group_members_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"not null;index:idx_group_members_created_by"`
	UpdatedBy uint32    `gorm:"not null"`

	GroupId uint32                   `gorm:"not null;index:idx_group_members_group_id,priority:1;index:idx_group_members_group_user,priority:1"`
	UserId  uint32                   `gorm:"not null;index:idx_group_members_user_id;index:idx_group_members_group_user,priority:2"`
	Role    entity.GroupMemberRole   `gorm:"not null;index:idx_group_members_role"`
	Status  entity.GroupMemberStatus `gorm:"not null;index:idx_group_members_status"`

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
	RoleId    *uint64    `gorm:"index:idx_group_members_role_id"`
}

func (GroupMemberModel) TableName() string {
	return "group_members"
}

func GroupMemberModelToEntity(model *GroupMemberModel) *entity.GroupMember {
	return &entity.GroupMember{
		Id:        model.Id,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		CreatedBy: model.CreatedBy,
		UpdatedBy: model.UpdatedBy,
		GroupId:   model.GroupId,
		UserId:    model.UserId,
		Role:      model.Role,
		Status:    model.Status,
		RoleId:    model.RoleId,
	}
}

func GroupMemberEntityToModel(entity *entity.GroupMember) *GroupMemberModel {
	return &GroupMemberModel{
		Id:        entity.Id,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
		GroupId:   entity.GroupId,
		UserId:    entity.UserId,
		Role:      entity.Role,
		Status:    entity.Status,
		RoleId:    entity.RoleId,
	}
}

// @bind: organization/internal/domain/repository.GroupMemberRepository
type GroupMemberPostgresRepository struct {
	db *gorm.DB
}

func NewGroupMemberPostgresRepository(db *gorm.DB) *GroupMemberPostgresRepository {
	return &GroupMemberPostgresRepository{db: db}
}

func (r *GroupMemberPostgresRepository) Create(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error) {
	model := GroupMemberEntityToModel(member)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupMemberModelToEntity(model), nil
}

func (r *GroupMemberPostgresRepository) CreateBatch(ctx context.Context, members []*entity.GroupMember) ([]*entity.GroupMember, error) {
	models := make([]*GroupMemberModel, len(members))
	for i, member := range members {
		models[i] = GroupMemberEntityToModel(member)
	}
	if err := r.db.WithContext(ctx).Create(models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.GroupMember, len(models))
	for i, model := range models {
		entities[i] = GroupMemberModelToEntity(model)
	}
	return entities, nil
}

func (r *GroupMemberPostgresRepository) Update(ctx context.Context, member *entity.GroupMember) (*entity.GroupMember, error) {
	model := GroupMemberEntityToModel(member)
	if err := r.db.WithContext(ctx).Model(&GroupMemberModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return GroupMemberModelToEntity(model), nil
}

func (r *GroupMemberPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&GroupMemberModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *GroupMemberPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.GroupMember, error) {
	var model *GroupMemberModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupMemberModelToEntity(model), nil
}

func (r *GroupMemberPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupMember, error) {
	var models []*GroupMemberModel
	if err := r.db.WithContext(ctx).Where("group_id = ? AND is_deleted = ?", groupID, false).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.GroupMember, len(models))
	for i, model := range models {
		entities[i] = GroupMemberModelToEntity(model)
	}
	return entities, nil
}

func (r *GroupMemberPostgresRepository) GetByUserID(ctx context.Context, userID uint32) ([]*entity.GroupMember, error) {
	var models []*GroupMemberModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_deleted = ?", userID, false).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.GroupMember, len(models))
	for i, model := range models {
		entities[i] = GroupMemberModelToEntity(model)
	}
	return entities, nil
}

func (r *GroupMemberPostgresRepository) UpdateStatus(ctx context.Context, id uint32, status entity.GroupMemberStatus) error {
	return r.db.WithContext(ctx).Model(&GroupMemberModel{}).Where("id = ? AND is_deleted = ?", id, false).Update("status", status).Error
}

func (r *GroupMemberPostgresRepository) UpdateRole(ctx context.Context, id uint32, role entity.GroupMemberRole) error {
	return r.db.WithContext(ctx).Model(&GroupMemberModel{}).Where("id = ? AND is_deleted = ?", id, false).Update("role", role).Error
}

func (r *GroupMemberPostgresRepository) GetByGroupIDAndUserID(ctx context.Context, groupID uint32, userID uint32) (*entity.GroupMember, error) {
	var model *GroupMemberModel
	if err := r.db.WithContext(ctx).Where("group_id = ? AND user_id = ? AND is_deleted = ?", groupID, userID, false).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return GroupMemberModelToEntity(model), nil
}

func (r *GroupMemberPostgresRepository) GetByGroupIDAndUserIDs(ctx context.Context, groupID uint32, userIDs []uint64) ([]*entity.GroupMember, error) {
	var models []*GroupMemberModel
	if err := r.db.WithContext(ctx).Where("group_id = ? AND user_id IN (?) AND is_deleted = ?", groupID, userIDs, false).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.GroupMember, len(models))
	for i, model := range models {
		entities[i] = GroupMemberModelToEntity(model)
	}
	return entities, nil
}
