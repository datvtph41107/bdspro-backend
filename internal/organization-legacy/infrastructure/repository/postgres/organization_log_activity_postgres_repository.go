package postgres

import (
	"context"
	"time"

	"organization/internal/domain/entity"

	"gorm.io/gorm"
)

type OrganizationLogActivityModel struct {
	Id        uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time `gorm:"autoCreateTime;index:idx_org_log_activities_created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"not null;index:idx_org_log_activities_created_by"`
	UpdatedBy uint32    `gorm:"not null"`

	OrganizationId uint32 `gorm:"not null;index:idx_org_log_activities_org_id"`
	ActorId        uint32 `gorm:"not null;index:idx_org_log_activities_actor_id"`
	LogType        string `gorm:"not null;index:idx_org_log_activities_type"`
	LogData        string `gorm:"not null"`
}

func (OrganizationLogActivityModel) TableName() string {
	return "organization_log_activities"
}

func OrganizationLogActivityModelToEntity(model *OrganizationLogActivityModel) *entity.OrganizationLogActivity {
	return &entity.OrganizationLogActivity{
		Id:             model.Id,
		CreatedAt:      model.CreatedAt,
		UpdatedAt:      model.UpdatedAt,
		CreatedBy:      model.CreatedBy,
		UpdatedBy:      model.UpdatedBy,
		OrganizationId: model.OrganizationId,
		ActorId:        model.ActorId,
		LogType:        model.LogType,
		LogData:        model.LogData,
	}
}

func OrganizationLogActivityEntityToModel(entity *entity.OrganizationLogActivity) *OrganizationLogActivityModel {
	return &OrganizationLogActivityModel{
		Id:             entity.Id,
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		CreatedBy:      entity.CreatedBy,
		UpdatedBy:      entity.UpdatedBy,
		OrganizationId: entity.OrganizationId,
		ActorId:        entity.ActorId,
		LogType:        entity.LogType,
		LogData:        entity.LogData,
	}
}

// @bind: organization/internal/domain/repository.OrganizationLogActivityRepository
type OrganizationLogActivityPostgresRepository struct {
	db *gorm.DB
}

func NewOrganizationLogActivityPostgresRepository(db *gorm.DB) *OrganizationLogActivityPostgresRepository {
	return &OrganizationLogActivityPostgresRepository{db: db}
}

func (r *OrganizationLogActivityPostgresRepository) Create(ctx context.Context, log *entity.OrganizationLogActivity) (*entity.OrganizationLogActivity, error) {
	model := OrganizationLogActivityEntityToModel(log)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return OrganizationLogActivityModelToEntity(model), nil
}

func (r *OrganizationLogActivityPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.OrganizationLogActivity, error) {
	var model OrganizationLogActivityModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return OrganizationLogActivityModelToEntity(&model), nil
}

func (r *OrganizationLogActivityPostgresRepository) GetByOrganizationID(ctx context.Context, organizationID uint32) ([]*entity.OrganizationLogActivity, error) {
	var models []*OrganizationLogActivityModel
	if err := r.db.WithContext(ctx).Where("organization_id = ?", organizationID).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.OrganizationLogActivity, len(models))
	for i, model := range models {
		entities[i] = OrganizationLogActivityModelToEntity(model)
	}
	return entities, nil
}

func (r *OrganizationLogActivityPostgresRepository) GetByActorID(ctx context.Context, actorID uint32) ([]*entity.OrganizationLogActivity, error) {
	var models []*OrganizationLogActivityModel
	if err := r.db.WithContext(ctx).Where("actor_id = ?", actorID).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.OrganizationLogActivity, len(models))
	for i, model := range models {
		entities[i] = OrganizationLogActivityModelToEntity(model)
	}
	return entities, nil
}

func (r *OrganizationLogActivityPostgresRepository) GetByOrganizationIDWithPagination(ctx context.Context, organizationID uint32, page, size int) ([]*entity.OrganizationLogActivity, uint32, error) {
	var models []*OrganizationLogActivityModel
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND is_deleted = ?", organizationID, false).Offset(page * size).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.OrganizationLogActivity, len(models))
	for i, model := range models {
		entities[i] = OrganizationLogActivityModelToEntity(model)
	}
	var total int64
	if err := r.db.WithContext(ctx).Model(&OrganizationLogActivityModel{}).Where("organization_id = ? AND is_deleted = ?", organizationID, false).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(total), nil
}

func (r *OrganizationLogActivityPostgresRepository) GetByActorIDWithPagination(ctx context.Context, actorID uint32, page, size int) ([]*entity.OrganizationLogActivity, uint32, error) {
	var models []*OrganizationLogActivityModel
	if err := r.db.WithContext(ctx).Where("actor_id = ? AND is_deleted = ?", actorID, false).Offset(page * size).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	entities := make([]*entity.OrganizationLogActivity, len(models))
	for i, model := range models {
		entities[i] = OrganizationLogActivityModelToEntity(model)
	}
	var total int64
	if err := r.db.WithContext(ctx).Model(&OrganizationLogActivityModel{}).Where("actor_id = ? AND is_deleted = ?", actorID, false).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return entities, uint32(total), nil
}
