package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type GroupLogActivityModel struct {
	Id uint32 `gorm:"primaryKey;autoIncrement"`	
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_group_log_activities_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32 `gorm:"not null;index:idx_group_log_activities_created_by"`
	UpdatedBy uint32 `gorm:"not null"`

	GroupId uint32 `gorm:"not null;index:idx_group_log_activities_group_id"`
	ActorId uint32 `gorm:"not null;index:idx_group_log_activities_actor_id"`
	LogType string `gorm:"not null;index:idx_group_log_activities_type"`
	LogData string `gorm:"not null"`
}

func (GroupLogActivityModel) TableName() string {
	return "group_log_activities"
}

func GroupLogActivityModelToEntity(model *GroupLogActivityModel) *entity.GroupLogActivity {
	return &entity.GroupLogActivity{
		Id: model.Id,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		CreatedBy: model.CreatedBy,
		UpdatedBy: model.UpdatedBy,
		GroupId: model.GroupId,
		ActorId: model.ActorId,
		LogType: model.LogType,
		LogData: model.LogData,
	}
}

func GroupLogActivityEntityToModel(entity *entity.GroupLogActivity) *GroupLogActivityModel {
	return &GroupLogActivityModel{
		Id: entity.Id,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
		GroupId: entity.GroupId,
		ActorId: entity.ActorId,
		LogType: entity.LogType,
		LogData: entity.LogData,
	}
}

type groupLogActivityPostgresRepository struct {
	db *gorm.DB
}

func NewGroupLogActivityPostgresRepository(db *gorm.DB) repository.GroupLogActivityRepository {
	return &groupLogActivityPostgresRepository{db: db}
}

func (r *groupLogActivityPostgresRepository) Create(ctx context.Context, log *entity.GroupLogActivity) (*entity.GroupLogActivity, error) {
	model := GroupLogActivityEntityToModel(log)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupLogActivityModelToEntity(model), nil
}

func (r *groupLogActivityPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.GroupLogActivity, error) {
	var model *GroupLogActivityModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return GroupLogActivityModelToEntity(model), nil
}

func (r *groupLogActivityPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupLogActivity, uint32, error) {
	var models []*GroupLogActivityModel
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Offset(page * size).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.GroupLogActivity, len(models))
	for i, model := range models {
		entities[i] = GroupLogActivityModelToEntity(model)
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&GroupLogActivityModel{}).Where("group_id = ?", groupID).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}

func (r *groupLogActivityPostgresRepository) GetByActorID(ctx context.Context, actorID uint32) ([]*entity.GroupLogActivity, error) {
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