package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"gorm.io/gorm"
)

type GroupDocumentModel struct {
	Id          uint32    `gorm:"primaryKey;autoIncrement"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	CreatedBy   uint32    `gorm:"not null"`
	UpdatedBy   uint32    `gorm:"not null"`

	GroupId     uint32 `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	FileUrl     string `gorm:"not null"`
	FileType    string `gorm:"not null"`
	FileSize    uint64 `gorm:"not null"`
}

func (GroupDocumentModel) TableName() string {
	return "group_documents"
}

func GroupDocumentModelToEntity(model *GroupDocumentModel) *entity.GroupDocument {
	return &entity.GroupDocument{
		Id:          model.Id,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		CreatedBy:   model.CreatedBy,
		UpdatedBy:   model.UpdatedBy,
		GroupId:     model.GroupId,
		Name:        model.Name,
		Description: model.Description,
		FileUrl:     model.FileUrl,
		FileType:    model.FileType,
		FileSize:    model.FileSize,
	}
}

func GroupDocumentEntityToModel(entity *entity.GroupDocument) *GroupDocumentModel {
	return &GroupDocumentModel{
		Id:          entity.Id,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
		GroupId:     entity.GroupId,
		Name:        entity.Name,
		Description: entity.Description,
		FileUrl:     entity.FileUrl,
		FileType:    entity.FileType,
		FileSize:    entity.FileSize,
	}
}

type groupDocumentPostgresRepository struct {
	db *gorm.DB
}

func NewGroupDocumentPostgresRepository(db *gorm.DB) repository.GroupDocumentRepository {
	return &groupDocumentPostgresRepository{db: db}
}

func (r *groupDocumentPostgresRepository) Create(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error) {
	model := GroupDocumentEntityToModel(document)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return GroupDocumentModelToEntity(model), nil
}

func (r *groupDocumentPostgresRepository) Update(ctx context.Context, document *entity.GroupDocument) (*entity.GroupDocument, error) {
	model := GroupDocumentEntityToModel(document)
	if err := r.db.WithContext(ctx).Model(&GroupDocumentModel{}).Where("id = ?", model.Id).Updates(model).Error; err != nil {
		return nil, err
	}
	return GroupDocumentModelToEntity(model), nil
}

func (r *groupDocumentPostgresRepository) Delete(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&GroupDocumentModel{}, id).Error
}

func (r *groupDocumentPostgresRepository) GetByID(ctx context.Context, id uint32) (*entity.GroupDocument, error) {
	var model *GroupDocumentModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return GroupDocumentModelToEntity(model), nil
}

func (r *groupDocumentPostgresRepository) GetByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupDocument, error) {
	var models []*GroupDocumentModel
	if err := r.db.WithContext(ctx).Where("group_id = ?", groupID).Find(&models).Error; err != nil {
		return nil, err
	}
	entities := make([]*entity.GroupDocument, len(models))
	for i, model := range models {
		entities[i] = GroupDocumentModelToEntity(model)
	}
	return entities, nil
} 