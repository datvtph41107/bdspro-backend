package postgres

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type GroupSettingModel struct {
	Id        uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_group_settings_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"not null;index:idx_group_settings_created_by"`
	UpdatedBy uint32    `gorm:"not null"`

	GroupId     uint32 `gorm:"not null;index:idx_group_settings_group_id"`
	ConfigKey   string `gorm:"not null;index:idx_group_settings_config_key"`
	ConfigValue string `gorm:"not null"`
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
	}
}

type groupSettingPostgresRepository struct {
	db *gorm.DB
}

func NewGroupSettingPostgresRepository(db *gorm.DB) repository.GroupSettingRepository {
	return &groupSettingPostgresRepository{db: db}
}

func (r *groupSettingPostgresRepository) Create(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error) {
	model := GroupSettingEntityToModel(setting)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupSettingModelToEntity(model), nil
}

func (r *groupSettingPostgresRepository) Update(ctx context.Context, setting *entity.GroupSetting) (*entity.GroupSetting, error) {
	model := GroupSettingEntityToModel(setting)
	if err := r.db.WithContext(ctx).Model(&GroupSettingModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return GroupSettingModelToEntity(model), nil
}

func (r *groupSettingPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&GroupSettingModel{}, id).Error
}

func (r *groupSettingPostgresRepository) GetByConfigKey(ctx context.Context, configKey string) (*entity.GroupSetting, error) {
	var model *GroupSettingModel
	if err := r.db.WithContext(ctx).Where("config_key = ?", configKey).First(&model).Error; err != nil {
		return nil, err
	}
	return GroupSettingModelToEntity(model), nil
}

func (r *groupSettingPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.GroupSetting, error) {
	var model *GroupSettingModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return GroupSettingModelToEntity(model), nil
}

func (r *groupSettingPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.GroupSetting, uint32, error) {
	var models []*GroupSettingModel
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Offset(page * size).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.GroupSetting, len(models))
	for i, model := range models {
		entities[i] = GroupSettingModelToEntity(model)
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&GroupSettingModel{}).Where("group_id = ?", groupID).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(count), nil
}