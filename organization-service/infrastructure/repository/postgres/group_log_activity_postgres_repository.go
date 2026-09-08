package postgres

import (
	"context"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type GroupLogActivityModel struct {
	Id        uint32     `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time  `gorm:"autoCreateTime;index:idx_group_log_activities_created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	CreatedBy uint32     `gorm:"not null;index:idx_group_log_activities_created_by"`
	UpdatedBy uint32     `gorm:"not null"`
	GroupId   uint32     `gorm:"not null;index:idx_group_log_activities_group_id"`
	ActorId   uint32     `gorm:"not null;index:idx_group_log_activities_actor_id"`
	LogType   string     `gorm:"not null;index:idx_group_log_activities_type"`
	LogData   string     `gorm:"not null"`
	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}

func (GroupLogActivityModel) TableName() string {
	return "group_log_activities"
}

func GroupLogActivityModelToEntity(model *GroupLogActivityModel) *entity.GroupLogActivity {
	return &entity.GroupLogActivity{
		Id:        model.Id,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		CreatedBy: model.CreatedBy,
		UpdatedBy: model.UpdatedBy,
		GroupId:   model.GroupId,
		ActorId:   model.ActorId,
		LogType:   model.LogType,
		LogData:   model.LogData,
	}
}

func GroupLogActivityEntityToModel(entity *entity.GroupLogActivity) *GroupLogActivityModel {
	return &GroupLogActivityModel{
		Id:        entity.Id,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
		GroupId:   entity.GroupId,
		ActorId:   entity.ActorId,
		LogType:   entity.LogType,
		LogData:   entity.LogData,
		IsDeleted: false,
	}
}

// @bind: organization/internal/domain/repository.GroupLogActivityRepository
type GroupLogActivityPostgresRepository struct {
	db *gorm.DB
}

func NewGroupLogActivityPostgresRepository(db *gorm.DB) *GroupLogActivityPostgresRepository {
	return &GroupLogActivityPostgresRepository{db: db}
}

func (r *GroupLogActivityPostgresRepository) Create(ctx context.Context, activity *entity.GroupLogActivity) (*entity.GroupLogActivity, error) {
	model := GroupLogActivityEntityToModel(activity)
	model.IsDeleted = false
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return GroupLogActivityModelToEntity(model), nil
}

func (r *GroupLogActivityPostgresRepository) Update(ctx context.Context, activity *entity.GroupLogActivity) (*entity.GroupLogActivity, error) {
	model := GroupLogActivityEntityToModel(activity)
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", activity.Id, false).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return GroupLogActivityModelToEntity(model), nil
}

func (r *GroupLogActivityPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&GroupLogActivityModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *GroupLogActivityPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.GroupLogActivity, error) {
	var model GroupLogActivityModel
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error
	if err != nil {
		return nil, err
	}
	return GroupLogActivityModelToEntity(&model), nil
}

func (r *GroupLogActivityPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupLogActivity, uint32, error) {
	var models []*GroupLogActivityModel
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND is_deleted = ?", groupID, false).
		Offset(page * size).
		Limit(size).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.GroupLogActivity, len(models))
	for i, model := range models {
		entities[i] = GroupLogActivityModelToEntity(model)
	}
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&GroupLogActivityModel{}).
		Where("group_id = ? AND is_deleted = ?", groupID, false).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}

func (r *GroupLogActivityPostgresRepository) GetByActorID(ctx context.Context, actorID uint32) ([]*entity.GroupLogActivity, error) {
	var models []*GroupLogActivityModel
	if err := r.db.WithContext(ctx).Where("actor_id = ?", actorID).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.GroupLogActivity, len(models))
	for i, model := range models {
		entities[i] = GroupLogActivityModelToEntity(model)
	}
	return entities, nil
}
