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
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	CreatedBy uint32    `gorm:"not null"`
	UpdatedBy uint32    `gorm:"not null"`

	Name        string `gorm:"not null"`
	Description string
	AvatarUrl   string
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
	return r.db.WithContext(ctx).Delete(&GroupModel{}, id).Error
}

func (r *groupPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.Group, error) {
	var model *GroupModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return GroupModelToEntity(model), nil
}

func (r *groupPostgresRepository) List(ctx context.Context, offset, limit int) ([]*entity.Group, error) {
	var models []*GroupModel
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.Group, len(models))
	for i, model := range models {
		entities[i] = GroupModelToEntity(model)
	}
	return entities, nil
}
