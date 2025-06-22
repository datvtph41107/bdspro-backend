package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type GroupModel struct {
	Id        uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_groups_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"not null;index:idx_groups_created_by"`
	UpdatedBy uint32    `gorm:"not null"`

	Name        string `gorm:"not null;index:idx_groups_name"`
	Description string
	AvatarUrl   string

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}

func (GroupModel) TableName() string {
	return "groups"
}

func GroupModelToEntity(model *GroupModel) *entity.Group {
	return &entity.Group{
		Id:          model.Id,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		CreatedBy:   model.CreatedBy,
		UpdatedBy:   model.UpdatedBy,
		Name:        model.Name,
		Description: model.Description,
		AvatarUrl:   model.AvatarUrl,
	}
}

func GroupEntityToModel(entity *entity.Group) *GroupModel {
	return &GroupModel{
		Id:          entity.Id,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
		Name:        entity.Name,
		Description: entity.Description,
		AvatarUrl:   entity.AvatarUrl,
	}
}

type groupPostgresRepository struct {
	db *gorm.DB
}

func NewGroupPostgresRepository(db *gorm.DB) repository.GroupRepository {
	return &groupPostgresRepository{db: db}
}

func (r *groupPostgresRepository) Create(ctx context.Context, group *entity.Group) (*entity.Group, error) {
	model := GroupEntityToModel(group)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupModelToEntity(model), nil
}

func (r *groupPostgresRepository) Update(ctx context.Context, group *entity.Group) (*entity.Group, error) {
	model := GroupEntityToModel(group)
	if err := r.db.WithContext(ctx).Model(&GroupModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return GroupModelToEntity(model), nil
}

func (r *groupPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&GroupModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted":  true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *groupPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.Group, error) {
	var model *GroupModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupModelToEntity(model), nil
}

func (r *groupPostgresRepository) List(ctx context.Context, page, size int, currentUserId uint32) ([]*entity.Group, uint32, error) {
	var models []*GroupModel
	if err := r.db.WithContext(ctx).Offset(page*size).Limit(size).Where("created_by = ? AND is_deleted = ?", currentUserId, false).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.Group, len(models))
	for i, model := range models {
		entities[i] = GroupModelToEntity(model)
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&GroupModel{}).Where("created_by = ? AND is_deleted = ?", currentUserId, false).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}

func (r *groupPostgresRepository) FindByUserId(ctx context.Context, userId uint32) ([]*entity.Group, error) {
	var groups []*GroupModel
	if err := r.db.WithContext(ctx).
		Joins("LEFT JOIN group_members ON groups.id = group_members.group_id").
		Where("groups.created_by = ? OR group_members.user_id = ? AND groups.is_deleted = ?", userId, userId, false).
		Distinct().
		Order("groups.created_at DESC").
		Find(&groups).Error; err != nil {
		return nil, err
	}
	return GroupModelsToEntities(groups), nil
}

func GroupModelsToEntities(models []*GroupModel) []*entity.Group {
	entities := make([]*entity.Group, len(models))
	for i, model := range models {
		entities[i] = GroupModelToEntity(model)
	}
	return entities
}
