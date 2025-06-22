package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type GroupNotificationModel struct {
	Id        uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_group_notifications_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"not null;index:idx_group_notifications_created_by"`
	UpdatedBy uint32    `gorm:"not null"`

	GroupId uint32 `gorm:"not null;index:idx_group_notifications_group_id"`
	Type    string `gorm:"not null;index:idx_group_notifications_type"`
	Content string `gorm:"not null"`

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool `gorm:"default:false"`
}

func (GroupNotificationModel) TableName() string {
	return "group_notifications"
}

func GroupNotificationModelToEntity(model *GroupNotificationModel) *entity.GroupNotification {
	return &entity.GroupNotification{
		Id:        model.Id,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		CreatedBy: model.CreatedBy,
		UpdatedBy: model.UpdatedBy,
		GroupId:   model.GroupId,
		Type:      model.Type,
		Content:   model.Content,
	}
}

func GroupNotificationEntityToModel(entity *entity.GroupNotification) *GroupNotificationModel {
	return &GroupNotificationModel{
		Id:        entity.Id,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
		GroupId:   entity.GroupId,
		Type:      entity.Type,
		Content:   entity.Content,
	}
}

type groupNotificationPostgresRepository struct {
	db *gorm.DB
}

func NewGroupNotificationPostgresRepository(db *gorm.DB) repository.GroupNotificationRepository {
	return &groupNotificationPostgresRepository{db: db}
}

func (r *groupNotificationPostgresRepository) Create(ctx context.Context, notification *entity.GroupNotification) (*entity.GroupNotification, error) {
	model := GroupNotificationEntityToModel(notification)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupNotificationModelToEntity(model), nil
}

func (r *groupNotificationPostgresRepository) Update(ctx context.Context, notification *entity.GroupNotification) (*entity.GroupNotification, error) {
	model := GroupNotificationEntityToModel(notification)
	if err := r.db.WithContext(ctx).Model(&GroupNotificationModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return GroupNotificationModelToEntity(model), nil
}

func (r *groupNotificationPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&GroupNotificationModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted":  true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *groupNotificationPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.GroupNotification, error) {
	var model *GroupNotificationModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupNotificationModelToEntity(model), nil
}

func (r *groupNotificationPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupNotification, error) {
	var models []*GroupNotificationModel
	if err := r.db.WithContext(ctx).Where("group_id = ? AND is_deleted = ?", groupID, false).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.GroupNotification, len(models))
	for i, model := range models {
		entities[i] = GroupNotificationModelToEntity(model)
	}
	return entities, nil
}

func (r *groupNotificationPostgresRepository) GetByUserID(ctx context.Context, userID uint32) ([]*entity.GroupNotification, error) {
	var models []*GroupNotificationModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_deleted = ?", userID, false).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.GroupNotification, len(models))
	for i, model := range models {
		entities[i] = GroupNotificationModelToEntity(model)
	}

	return entities, nil
}
