package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"organization/internal/domain/entity"
)

type GroupSettingModel struct {
	Id          uint32     `gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time  `gorm:"autoCreateTime;index:idx_group_settings_created_at"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
	CreatedBy   uint32     `gorm:"not null;index:idx_group_settings_created_by"`
	UpdatedBy   uint32     `gorm:"not null"`
	GroupId     uint32     `gorm:"not null;index:idx_group_settings_group_id"`
	ConfigKey   string     `gorm:"not null;index:idx_group_settings_config_key"`
	ConfigValue string     `gorm:"not null"`
	DeletedAt   *time.Time `gorm:"index"`
	IsDeleted   bool       `gorm:"default:false"`
}

func (GroupSettingModel) TableName() string {
	return "group_settings"
}

func GroupSettingModelToEntity(model *GroupSettingModel) *entity.GroupSetting {
	return &entity.GroupSetting{
		Id:          model.Id,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		CreatedBy:   model.CreatedBy,
		UpdatedBy:   model.UpdatedBy,
		GroupId:     model.GroupId,
		ConfigKey:   model.ConfigKey,
		ConfigValue: model.ConfigValue,
	}
}

func GroupSettingEntityToModel(entity *entity.GroupSetting) *GroupSettingModel {
	return &GroupSettingModel{
		Id:          entity.Id,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
		GroupId:     entity.GroupId,
		ConfigKey:   entity.ConfigKey,
		ConfigValue: entity.ConfigValue,
		IsDeleted:   false,
	}
}

// @bind: organization/internal/domain/repository.GroupSettingRepository
type GroupSettingPostgresRepository struct {
	db *gorm.DB
}

func NewGroupSettingPostgresRepository(db *gorm.DB) *GroupSettingPostgresRepository {
	return &GroupSettingPostgresRepository{db: db}
}

func (r *GroupSettingPostgresRepository) Create(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error) {
	model := GroupSettingEntityToModel(setting)
	model.IsDeleted = false
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return nil, err
	}
	return GroupSettingModelToEntity(model), nil
}

func (r *GroupSettingPostgresRepository) Update(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error) {
	model := GroupSettingEntityToModel(setting)
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", setting.Id, false).Updates(model).Error
	if err != nil {
		return nil, err
	}
	return GroupSettingModelToEntity(model), nil
}

func (r *GroupSettingPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Model(&GroupSettingModel{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *GroupSettingPostgresRepository) GetByConfigKey(ctx context.Context, configKey string) (*entity.GroupSetting, error) {
	var model *GroupSettingModel
	if err := r.db.WithContext(ctx).Where("config_key = ? AND is_deleted = ?", configKey, false).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupSettingModelToEntity(model), nil
}

func (r *GroupSettingPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.GroupSetting, error) {
	var model GroupSettingModel
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error
	if err != nil {
		return nil, err
	}
	return GroupSettingModelToEntity(&model), nil
}

func (r *GroupSettingPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupSetting, uint32, error) {
	var models []*GroupSettingModel
	if err := r.db.WithContext(ctx).
		Where("group_id = ? AND is_deleted = ?", groupID, false).
		Offset(page * size).
		Limit(size).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.GroupSetting, len(models))
	for i, model := range models {
		entities[i] = GroupSettingModelToEntity(model)
	}
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&GroupSettingModel{}).
		Where("group_id = ? AND is_deleted = ?", groupID, false).
		Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}
