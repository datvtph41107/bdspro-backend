package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type OrganizationLogActivityModel struct {
	Id uint32 `gorm:"primaryKey;autoIncrement"`	
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32 `gorm:"not null"`
	UpdatedBy uint32 `gorm:"not null"`

	OrganizationId uint32 `gorm:"not null"`
	ActorId uint32 `gorm:"not null"`
	LogType string `gorm:"not null"`
	LogData string `gorm:"not null"`
}

func (OrganizationLogActivityModel) TableName() string {
	return "organization_log_activities"
}

func OrganizationLogActivityModelToEntity(model *OrganizationLogActivityModel) *entity.OrganizationLogActivity {
	return &entity.OrganizationLogActivity{
		Id: model.Id,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		CreatedBy: model.CreatedBy,
		UpdatedBy: model.UpdatedBy,
		OrganizationId: model.OrganizationId,
		ActorId: model.ActorId,
		LogType: model.LogType,
		LogData: model.LogData,
	}
}

func OrganizationLogActivityEntityToModel(entity *entity.OrganizationLogActivity) *OrganizationLogActivityModel {
	return &OrganizationLogActivityModel{
		Id: entity.Id,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		CreatedBy: entity.CreatedBy,
		UpdatedBy: entity.UpdatedBy,
		OrganizationId: entity.OrganizationId,
		ActorId: entity.ActorId,
		LogType: entity.LogType,
		LogData: entity.LogData,
	}
}

type organizationLogActivityPostgresRepository struct {
	db *gorm.DB
}

func NewOrganizationLogActivityPostgresRepository(db *gorm.DB) repository.OrganizationLogActivityRepository {
	return &organizationLogActivityPostgresRepository{db: db}
}

func (r *organizationLogActivityPostgresRepository) Create(ctx context.Context, log *entity.OrganizationLogActivity) (*entity.OrganizationLogActivity, error) {
	model := OrganizationLogActivityEntityToModel(log)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return OrganizationLogActivityModelToEntity(model), nil
}

func (r *organizationLogActivityPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.OrganizationLogActivity, error) {
	var model OrganizationLogActivityModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return OrganizationLogActivityModelToEntity(&model), nil
}

func (r *organizationLogActivityPostgresRepository) GetByOrganizationID(ctx context.Context, organizationID uint32) ([]*entity.OrganizationLogActivity, error) {
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

func (r *organizationLogActivityPostgresRepository) GetByActorID(ctx context.Context, actorID uint32) ([]*entity.OrganizationLogActivity, error) {
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